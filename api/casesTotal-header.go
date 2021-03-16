package api

type casesTotal struct {
	All struct {
		Confirmed         int          `json:"confirmed"`
		Recovered         int          `json:"recovered"`
		Deaths            int          `json:"deaths"`
		Country           string       `json:"country"`
		Population        int          `json:"population"`
		SqKmArea          int          `json:"sq_km_area"`
		LifeExpectancy    *interface{} `json:"life_expectancy"`
		ElevationInMeters int          `json:"elevation_in_meters"`
		Continent         string       `json:"continent"`
		Abbreviation      string       `json:"abbreviation"`
		Location          string       `json:"location"`
		Iso               int          `json:"iso"`
		CapitalCity       string       `json:"capital_city"`
		Lat               string       `json:"lat"`
		Long              string       `json:"long"`
		Updated           string       `json:"updated"`
	} `json:"all"`
}
