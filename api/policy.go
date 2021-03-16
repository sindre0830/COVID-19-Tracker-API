package api

import "net/http"

// Handler will handle http request for COVID policies.
func (policy *Policy) Handler(w http.ResponseWriter, r *http.Request) {

}
// get will update Cases based on input.
func (policy *Policy) get(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// getTotal will get all available COVID policies.
func (policy *Policy) getTotal(country string) (int, error) {
	return http.StatusOK, nil
}
// getHistory will get COVID policies between two dates.
func (policy *Policy) getHistory(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// update sets new data in cases.
func (policy *Policy) update(country string, scope string, stringency float64, trend float64, Update string) {

}
