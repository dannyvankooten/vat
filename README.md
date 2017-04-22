Package vat
===

[![Go Report Card](https://goreportcard.com/badge/github.com/dannyvankooten/vat)](https://goreportcard.com/report/github.com/dannyvankooten/vat)
[![GoDoc](https://godoc.org/github.com/dannyvankooten/vat?status.svg)](https://godoc.org/github.com/dannyvankooten/vat)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/dannyvankooten/vat/master/LICENSE)

Package for validating VAT numbers & retrieving VAT rates in Go.

## Installation

Use go get.

```
go get github.com/dannyvankooten/vat
```

Then import the package into your own code.

```
import "github.com/dannyvankooten/vat"
```

## Usage

### Validating VAT numbers

VAT numbers can be validated by format, existence or both. VAT numbers are looked up using the [VIES VAT validation API](http://ec.europa.eu/taxation_customs/vies/).

```go
package main

import "github.com/dannyvankooten/vat"

func main() {
  // Validate number by format + existence
  validity, err := vat.ValidateNumber("NL123456789B01")

  // Validate number format
  validity, err := vat.ValidateNumberFormat("NL123456789B01")

  // Validate number existence
  validity, err := vat.ValidateNumberExistence("NL123456789B01")
}
```

### Retrieving VAT rates

To get VAT rate periods for a country, first get a CountryRates struct using the country's ISO-3166-1-alpha2 code.

You can get the rate that is currently in effect using the `GetRate` function.

```go
package main

import (
  "fmt"
  "github.com/dannyvankooten/vat"
)

func main() {
  c, err := vat.GetCountryRates("NL")
  r, err := c.GetRate("standard")

  fmt.Printf("Standard VAT rate for NL is %.2f", r)
  // Output: Standard VAT rate for NL is 21.00
}
```

## License

MIT licensed. See the LICENSE file for details.
