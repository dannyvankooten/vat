package vat

import (
	"fmt"
	"testing"
)

var tests = []struct {
	number string
	valid  bool
}{
	{"", false},
	{"A", false},
	{"AB123A01", false},
	{"ATU12345678", true},
	{"ATU15673009", true},
	{"ATU1234567", false},
	{"BE0123456789", true},
	{"BE1234567891", true},
	{"BE0999999999", true},
	{"BE9999999999", true},
	{"BE012345678", false},
	{"BE123456789", false},
	{"BG123456789", true},
	{"BG1234567890", true},
	{"BG1234567", false},
	{"CHE-156.730.098 MWST", true},
	{"CHE-156.730.098", true},
	{"CHE156730098MWST", true},
	{"CHE156730098", true},
	{"CY12345678X", true},
	{"CY15673009L", true},
	{"CY1234567X", false},
	{"CZ12345678", true},
	{"CZ1234567", false},
	{"DE123456789", true},
	{"DE12345678", false},
	{"DK12345678", true},
	{"DK1234567", false},
	{"EE123456789", true},
	{"EE12345678", false},
	{"EL123456789", true},
	{"EL12345678", false},
	{"ESX12345678", true},
	{"ESX1234567", false},
	{"FI1234567", false},
	{"FI12345678", true},
	{"FR12345678901", true},
	{"FR1234567890", false},
	{"GB999999973", true},
	{"GB156730098481", true},
	{"GBGD549", true},
	{"GBHA549", true},
	{"GB99999997", false},
	{"HU12345678", true},
	{"HU1234567", false},
	{"HR12345678901", true},
	{"HR1234567890", false},
	{"IE1234567X", true},
	{"IE123456X", false},
	{"IT12345678901", true},
	{"IT1234567890", false},
	{"LT123456789", true},
	{"LT12345678", false},
	{"LU26375245", true},
	{"LU12345678", true},
	{"LU1234567", false},
	{"LV12345678901", true},
	{"LV1234567890", false},
	{"MT12345678", true},
	{"MT1234567", false},
	{"NL123456789B01", true},
	{"NL123456789B12", true},
	{"NL12345678B12", false},
	{"PL1234567890", true},
	{"PL123456789", false},
	{"PT123456789", true},
	{"PT12345678", false},
	{"RO123456789", true},
	{"RO1", false}, // Romania has a really weird VAT format...
	{"SE123456789012", true},
	{"SE12345678901", false},
	{"SI12345678", true},
	{"SI1234567", false},
	{"SK1234567890", true},
	{"SK123456789", false},
}

func BenchmarkValidateFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateNumberFormat("NL" + string(i))
	}
}

func TestValidateNumber(t *testing.T) {
	for _, test := range tests {
		valid, err := ValidateNumberFormat(test.number)
		if err != nil {
			panic(err)
		}

		if test.valid != valid {
			t.Errorf("Expected %v for %v, got %v", test.valid, test.number, valid)
		}
	}
}

func ExampleValidateNumber() {
	vatNumber := "IE6388047V"
	valid, _ := ValidateNumber(vatNumber)
	fmt.Printf("Is %s valid: %t", vatNumber, valid)
	// Output: Is IE6388047V valid: true
}

func TestValidateNumberFormat(t *testing.T) {
	for _, test := range tests {
		valid, err := ValidateNumberFormat(test.number)

		if err != nil {
			panic(err)
		}

		if test.valid != valid {
			t.Errorf("Expected %v for %v, got %v", test.valid, test.number, valid)
		}

	}
}

func TestValidateNumberExistence(t *testing.T) {
	valid, _ := ValidateNumberExistence("IE6388047V")
	if !valid {
		t.Error("IE6388047V is a valid VAT number.")
	}

	valid, _ = ValidateNumberExistence("NL123456789B01")
	if valid {
		t.Error("NL123456789B01 is not a valid VAT number.")
	}
}
