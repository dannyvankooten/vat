// Package vat provides VAT number verification for Golang.
package vat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type viesResponse struct {
	CountryCode string
	VATnumber   string
	RequestDate time.Time
	Valid       bool
	Name        string
	Address     string
}

const serviceURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"
const envelope = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://schemas.conversesolutions.com/xsd/dmticta/v1">
<soapenv:Header/>
<soapenv:Body>
  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
    <countryCode>{{.countryCode}}</countryCode>
    <vatNumber>{{.vatNumber}}</vatNumber>
  </checkVat>
</soapenv:Body>
</soapenv:Envelope>
`

var (
	ErrInvalidVATNumber   = errors.New("VAT number is invalid.")
	ErrServiceUnreachable = errors.New("Validation service is offline.")
)

// Validate validates a VAT number by format and existence.
//
// The existence check uses the VIES VAT validation SOAP API and will only run when format validation passes.
func Validate(n string) (bool, error) {
	format, err := ValidateFormat(n)
	existence := false

	if format {
		existence, err = ValidateExistence(n)
	}

	return (format && existence), err
}

// ValidateFormat validates a VAT number by its format.
func ValidateFormat(n string) (bool, error) {
	patterns := map[string]string{
		"AT": "U[A-Z\\d]{8}",
		"BE": "(0\\d{9}|\\d{10})",
		"BG": "\\d{9,10}",
		"CY": "\\d{8}[A-Z]",
		"CZ": "\\d{8,10}",
		"DE": "\\d{9}",
		"DK": "(\\d{2} ?){3}\\d{2}",
		"EE": "\\d{9}",
		"EL": "\\d{9}",
		"ES": "[A-Z]\\d{7}[A-Z]|\\d{8}[A-Z]|[A-Z]\\d{8}",
		"FI": "\\d{8}",
		"FR": "([A-Z]{2}|\\d{2})\\d{9}",
		"GB": "\\d{9}|\\d{12}|(GD|HA)\\d{3}",
		"HR": "\\d{11}",
		"HU": "\\d{8}",
		"IE": "[A-Z\\d]{8}|[A-Z\\d]{9}",
		"IT": "\\d{11}",
		"LT": "(\\d{9}|\\d{12})",
		"LU": "\\d{8}",
		"LV": "\\d{11}",
		"MT": "\\d{8}",
		"NL": "\\d{9}B\\d{2}",
		"PL": "\\d{10}",
		"PT": "\\d{9}",
		"RO": "\\d{2,10}",
		"SE": "\\d{12}",
		"SI": "\\d{8}",
		"SK": "\\d{10}",
	}

	if len(n) < 3 {
		return false, nil
	}

	n = strings.ToUpper(n)
	pattern, ok := patterns[n[0:2]]
	if !ok {
		return false, nil
	}

	matched, err := regexp.MatchString(pattern, n[2:])
	return matched, err
}

// ValidateExistence validates a VAT number by its existence using the VIES VAT API (using SOAP)
func ValidateExistence(n string) (bool, error) {
	r, err := checkVAT(n)
	return r.Valid, err
}

// Check returns *VATresponse for vat number
func checkVAT(vatNumber string) (*viesResponse, error) {
	if len(vatNumber) < 3 {
		return nil, ErrInvalidVATNumber
	}

	e, err := getEnvelope(vatNumber)
	if err != nil {
		return nil, err
	}
	eb := bytes.NewBufferString(e)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Post(serviceURL, "text/xml;charset=UTF-8", eb)
	if err != nil {
		return nil, ErrServiceUnreachable
	}
	defer res.Body.Close()

	// TODO: Use reader XML decoder
	xmlRes, err := ioutil.ReadAll(res.Body)

	// check if response contains "INVALID_INPUT" string
	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return nil, ErrInvalidVATNumber
	}

	var rd struct {
		XMLName xml.Name `xml:"Envelope"`
		Soap    struct {
			XMLName xml.Name `xml:"Body"`
			Soap    struct {
				XMLName     xml.Name `xml:"checkVatResponse"`
				CountryCode string   `xml:"countryCode"`
				VATnumber   string   `xml:"vatNumber"`
				RequestDate string   `xml:"requestDate"` // 2015-03-06+01:00
				Valid       bool     `xml:"valid"`
				Name        string   `xml:"name"`
				Address     string   `xml:"address"`
			}
		}
	}
	if err := xml.Unmarshal(xmlRes, &rd); err != nil {
		return nil, err
	}

	pDate, err := time.Parse("2006-01-02-07:00", rd.Soap.Soap.RequestDate)
	if err != nil {
		return nil, err
	}

	r := &viesResponse{
		CountryCode: rd.Soap.Soap.CountryCode,
		VATnumber:   rd.Soap.Soap.VATnumber,
		RequestDate: pDate,
		Valid:       rd.Soap.Soap.Valid,
		Name:        rd.Soap.Soap.Name,
		Address:     rd.Soap.Soap.Address,
	}

	return r, nil
}

// getEnvelope parses envelope template
func getEnvelope(vatNumber string) (string, error) {
	t, err := template.New("envelope").Parse(envelope)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	if err := t.Execute(&result, map[string]string{
		"countryCode": strings.ToUpper(vatNumber[0:2]),
		"vatNumber":   vatNumber[2:],
	}); err != nil {
		return "", err
	}
	return result.String(), nil
}
