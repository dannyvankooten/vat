vat
===

Package for validating VAT numbers in Golang.

```go
// validates format + existence
validity := vat.Validate("NL123456789B01")

// validate format
validity := vat.ValidateFormat("NL123456789B01")

// validate existence
validity := vat.ValidateExistence("NL123456789B01")
```

### License

MIT licensed. See the LICENSE file for details.
