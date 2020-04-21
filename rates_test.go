package vat

import (
	"testing"
	"time"
)

func TestCountryRates_GetRate(t *testing.T) {
	c, _ := GetCountryRates("NL")

	if r, _ := c.GetRate("standard"); r != 21 {
		t.Errorf("Standard VAT rate for NL is supposed to be 21. Got %.2f", r)
	}

	if r, _ := c.GetRate("reduced"); r != 9 {
		t.Errorf("Reduced VAT rate for NL is supposed to be 9. Got %.2f", r)
	}

	c, _ = GetCountryRates("RO")
	if r, _ := c.GetRate("standard"); r != 19 {
		t.Errorf("Standard VAT rate for RO is supposed to be 19. Got %.2f", r)
	}
}

func TestCountryRates_GetRateOn(t *testing.T) {
	c, _ := GetCountryRates("NL")
	time, _ := time.Parse("2006-01-01", "2002-01-01")
	if r, _ := c.GetRateOn(time, "standard"); r != 19 {
		t.Errorf("Standard VAT rate for NL in 2002 is supposed to be 19. Got %.2f", r)
	}
}
