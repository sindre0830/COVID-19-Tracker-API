package api

import "net/http"

// get will update policyCurrent based on input.
func (policyNow *policyCurrent) get(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (policyNow *policyCurrent) req(url string) (int, error) {
	return http.StatusOK, nil
}
