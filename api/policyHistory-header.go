package api

// policyHistory stores data about COVID policies for all countries between two dates.
//
// Functionality: get, req
type policyHistory struct {
	Scale     map[string]map[string]int `json:"scale"`
	Countries []string                  `json:"countries"`
	Data      map[string]map[string]struct {
		DataValue            string  `json:"date_value"`
		CountryCode          string  `json:"country_code"`
		Confirmed            int     `json:"confirmed"`
		Deaths               int     `json:"deaths"`
		StringencyActual     float64 `json:"stringency_actual"`
		Stringency           float64 `json:"stringency"`
		StringencyLegacy     float64 `json:"stringency_legacy"`
		StringencyLegacyDisp float64 `json:"stringency_legacy_disp"`
	} `json:"data"`
}
