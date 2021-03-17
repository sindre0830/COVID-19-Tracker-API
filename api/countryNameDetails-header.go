package api

// CountryNameDetails stores country name details.
//
// CountryNameDetails: Get, req
type CountryNameDetails []struct {
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
}
