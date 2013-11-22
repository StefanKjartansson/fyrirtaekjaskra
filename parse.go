package fyrirtaekjaskra

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	shortForm = "02.01.2006"
)

var (
	deregRegex         = regexp.MustCompile("(i?)Félag afskráð")
	notInBusinessRegex = regexp.MustCompile("(i?)Rekstri hætt")
	ehfRegex           = regexp.MustCompile("(i?)ehf")
)

type Scraper struct {
	CompanyChan chan Company
	ErrChan     chan error
}

func NewScraper() *Scraper {
	return &Scraper{
		CompanyChan: make(chan Company),
		ErrChan:     make(chan error),
	}
}

//Parses details page, TODO: better error handling
func (s *Scraper) ParseDetails(r io.Reader, c *Company) (err error) {

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	doc.Find(".company .boxbody table:nth-of-type(1)>tbody>tr>td").Each(func(i int, s *goquery.Selection) {

		content := strings.Trim(s.Text(), " ")
		log.Println(content)
		switch i {
		case 0:
			(*c).PostAddress, err = ParseAddress(content)
		case 1:
			if content != "" {
				(*c).LegalAddress, err = ParseAddress(content)
			} else {
				(*c).LegalAddress = c.PostAddress
			}
		case 3:
			(*c).Type = content
		}
	})

	vnr := VATNumber{}
	doc.Find(".company .boxbody table.nolines>tbody>tr>td").Each(func(i int, s *goquery.Selection) {

		content := strings.Trim(s.Text(), " ")

		if i > 0 && i%4 == 0 {
			(*c).VATNumbers = append((*c).VATNumbers, vnr)
			vnr = VATNumber{}
		}

		switch i % 4 {
		case 0:
			vnr.ID, _ = strconv.Atoi(content)
		case 1:
			vnr.DateOpened, _ = time.Parse(shortForm, content)
		case 2:
			vnr.DateClosed, _ = time.Parse(shortForm, content)
		case 3:
			vnr.ISATTypes, _ = ParseISATTypes(content)
		}
	})

	// Add last VATNumber
	if vnr.ID > 0 {
		(*c).VATNumbers = append((*c).VATNumbers, vnr)
	}

	address := c.GuessDomain()
	res, xerr := net.LookupHost(address)
	if xerr == nil {
		(*c).Domain = address
		(*c).IPS = res
	}

	return
}

func (s *Scraper) FetchDetails(c Company) {

	content, err := ReadOrGetSSID(c.Ssid)
	if err != nil {
		s.ErrChan <- err
		return
	}
	err = s.ParseDetails(content, &c)
	if err != nil {
		s.ErrChan <- err
		return
	}

	s.CompanyChan <- c
}

func (s *Scraper) ParseSearchResults(r io.Reader) {

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return
	}

	company := Company{Type: "Unknown"}
	doc.Find(".companies .boxbody table>tbody>tr>td").Each(func(idx int, sel *goquery.Selection) {

		if idx > 0 && idx%3 == 0 {
			if company.ShouldGetDetails() {
				go s.FetchDetails(company)
			} else {
				s.CompanyChan <- company
			}
		}

		switch idx % 3 {
		case 0:
			company.Ssid = sel.Find("a").Text()
		case 1:
			content := sel.Text()
			company.Name = strings.Split(content, "\n")[0]
			if deregRegex.MatchString(content) {
				company.State = Deregistered
			} else if notInBusinessRegex.MatchString(content) {
				company.State = NotInBusiness
			}
		case 2:
			content := sel.Text()
			company.PostAddress, _ = ParseAddress(content)
			company.LegalAddress, _ = ParseAddress(content)
		}
	})

	s.CompanyChan <- company
}

func (s *Scraper) ScrapeList(streets []string) {

	for _, street := range streets {

		content, err := ReadOrGetSearch(street)
		if err != nil {
			s.ErrChan <- err
		} else {
			go s.ParseSearchResults(content)
		}
	}

}