package api

type cases struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int32   `json:"confirmed"`
	Recovered            int32   `json:"recovered"`
	PopulationPercentage float32 `json:"population_percentage"`
}
