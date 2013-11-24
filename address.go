package fyrirtaekjaskra

import (
	"regexp"
	"strconv"
)

var (
	reAddress = regexp.MustCompile(`(?P<street>[\p{Latin}-]+)\s(?:(?P<number>[a-zA-Z0-9]+))?,?\s+(?P<postcode>\d{3}) (?P<place>[\p{Latin}]+)`)
)

// ParseAddress parses a string and returns an address
func ParseAddress(s string) (a Address, err error) {

	if reAddress.MatchString(s) {
		parts := reAddress.FindStringSubmatch(s)[1:]
		a.Street = parts[0]
		a.HouseNumber = parts[1]
		a.Postcode, err = strconv.Atoi(parts[2])
		a.Place = parts[3]
	} else {
		logger.Warningf("Parse address: \"%s\"", s)
	}

	return
}
