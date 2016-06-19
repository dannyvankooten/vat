package rates

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// CountryRates holds the various differing VAT rate periods for a given country
type CountryRates struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"country_code"`
	Periods     []struct {
		EffectiveFrom string `json:"effective_from"`
		Rates         map[string]float32
	}
}

type apiResponse struct {
	Details string
	Version string
	Rates   []CountryRates
}

var countryRates []CountryRates

// ErrServiceUnavailable will be returned when jsonvat.com is unreachable
var ErrServiceUnavailable = errors.New("VAT Rates API is unavailable")

// ErrInvalidCountryCode will be returned when calling Country with an invalid country code
var ErrInvalidCountryCode = errors.New("Unknown country code.")

// ErrInvalidRateLevel will be returned when getting wrong rate level
var ErrInvalidRateLevel = errors.New("Unknown rate level")

// Rate returns the currently active rate
func (r *CountryRates) Rate(level string) (float32, error) {
	now := time.Now()

	// return the first rate where EffectiveFrom is in the past.
	for _, p := range r.Periods {
		date, _ := time.Parse("0000-01-01", p.EffectiveFrom)
		if now.After(date) {
			return p.Rates[level], nil
		}
	}

	return 0.00, ErrInvalidRateLevel
}

// Country gets the CountryRates struct for a country by its ISO-3166-1-alpha2 country code.
func Country(c string) (CountryRates, error) {
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
