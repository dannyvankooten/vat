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
	{"ATU1234567", false},
	{"BE012345678", false},
	{"BE123456789", false},
	{"BG1234567", false},
	{"CY1234567X", false},
	{"CZ1234567", false},
	{"DE12345678", false},
	{"DK1234567", false},
	{"EE12345678", false},
	{"EL12345678", false},
	{"ESX1234567", false},
	{"FI1234567", false},
	{"FR1234567890", false},
	{"GB99999997", false},
	{"HU1234567", false},
	{"HR1234567890", false},
	{"IE123456X", false},
	{"IT1234567890", false},
	{"LT12345678", false},
	{"LU1234567", false},
	{"LV1234567890", false},
	{"MT1234567", false},
	{"NL12345678B12", false},
	{"PL123456789", false},
	{"PT12345678", false},
	{"RO1", false}, // Romania has a really weird VAT format...
	{"SE12345678901", false},
	{"SI1234567", false},
	{"SK123456789", false},
	{"AB123A01", false},
	{"LU26375245", true},
	{"NL123456789B01", true},
	{"ATU12345678", true},
	{"BE0123456789", true},
	{"BE1234567891", true},
	{"BG123456789", true},
	{"BG1234567890", true},
	{"CY12345678X", true},
	{"CZ12345678", true},
	{"DE123456789", true},
	{"DK12345678", true},
	{"EE123456789", true},
	{"EL123456789", true},
	{"ESX12345678", true},
	{"FI12345678", true},
	{"FR12345678901", true},
	{"GB999999973", true},
	{"HU12345678", true},
	{"HR12345678901", true},
	{"IE1234567X", true},
	{"IT12345678901", true},
	{"LT123456789", true},
	{"LU12345678", true},
	{"LV12345678901", true},
	{"MT12345678", true},
	{"NL123456789B12", true},
	{"PL1234567890", true},
	{"PT123456789", true},
	{"RO123456789", true},
	{"SE123456789012", true},
	{"SI12345678", true},
	{"SK1234567890", true},
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
