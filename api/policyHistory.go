package api

import "net/http"

// get will update Policy based on input.
func (policyHis *policyHistory) get(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (policyHis *policyHistory) req(url string) (int, error) {
	return http.StatusOK, nil
}
