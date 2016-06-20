package vat

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// RatePeriod represents a time and the various activate rates at that time.
type RatePeriod struct {
	EffectiveFromStr string `json:"effective_from"`
	EffectiveFrom    time.Time
	Rates            map[string]float32
}

// CountryRates holds the various differing VAT rate periods for a given country
type CountryRates struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"country_code"`
	Periods     []RatePeriod
}

var countriesRates []CountryRates

// ErrInvalidCountryCode will be returned when calling GetCountryRates with an invalid country code
var ErrInvalidCountryCode = errors.New("Unknown country code.")

// ErrInvalidRateLevel will be returned when getting wrong rate level
var ErrInvalidRateLevel = errors.New("Unknown rate level")

// RateOn returns the effective VAT rate on a given date
func (cr *CountryRates) RateOn(t time.Time, level string) (float32, error) {
	var activePeriod RatePeriod

	// find active period for the given time
	for _, p := range cr.Periods {
		if t.After(p.EffectiveFrom) && p.EffectiveFrom.After(activePeriod.EffectiveFrom) {
			activePeriod = p
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
	var rate CountryRates
	rates, err := GetRates()

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

// GetRates returns the in-memory VAT rates
func GetRates() ([]CountryRates, error) {
	var err error

	if countriesRates == nil {
		countriesRates, err = FetchRates()
	}

	return countriesRates, err
}

// FetchRates fetches the latest VAT rates from jsonvat.com and updates the in-memory rates
func FetchRates() ([]CountryRates, error) {

	r, err := http.Get("https://jsonvat.com/")
	if err != nil {
		return nil, err
	}

	type CR CountryRates
	apiResponse := &struct {
		Details string
		Version string
		Rates   []CountryRates
	}{}

	err = json.NewDecoder(r.Body).Decode(&apiResponse)

	// convert EffectiveFrom to a proper time.Time for each rate period
	for idx1, cr := range apiResponse.Rates {
		for idx2, crp := range cr.Periods {
			apiResponse.Rates[idx1].Periods[idx2].EffectiveFrom, _ = time.Parse("2006-01-01", crp.EffectiveFromStr)
		}
	}

	countriesRates = apiResponse.Rates
	return countriesRates, err
}
