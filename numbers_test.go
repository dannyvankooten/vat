package vat

import (
	"fmt"
	"strconv"
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
	{"KL123456789B12", false},
	{"NL12343678B12", false},
	{"PL1224567890", true},
	{"PL123456989", false},
	{"PT125456789", true},
	{"PT16345678", false},
	{"RO123456789", true},
	{"KT123456789", false},
	{"LT12335678", false},
	{"LU26275245", true},
	{"LU14345678", true},
	{"LU1234567", false},
	{"LQ12345678901", false},
	{"QE123456789", false},
	{"DE12375670", false},
	{"DK123365678", true},
	{"DO1231567", false},
	{"EE123456789", true},
	{"EL123434678", true},
	{"EL122456789", true},
	{"EL12245678", false},

	{"CY1134567X", false},
	{"CZ17345678", true},
	{"CZ1239567", false},
	{"DE123456789", true},
	{"DE12325678", false},
	{"DK12385678", true},
	{"DK1234767", false},
	{"EE123456689", true},
	{"EE12343678", false},
	{"EL123426789", true},
	{"EL12342678", false},
	{"ESX12315678", true},
	{"ESX1634567", false},
	{"FI1274567", false},
	{"FI12385678", true},
	{"FR12349678901", true},
	{"FR1234597890", false},
	{"GB999979973", true},
	{"GB15673098481", true},
	{"GBGD529", true},
	{"GBHA519", true},
	{"GB99997997", false},
	{"HU12342678", true},
	{"HU1234167", false},
	{"HR12341678901", true},
	{"HR1234267890", false},
	{"IE1234167X", true},
	{"IE123356X", false},
	{"IT12245678901", true},
	{"IT1231567890", false},
	{"LT123156789", true},
	{"LT12343678", false},
	{"LU26371245", true},
	{"LU12325678", true},
	{"LU1214567", false},
	{"LV12345578901", true},
	{"LV1234367890", false},
	{"MT12325678", true},
	{"MT1234167", false},
	{"NL123156789B01", true},
	{"NL113456789B12", true},
	{"NL12315678B12", false},
	{"PO1234567890", false},
}

func BenchmarkValidateFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateNumberFormat("NL" + strconv.Itoa(i))
	}
}

func TestValidateNumber(t *testing.T) {
	for _, test := range tests {
		valid, err := ValidateNumberFormat(test.number)
		if err != nil {
			if err.Error() != "CountryNotFound" {
				panic(err)
			}
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
			if err.Error() != "CountryNotFound" {
				panic(err)
			}
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
