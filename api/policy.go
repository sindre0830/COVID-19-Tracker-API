package api

import "net/http"

// Policy stores data about COVID policies based on user input.
//
// Functionality: Handler, get, getTotal, getHistory, update
type Policy struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend	   float64 `json:"trend"`
	Update 	   string  `json:"update"`
}
// Handler will handle http request for COVID policies.
func (policy *Policy) Handler(w http.ResponseWriter, r *http.Request) {

}
// get will update Policy based on input.
func (policy *Policy) get(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// getCurrent will get current available COVID policies.
func (policy *Policy) getCurrent(country string) (int, error) {
	return http.StatusOK, nil
}
// getHistory will get COVID policies between two dates.
func (policy *Policy) getHistory(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// update sets new data in cases.
func (policy *Policy) update(country string, scope string, stringency float64, trend float64, Update string) {

}
