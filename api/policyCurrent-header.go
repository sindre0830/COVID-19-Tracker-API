package api

// policyCurrent stores data about current COVID policies based on a country.
//
// Functionality: get, req
type policyCurrent struct {
	Policyactions []struct {
		PolicyTypeCode          string      `json:"policy_type_code"`
		PolicyTypeDisplay       string      `json:"policy_type_display"`
		Policyvalue             string      `json:"policyvalue"`
		PolicyvalueActual       int         `json:"policyvalue_actual"`
		Flagged                 bool        `json:"flagged"`
		IsGeneral               bool        `json:"is_general"`
		Notes                   interface{} `json:"notes"`
		FlagValueDisplayField   string      `json:"flag_value_display_field,omitempty"`
		PolicyValueDisplayField string      `json:"policy_value_display_field"`
	} `json:"policyActions"`
	Stringencydata struct {
		DateValue        string  `json:"date_value"`
		CountryCode      string  `json:"country_code"`
		Confirmed        int     `json:"confirmed"`
		Deaths           int     `json:"deaths"`
		StringencyActual float64 `json:"stringency_actual"`
		Stringency       float64 `json:"stringency"`
	} `json:"stringencyData"`
}
