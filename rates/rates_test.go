package rates

import "testing"

func TestFetch(t *testing.T) {
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
