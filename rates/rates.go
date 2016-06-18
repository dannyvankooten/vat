package rates

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Rate struct {
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
	Rates   []Rate
}

// ErrServiceUnavailable will be returned when jsonvat.com is unreachable
var ErrServiceUnavailable = errors.New("VAT Rates API is unavailable")

// ErrInvalidCountryCode will be returned when calling Country with an invalid country code
var ErrInvalidCountryCode = errors.New("Unknown country code.")

// Country gets the Rate struct for a country by its ISO-3166-1-alpha2 code.
func Country(c string) (Rate, error) {
	rates, err := fetch()
	var rate Rate

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

func fetch() ([]Rate, error) {
	r, err := http.Get("https://jsonvat.com/")
	if err != nil {
		return nil, err
	}

	data := new(apiResponse)
	err = json.NewDecoder(r.Body).Decode(&data)

	// TODO: Parse current rate into easily accessible struct field.

	return data.Rates, err
}
