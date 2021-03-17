package api

// PolicyCurrent stores data about current COVID policies based on a country.
//
// Functionality: Get, req, isEmpty
type PolicyCurrent struct {
	Policyactions []struct {
		PolicyTypeCode          string      `json:"policy_type_code"`
		PolicyTypeDisplay       string      `json:"policy_type_display"`
		Policyvalue             interface{}      `json:"policyvalue"`
		PolicyvalueActual       *interface{}         `json:"policyvalue_actual"`
		Flagged                 interface{}        `json:"flagged"`
		IsGeneral               *interface{}        `json:"is_general"`
		Notes                   interface{} `json:"notes"`
		FlagValueDisplayField   string      `json:"flag_value_display_field"`
		PolicyValueDisplayField string      `json:"policy_value_display_field"`
	} `json:"policyActions"`
	Stringencydata struct {
		DateValue        *interface{}  `json:"date_value"`
		CountryCode      *interface{}  `json:"country_code"`
		Confirmed        *interface{}     `json:"confirmed"`
		Deaths           *interface{}     `json:"deaths"`
		StringencyActual *interface{} `json:"stringency_actual"`
		Stringency       *interface{} `json:"stringency"`
		Msg              *interface{} `json:"msg"`
	} `json:"stringencyData"`
}
