package rates

import (
	"fmt"
	"testing"
)

func TestCountryRates_Rate(t *testing.T) {
	c, _ := Country("NL")

	if r, _ := c.Rate("standard"); r != 21 {
		t.Errorf("Standard VAT rate for NL is supposed to be 21. Got %.2f", r)
	}

	if r, _ := c.Rate("reduced"); r != 6 {
		t.Errorf("Reduced VAT rate for NL is supposed to be 6. Got %.2f", r)
	}

	c, _ = Country("RO")
	if r, _ := c.Rate("standard"); r != 20 {
		t.Errorf("Standard VAT rate for RO is supposed to be 20. Got %.2f", r)
	}
}

func ExampleCountryRates_Rate() {
	c, _ := Country("NL")
	r, _ := c.Rate("standard")

	fmt.Printf("Standard VAT rate for %s is %.2f", c.Name, r)
	// Output: Standard VAT rate for Netherlands is 21.00
}
