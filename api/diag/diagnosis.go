package diag

import "net/http"

// Diagnosis structure keeps version, uptime and status codes on used API's.
type Diagnosis struct {
	Mmediagroupapi   int    `json:"mmediagroupapi"`
	Covidtrackerapi  int    `json:"covidtrackerapi"`
	Restcountriesapi int    `json:"restcountriesapi"`
	Registered       int    `json:"registered"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}

func (diagnosis *Diagnosis) Handler(w http.ResponseWriter, r *http.Request) {
	
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
