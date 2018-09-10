package vat

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
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

var mutex = &sync.Mutex{} // protect countriesRates
var countriesRates []CountryRates

// ErrInvalidCountryCode will be returned when calling GetCountryRates with an invalid country code
var ErrInvalidCountryCode = errors.New("vat: unknown country code")

// ErrInvalidRateLevel will be returned when getting wrong rate level
var ErrInvalidRateLevel = errors.New("vat: unknown rate level")

// GetRateOn returns the effective VAT rate on a given date
func (cr *CountryRates) GetRateOn(t time.Time, level string) (float32, error) {
	var activePeriod RatePeriod

	// find active period for the given time
	for _, p := range cr.Periods {
		if t.After(p.EffectiveFrom) && (activePeriod.EffectiveFrom.IsZero() || p.EffectiveFrom.After(activePeriod.EffectiveFrom)) {
			activePeriod = p
		}
	}

	activeRate, ok := activePeriod.Rates[level]
	if !ok {
		return 0.00, ErrInvalidRateLevel
	}

	return activeRate, nil
}

// GetRate returns the currently active rate
func (cr *CountryRates) GetRate(level string) (float32, error) {
	now := time.Now()
	return cr.GetRateOn(now, level)
}

// GetCountryRates gets the CountryRates struct for a country by its ISO-3166-1-alpha2 country code.
func GetCountryRates(countryCode string) (CountryRates, error) {
	var rate CountryRates
	rates, err := GetRates()

	if err != nil {
		return rate, err
	}

	for _, r := range rates {
		if r.CountryCode == countryCode {
			return r, nil
		}
	}

	return rate, ErrInvalidCountryCode
}

// GetRates returns the in-memory VAT rates
func GetRates() ([]CountryRates, error) {
	var err error

	mutex.Lock()
	if countriesRates == nil {
		countriesRates, err = FetchRates()
	}
	mutex.Unlock()

	return countriesRates, err
}

// FetchRates fetches the latest VAT rates from jsonvat.com and updates the in-memory rates
func FetchRates() ([]CountryRates, error) {

	client := http.Client{
		Timeout: (time.Duration(ServiceTimeout) * time.Second),
	}
	r, err := client.Get("https://jsonvat.com/")
	if err != nil {
		return nil, err
	}

	apiResponse := &struct {
		Details string
		Version string
		Rates   []CountryRates
	}{}

	err = json.NewDecoder(r.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	// convert EffectiveFrom to a proper time.Time for each rate period
	for idx1, cr := range apiResponse.Rates {
		for idx2, crp := range cr.Periods {
			crp.EffectiveFromStr = strings.Replace(crp.EffectiveFromStr, "0000-", "2000-", 1)
			apiResponse.Rates[idx1].Periods[idx2].EffectiveFrom, _ = time.Parse("2006-01-02", crp.EffectiveFromStr)
		}
	}

	countriesRates = apiResponse.Rates
	return countriesRates, err
}
