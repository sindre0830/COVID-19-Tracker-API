package api

// Cases stores data about COVID cases based on user input.
//
// Functionality: Handler, get, getTotal, getHistory, update
type Cases struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage float64 `json:"population_percentage"`
}
