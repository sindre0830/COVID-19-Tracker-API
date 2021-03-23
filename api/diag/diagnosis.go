package diag

import (
	"encoding/json"
	"main/api/notification"
	"main/debug"
	"math"
	"net/http"
	"time"
)

// StartTime defines the time the service started.
var StartTime time.Time

// Diagnosis structure stores information about the REST service.
//
// Functionality: Handler, get, req, getUptime
type Diagnosis struct {
	Mmediagroupapi   int    `json:"mmediagroupapi"`
	Covidtrackerapi  int    `json:"covidtrackerapi"`
	Restcountriesapi int    `json:"restcountriesapi"`
	Registered       int    `json:"registered"`
	Version          string `json:"version"`
	Uptime           int    `json:"uptime"`
}

// Handler will handle http request for REST service.
func (diagnosis *Diagnosis) Handler(w http.ResponseWriter, r *http.Request) {
	//get status codes for used REST services and branch if an error occured
	status, err := diagnosis.get()
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"HandlerDiagnosis() -> Getting status codes from used REST services",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//get the rest of the data for Diagnosis
	diagnosis.Registered = len(notification.Notifications)
	diagnosis.Version = "v1"
	diagnosis.Uptime = diagnosis.getUptime()
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//send output to user and branch if an error occured
	err = json.NewEncoder(w).Encode(diagnosis)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Diagnosis.Handler() -> Sending output to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

// get will get data for structure.
func (diagnosis *Diagnosis) get() (int, error) {
	//request statuscodes for each REST serverice and branch if an error occured
	var err error
	diagnosis.Mmediagroupapi, err = diagnosis.req("https://covid-api.mmediagroup.fr/v1/cases")
	if err != nil {
		return diagnosis.Mmediagroupapi, err
	}
	diagnosis.Covidtrackerapi, err = diagnosis.req("https://restcountries.eu/rest/v2/all")
	if err != nil {
		return diagnosis.Covidtrackerapi, err
	}
	diagnosis.Restcountriesapi, err = diagnosis.req("https://restcountries.eu/rest/v2/all")
	if err != nil {
		return diagnosis.Restcountriesapi, err
	}
	return http.StatusOK, nil
}

// req will request data from API.
func (diagnosis Diagnosis) req(url string) (int, error) {
	rsp, err := http.Get(url)
	//only interested in errors where the status code is not relevent (i.e. not API related)
	if err != nil && rsp.StatusCode == http.StatusOK {
		return rsp.StatusCode, err
	}
	return rsp.StatusCode, nil
}

// getUptime calculates the difference between the start of the service and now.
func (diagnosis Diagnosis) getUptime() int {
	return int(math.Floor(time.Since(StartTime).Seconds()))
}
