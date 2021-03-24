package cases

import (
	"encoding/json"
	"main/api/countryinfo"
	"main/debug"
	"main/dict"
	"main/fun"
	"net/http"
)

// Cases structure stores information about COIVD cases for a country.
//
// Functionality: Handler, get, getTotal, getHistory
type Cases struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage float64 `json:"population_percentage"`
}

// Handler will handle http request for REST service.
func (cases Cases) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	country, scope, status, err := fun.ParseURL(r.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Cases.Handler() -> ParseURL() -> Parsing URL",
			err.Error(),
			"URL format. Expected format: '.../country?scope=start_at-end_at'. Example: '.../norway?scope=2020-01-20-2021-02-01'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//get alphacode and country name by RestCountry and branch if an error occurred
	var countryNameDetails countryinfo.CountryNameDetails
	status, err = countryNameDetails.Get(country)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Cases.Handler() -> CountryNameDetails.Get() -> Getting country details",
			err.Error(),
			"Country format. Expected format: '.../country'. Example: '.../norway'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//branch if countrycode is an edgecase and set custom country name as defined in the dictionary, 
	//otherwise use RestCountry country name
	if countryName, ok := dict.CountryEdgeCases[countryNameDetails[0].Alpha3Code]; ok {
		//set edgecase and branch if it is marked as invalid
		country = countryName
		err := fun.ValidateCountry(country)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusNotFound,
				"Cases.Handler() -> ValidatingCountry() -> Checking if inputed country is valid",
				err.Error(),
				"Country format. Expected format: '.../country'. Example: '.../norway'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		country = countryNameDetails[0].Name
	}
	//set default start- and end date variables (total) and check if user inputted scope
	startDate := ""
	endDate := ""
	if len(scope) > 0 {
		//validate scope and branch if an error occurred
		err = fun.ValidateDates(scope)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusBadRequest, 
				"Cases.Handler() -> ValidateDates() -> Checking if inputed dates are valid",
				err.Error(),
				"Date format. Expected format: '...?scope=start_at-end_at'. Example: '...?scope=2020-01-20-2021-02-01'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		startDate = scope[:10]
		endDate = scope[11:]
	}
	//get data based on country and scope and branch if an error occured
	status, err = cases.get(country, startDate, endDate)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Cases.Handler() -> Cases.get() -> Getting covid cases data",
			err.Error(),
			"Country format. Either country doesn't exist in our database or it's mistyped",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//send output to user and branch if an error occured
	err = json.NewEncoder(w).Encode(cases)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Cases.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

// get will get data for structure.
func (cases Cases) get(country string, startDate string, endDate string) (int, error) {
	//branch if scope parameter is used
	if startDate == "" {
		//get all available data and branch if an error occurred
		status, err := cases.getTotal(country)
		if err != nil {
			return status, err
		}
	} else {
		//get data between within scope and branch if an error occurred
		status, err := cases.getHistory(country, startDate, endDate)
		if err != nil {
			return status, err
		}
	}
	return http.StatusOK, nil
}

// getTotal will get all available data.
func (cases *Cases) getTotal(country string) (int, error) {
	var casesTotal CasesTotal
	//get total cases and branch if an error occurred
	status, err := casesTotal.Get(country)
	if err != nil {
		return status, err
	}
	//set data in structure
	cases.Country = casesTotal.All.Country
	cases.Continent = casesTotal.All.Continent
	cases.Scope = "total"
	cases.Confirmed = casesTotal.All.Confirmed
	cases.Recovered = casesTotal.All.Recovered
	cases.PopulationPercentage = fun.LimitDecimals(float64(casesTotal.All.Confirmed) / float64(casesTotal.All.Population))
	return http.StatusOK, nil
}

// getHistory will get data within scope.
func (cases *Cases) getHistory(country string, startDate string, endDate string) (int, error) {
	var casesHistory CasesHistory
	//get confirmed cases between two dates and branch if an error occurred
	status, err := casesHistory.Get("Confirmed", country)
	if err != nil {
		return status, err
	}
	confirmed := casesHistory.All.Dates[endDate] - casesHistory.All.Dates[startDate]
	//get recovered cases between two dates and branch if an error occurred
	status, err = casesHistory.Get("Recovered", country)
	if err != nil {
		return status, err
	}
	recovered := casesHistory.All.Dates[endDate] - casesHistory.All.Dates[startDate]
	//set data in structure
	cases.Country = casesHistory.All.Country
	cases.Continent = casesHistory.All.Continent
	cases.Scope = startDate + "-" + endDate
	cases.Confirmed = confirmed
	cases.Recovered = recovered
	cases.PopulationPercentage = fun.LimitDecimals(float64(confirmed) / float64(casesHistory.All.Population))
	return http.StatusOK, nil
}
