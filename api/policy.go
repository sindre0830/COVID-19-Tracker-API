package api

import "net/http"

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
