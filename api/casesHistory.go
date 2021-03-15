package api

type casesHistory struct {
	All struct {
		Country           string 		 `json:"country"`
		Population        int    		 `json:"population"`
		SqKmArea          int    		 `json:"sq_km_area"`
		LifeExpectancy    string 		 `json:"life_expectancy"`
		ElevationInMeters int    		 `json:"elevation_in_meters"`
		Continent         string 		 `json:"continent"`
		Abbreviation      string 		 `json:"abbreviation"`
		Location          string 		 `json:"location"`
		Iso               int    		 `json:"iso"`
		CapitalCity       string 		 `json:"capital_city"`
		Dates 		      map[string]int `json:"dates"`
	} `json:"All"`
}
