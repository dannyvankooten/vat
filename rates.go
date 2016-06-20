package vat

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type RatePeriod struct {
	EffectiveFrom string `json:"effective_from"`
	Rates         map[string]float32
}

// CountryRates holds the various differing VAT rate periods for a given country
type CountryRates struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"country_code"`
	Periods     []RatePeriod
}

type apiResponse struct {
	Details string
	Version string
	Rates   []CountryRates
}

var countryRates []CountryRates

// ErrInvalidCountryCode will be returned when calling GetCountryRates with an invalid country code
var ErrInvalidCountryCode = errors.New("Unknown country code.")

// ErrInvalidRateLevel will be returned when getting wrong rate level
var ErrInvalidRateLevel = errors.New("Unknown rate level")

// RateOn returns the effective VAT rate on a given date
func (cr *CountryRates) RateOn(t time.Time, level string) (float32, error) {
	var activePeriod RatePeriod
	var activePeriodDate time.Time

	// find active period for the given time
	for _, p := range cr.Periods {
		periodDate, _ := time.Parse("2006-01-01", p.EffectiveFrom)
		if t.After(periodDate) && periodDate.After(activePeriodDate) {
			activePeriod = p
			activePeriodDate = periodDate
		}
	}

	activeRate, ok := activePeriod.Rates[level]
	if !ok {
		return 0.00, ErrInvalidRateLevel
	}

	return activeRate, nil
}

// Rate returns the currently active rate
func (cr *CountryRates) Rate(level string) (float32, error) {
	now := time.Now()
	return cr.RateOn(now, level)
}

// GetCountryRates gets the CountryRates struct for a country by its ISO-3166-1-alpha2 country code.
func GetCountryRates(c string) (CountryRates, error) {
	rates, err := fetch()
	var rate CountryRates

	if err != nil {
		return rate, err
	}

	for _, r := range rates {
		if r.CountryCode == c {
			return r, nil
		}
	}

	return rate, ErrInvalidCountryCode
}

func fetch() ([]CountryRates, error) {
	if countryRates != nil {
		return countryRates, nil
	}

	r, err := http.Get("https://jsonvat.com/")
	if err != nil {
		return nil, err
	}

	data := new(apiResponse)
	err = json.NewDecoder(r.Body).Decode(&data)
	countryRates = data.Rates

	return countryRates, err
}
