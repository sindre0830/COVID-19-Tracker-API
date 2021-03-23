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

// Policy stores data about COVID policies based on user input.
//
// Functionality: Handler, get, getTotal, getHistory, update, decreaseDate
type Policy struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend	   float64 `json:"trend"`
}
// Handler will handle http request for COVID policies.
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
	//get alphacode and countryname by RestCountry definition and branch if an error occurred
	var countryNameDetails countryinfo.CountryNameDetails
	status, err = countryNameDetails.Get(country)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Cases.Handler() -> Getting alphacode",
			err.Error(),
			"Country format. Expected format: '.../country'. Example: '.../norway'",
		)
		debug.ErrorMessage.Print(w)
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
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"Policy.Handler() -> ValidatingCountry() -> Checking if inputted country is valid",
			err.Error(),
			"Country format. Expected format: '.../country'. Example: '.../norway'",
		)
		debug.ErrorMessage.Print(w)
		return
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
				"Policy.Handler() -> Checking if inputed dates are valid",
				err.Error(),
				"Date format. Expected format: '...?start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '...?2020-01-20-2021-02-01'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		startDate = scope[:10]
		endDate = scope[11:]
	}
	//get data based on country and scope
	status, err = policy.get(country, startDate, endDate)
	//branch if there is an error
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
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(policy)
	//branch if something went wrong with output
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
// get will update Policy based on input.
func (policy *Policy) get(country string, startDate string, endDate string) (int, error) {
	//get country name details and branch if an error occurred
	var countryNameDetails countryinfo.CountryNameDetails
	status, err := countryNameDetails.Get(country)
	if err != nil {
		return status, err
	}
	//branch if scope parameter is used
	if startDate == "" {
		//get all available data and branch if an error occurred
		status, err := policy.getCurrent(countryNameDetails[0].Alpha3Code)
		if err != nil {
			return status, err
		}
	} else {
		//get data between two dates and branch if an error occurred
		status, err := policy.getHistory(countryNameDetails[0].Alpha3Code, startDate, endDate)
		if err != nil {
			return status, err
		}
	}
	policy.Country = country
	return http.StatusOK, nil
}
// getCurrent will get current available COVID policies.
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
	policy.update(
		"total", 
		stringency, 
		0,
	)
	return http.StatusOK, nil
}
// getHistory will get COVID policies between two dates.
func (policy *Policy) getHistory(country string, startDate string, endDate string) (int, error) {
	var policyHistory PolicyHistory
	//increases both dates by 10 days and branch if an error occurred
	increasedStartDate, err := policy.modifyDate(startDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	increasedEndDate, err := policy.modifyDate(endDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	//get total cases and branch if an error occurred
	status, err := policyHistory.Get(increasedStartDate, increasedEndDate)
	if err != nil {
		return status, err
	}
	//set stringency to StringencyActual and branch if it isn't filled and fall back to Stringency
	stringencyStart := policyHistory.Data[increasedStartDate][country].StringencyActual
	stringencyEnd := policyHistory.Data[increasedEndDate][country].StringencyActual
	if stringencyEnd == 0 {
		stringencyStart = policyHistory.Data[increasedStartDate][country].Stringency
		stringencyEnd = policyHistory.Data[increasedEndDate][country].Stringency
	}
	//set default value of trend and branch if there is valid data
	trend := 0.00
	if stringencyEnd != 0 {
		trend = stringencyEnd - stringencyStart
		trend = fun.LimitDecimals(trend)
	} else {
		stringencyEnd = -1
	}
	//set data in cases
	policy.update(
		startDate + "-" + endDate, 
		stringencyEnd, 
		trend,
	)
	return http.StatusOK, nil
}
// update sets new data in cases.
func (policy *Policy) update(scope string, stringency float64, trend float64) {
	policy.Scope = scope
	policy.Stringency = stringency
	policy.Trend = trend
}
// modifyDate increases the date by 10 days of current date if it is within the buffer.
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
