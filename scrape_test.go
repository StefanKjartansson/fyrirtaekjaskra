package fyrirtaekjaskra

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func read(filename string) []byte {

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("unable to read a file")
	}
	return contents
}

func TestParseDetails(t *testing.T) {

	ssid := "5902697199"

	c := Company{
		Ssid: ssid,
	}

	scraper := NewScraper()
	scraper.ParseDetails(read(fmt.Sprintf("./test/fskra-%s.html", ssid)), &c)
	err := <-scraper.ErrChan

	if err != nil {
		t.Error("Parsing has error: %s", err.Error())
		return
	}

	if c.VATNumbers[0].ID != 10487 {
		t.Errorf("vatnumber is weird: %v\n", c.VATNumbers)
	}

}

func TestXpathSearchTable(t *testing.T) {

	scraper := NewScraper()
	scraper.ParseSearchResults(read("./test/fskra-leit.html"))
	err := <-scraper.ErrChan

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	c := <-scraper.CompanyChan
	if c.Ssid != "5407051000" ||
		c.Name != "A Einn ehf" {
		t.Errorf("Parsing error, %+v", c)
		return
	}

}
