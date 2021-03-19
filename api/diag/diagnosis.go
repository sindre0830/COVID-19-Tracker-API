package diag

import (
	"encoding/json"
	"fmt"
	"main/debug"
	"net/http"
	"time"
)

// StartTime is a variable declared at the start of the program to calculate uptime.
var StartTime time.Time
// Diagnosis structure keeps version, uptime and status codes on used API's.
//
// Functionality: Handler, get, req, update, getUptime
type Diagnosis struct {
	Mmediagroupapi   int    `json:"mmediagroupapi"`
	Covidtrackerapi  int    `json:"covidtrackerapi"`
	Restcountriesapi int    `json:"restcountriesapi"`
	Registered       int    `json:"registered"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}

func (diagnosis *Diagnosis) Handler(w http.ResponseWriter, r *http.Request) {
	//status codes for used REST services
	status, err := diagnosis.get()
	//branch if there is an error
	if err != nil {
		debug.ErrorMessag.Update(
			status, 
			"HandlerDiagnosis() -> Getting status codes from used REST services",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	//amount of registerd webhooks (not implemented yet)
	diagnosis.Registered = 0
	diagnosis.Version = "v1"
	//get uptime
	diagnosis.Uptime = fmt.Sprintf("%f", diagnosis.getUptime())
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(diagnosis)
	//branch if something went wrong with output
	if err != nil {
		debug.ErrorMessag.Update(
			http.StatusInternalServerError, 
			"Diagnosis.Handler() -> Sending output to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessag.Print(w)
	}
}

func (diagnosis *Diagnosis) get() (int, error) {
	mmediagroupStatus, err := diagnosis.req("https://covid-api.mmediagroup.fr/v1/cases")
	if err != nil {
		return mmediagroupStatus, err
	}
	covidtrackerStatus, err := diagnosis.req("https://restcountries.eu/rest/v2/all")
	if err != nil {
		return covidtrackerStatus, err
	}
	restcountriesStatus, err := diagnosis.req("https://restcountries.eu/rest/v2/all")
	if err != nil {
		return restcountriesStatus, err
	}
	diagnosis.update(mmediagroupStatus, covidtrackerStatus, restcountriesStatus)
	return http.StatusOK, nil
}

func (diagnosis *Diagnosis) req(url string) (int, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return rsp.StatusCode, err
	}
	return rsp.StatusCode, nil
}

func (diagnosis *Diagnosis) update(mmediagroupStatus int, covidtrackerStatus int, restcountriesStatus int) {
	diagnosis.Mmediagroupapi = mmediagroupStatus
	diagnosis.Covidtrackerapi = covidtrackerStatus
	diagnosis.Restcountriesapi = restcountriesStatus
}
// getUptime calculates uptime based on start time and current time.
func (diagnosis *Diagnosis) getUptime() float64 {
	return time.Since(StartTime).Seconds()
}
