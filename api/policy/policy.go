package policy

import (
	"encoding/json"
	"main/api/countryinfo"
	"main/debug"
	"main/dict"
	"main/fun"
	"math"
	"net/http"
	"time"
)

// Policy structure stores information about COVID policies for a country.
//
// Functionality: Handler, get, getcurrent, getHistory, modifyDate
type Policy struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend	   float64 `json:"trend"`
}

// Handler will handle http request for REST service.
func (policy *Policy) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	country, scope, status, err := fun.ParseURL(r.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Policy.Handler() -> Parsing URL",
			err.Error(),
			"URL format. Expected format: '.../country?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../norway?2020-01-20-2021-02-01'",
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
			"Policy.Handler() -> Getting alphacode",
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
		err = fun.ValidateCountry(country)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusNotFound,
				"Policy.Handler() -> ValidatingCountry() -> Checking if inputted country is valid",
				err.Error(),
				"Country format. Expected format: '.../country'. Example: '.../norway'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		country = countryNameDetails[0].Name
	}
	policy.Country = country
	//set default start- and end date variables (total) and check if user inputted scope
	startDate := ""
	endDate := ""
	if len(scope) > 0 {
		//validate scope and branch if an error occurred
		err = fun.ValidateDates(scope)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusBadRequest, 
				"Policy.Handler() -> Checking if inputed dates are valid",
				err.Error(),
				"Date format. Expected format: '...?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '...?scope=2020-01-20-2021-02-01'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		startDate = scope[:10]
		endDate = scope[11:]
	}
	//get data based on country and scope and branch if an error occured
	status, err = policy.get(countryNameDetails[0].Alpha3Code, startDate, endDate)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Policy.Handler() -> Policy.get() -> Getting covid policies data",
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
	err = json.NewEncoder(w).Encode(policy)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Policy.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

// get will get data for structure.
func (policy *Policy) get(country string, startDate string, endDate string) (int, error) {
	//branch if scope parameter is used
	if startDate == "" {
		//get all available data and branch if an error occurred
		status, err := policy.getCurrent(country)
		if err != nil {
			return status, err
		}
	} else {
		//get data between two dates and branch if an error occurred
		status, err := policy.getHistory(country, startDate, endDate)
		if err != nil {
			return status, err
		}
	}
	return http.StatusOK, nil
}

// getCurrent will get current available data.
func (policy *Policy) getCurrent(country string) (int, error) {
	var policyCurrent PolicyCurrent
	//get current time and reduce by 10 days
	currentTime := time.Now()
	currentTime = currentTime.AddDate(0, 0, -10)
	date := currentTime.Format("2006-01-02")
	//get total cases and branch if an error occurred
	status, err := policyCurrent.Get(country, date)
	if err != nil {
		return status, err
	}
	//set stringency to StringencyActual and branch if it isn't filled and fall back to Stringency
	//if neither fields are filled, set stringency to -1 as stated by the assignment
	stringency := policyCurrent.Stringencydata.StringencyActual
	if stringency == 0 {
		stringency = policyCurrent.Stringencydata.Stringency
		if stringency == 0 {
			stringency = -1
		}
	}
	//set data in cases
	policy.Scope = "total"
	policy.Stringency = stringency
	policy.Trend = 0
	return http.StatusOK, nil
}

// getHistory will get data within scope.
func (policy *Policy) getHistory(country string, startDate string, endDate string) (int, error) {
	var policyHistory PolicyHistory
	//check if dates are within the 10 day buffer and modify accordingly, branch if an error occurred
	startDate, err := policy.modifyDate(startDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	endDate, err = policy.modifyDate(endDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	//get data within scope and branch if an error occurred
	status, err := policyHistory.Get(startDate, endDate)
	if err != nil {
		return status, err
	}
	//set stringency to StringencyActual and branch if it isn't filled and fall back to Stringency
	stringencyStart := policyHistory.Data[startDate][country].StringencyActual
	stringencyEnd := policyHistory.Data[endDate][country].StringencyActual
	if stringencyEnd == 0 {
		stringencyStart = policyHistory.Data[startDate][country].Stringency
		stringencyEnd = policyHistory.Data[endDate][country].Stringency
	}
	//set default value of trend and branch if there is valid data
	trend := 0.00
	if stringencyEnd != 0 {
		trend = stringencyEnd - stringencyStart
		trend = fun.LimitDecimals(trend)
	} else {
		//if there is no data, set it to -1
		stringencyEnd = -1
	}
	//set data in structure
	policy.Scope = startDate + "-" + endDate
	policy.Stringency = stringencyEnd
	policy.Trend = trend
	return http.StatusOK, nil
}

// modifyDate decreases the date up to 10 days of current date if it is within the buffer.
// This is because the information from the first 10 days are stated inaccurate by the assignment.
func (policy *Policy) modifyDate(date string) (string, error) {
	//parse date to time format and branch if an error occurred
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	//get current time and subtract inputted date
	currentTime := time.Now()
	diffTime := currentTime.Sub(dateTime)
	//convert the difference to integer of days and branch if it's within the buffer
	diffDays := int(math.Floor(diffTime.Hours() / 24.0))
	if diffDays >= 0 && diffDays < 10 {
		//set amount of days to subtract
		diffDays = 10 - diffDays
		dateTime = dateTime.AddDate(0, 0, -(diffDays))
	}
	date = dateTime.Format("2006-01-02")
	return date, nil
}
