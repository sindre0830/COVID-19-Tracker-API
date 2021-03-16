package api

// countryCode stores country name details.
//
// countryCode: get, req
type countryCode []struct {
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
}
