package api

import "net/http"

// get will update countryCode based on input.
func (code *countryCode) get(country string, startDate string, endDate string) (int, error) {
	return http.StatusOK, nil
}
// req will request from API based on URL.
func (code *countryCode) req(url string) (int, error) {
	return http.StatusOK, nil
}
