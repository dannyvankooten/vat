/*
Package vat helps you deal with European VAT in Go.

It offers VAT number validation using the VIES VAT validation API & VAT rates retrieval using jsonvat.com

Validate a VAT number
		validity := vat.ValidateNumber("NL123456789B01")

Get VAT rate that is currently in effect for a given country
		c, _ := vat.GetCountryRates("NL")
		r, _ := c.GetRate("standard")
*/
package vat

import "errors"

// ErrServiceUnavailable will be returned when VIES VAT validation API or jsonvat.com is unreachable.
var ErrServiceUnavailable = errors.New("vat: service is unreachable")

// ServiceTimeout indicates the number of seconds before a service request times out.
var ServiceTimeout = 10
