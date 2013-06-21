package main

import (
	"strconv"
	"time"
)

type CompanyState int

const (
	Active CompanyState = iota
	Deregistered
	NotInBusiness
)

type Address struct {
	Street      string `json:"street"`
	HouseNumber int    `json:"number"`
	Postcode    int    `json:"postcode"`
	Place       string `json:"place"`
}

type ISATType struct {
	Number      int    `json:"number"`
	Description string `json:"description"`
	Main        bool   `json:"is_main"`
}

type VATNumber struct {
	ID         int        `json:"id"`
	DateOpened time.Time  `json:"date_opened"`
	DateClosed time.Time  `json:"date_closed,omitempty"`
	ISATTypes  []ISATType `json:"isat_types"`
}

type Company struct {
	Ssid         string       `json:"ssid"`
	Name         string       `json:"name"`
	PostAddress  Address      `json:"post_address,omitempty"`
	LegalAddress Address      `json:"legal_address,omitempty"`
	Type         string       `json:"company_type"`
	VATNumbers   []VATNumber  `json:"vat_numbers,omitempty"`
	State        CompanyState `json:"company_state"`
}

func (c Company) ShouldGetDetails() bool {
	if c.State != Active {
		return false
	}
	x, _ := strconv.Atoi(string(c.Ssid[0]))
	if x < 4 {
		return false
	}
	return true
}
