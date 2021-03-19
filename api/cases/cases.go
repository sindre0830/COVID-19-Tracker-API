package cases

import (
	"encoding/json"
	"main/api"
	"main/debug"
	"main/dict"
	"main/fun"
	"net/http"
)

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
// Handler will handle http request for COVID cases.
func (cases *Cases) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	country, scope, status, err := fun.ParseURL(r.URL)
	if err != nil {
		debug.UpdateErrorMessage(
			status, 
			"Cases.Handler() -> Parsing URL",
			err.Error(),
			"URL format. Expected format: '.../country?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../norway?2020-01-20-2021-02-01'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//get alphacode and countryname by RestCountry definition and branch if an error occurred
	var countryNameDetails api.CountryNameDetails
	status, err = countryNameDetails.Get(country)
	if err != nil {
		debug.UpdateErrorMessage(
			status, 
			"Cases.Handler() -> Getting alphacode",
			err.Error(),
			"Country format. Expected format: '.../country'. Example: '.../norway'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//branch if countrycode is an edgecase and set custom country name as defined in the dictionary, otherwise use RestCountry country name
	if countryName, ok := dict.Country[countryNameDetails[0].Alpha3Code]; ok {
		country = countryName
	} else {
		country = countryNameDetails[0].Name
	}
	//validate country name and branch if an error occurred
	err = fun.ValidateCountry(country)
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusBadRequest,
			"Cases.Handler() -> ValidatingCountry() -> Checking if inputed country is valid",
			err.Error(),
			"Country format. Expected format: '.../country'. Example: '.../norway'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set default start- and end date variables (total) and check if user inputted scope
	startDate := ""
	endDate := ""
	if len(scope) > 0 {
		//validate scope and branch if an error occurred
		err = fun.ValidateDates(scope)
		if err != nil {
			debug.UpdateErrorMessage(
				http.StatusBadRequest, 
				"Cases.Handler() -> Checking if inputed dates are valid",
				err.Error(),
				"Date format. Expected format: '...?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '...?2020-01-20-2021-02-01'",
			)
			debug.PrintErrorInformation(w)
			return
		}
		startDate = scope[:10]
		endDate = scope[11:]
	}
	//get data based on country and scope
	status, err = cases.get(country, startDate, endDate)
	//branch if there is an error
	if err != nil {
		reason := "Unknown"
		if status == http.StatusBadRequest {
			reason = "Country format. Either country doesn't exist in our database or it's mistyped"
		}
		debug.UpdateErrorMessage(
			status, 
			"Cases.Handler() -> Cases.get() -> Getting covid cases data",
			err.Error(),
			reason,
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(cases)
	//branch if something went wrong with output
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"Cases.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.PrintErrorInformation(w)
	}
}
// get will update Cases based on input.
func (cases *Cases) get(country string, startDate string, endDate string) (int, error) {
	//branch if scope parameter is used
	if startDate == "" {
		//get all available data and branch if an error occurred
		status, err := cases.getTotal(country)
		if err != nil {
			return status, err
		}
	} else {
		//get data between two dates and branch if an error occurred
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
	//set data in cases
	cases.update(casesTotal.All.Country, casesTotal.All.Continent, "total", casesTotal.All.Confirmed, casesTotal.All.Recovered, casesTotal.All.Population)
	return http.StatusOK, nil
}
// getHistory will get data between two dates.
func (cases *Cases) getHistory(country string, startDate string, endDate string) (int, error) {
	var casesHistory CasesHistory
	//get cases between two dates and branch if an error occurred
	confirmed, recovered, status, err := casesHistory.Get(country, startDate, endDate)
	if err != nil {
		return status, err
	}
	//set data in cases
	cases.update(casesHistory.All.Country, casesHistory.All.Continent, startDate + "-" + endDate, confirmed, recovered, casesHistory.All.Population)
	return http.StatusOK, nil
}
// update sets new data in cases.
func (cases *Cases) update(country string, continent string, scope string, confirmed int, recovered int, population int) {
	cases.Country = country
	cases.Continent = continent
	cases.Scope = scope
	cases.Confirmed = confirmed
	cases.Recovered = recovered
	cases.PopulationPercentage = fun.LimitDecimals(float64(confirmed) / float64(population))
}
