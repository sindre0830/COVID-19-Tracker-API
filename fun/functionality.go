package fun

import (
	"errors"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// LimitDecimals limits decimals to two.
func LimitDecimals(number float64) float64 {
	return (math.Round(number * 100) / 100)
}
// ParseURL parses URL to country and possible scope.
func ParseURL(inpURL *url.URL) (string, string, int, error) {
	//split URL path by '/'
	arrPath := strings.Split(inpURL.Path, "/")
	//branch if there aren't enough elements in URL and return error
	if len(arrPath) != 5 {
		err := errors.New("url validation: either too many or too few arguments in url path")
		return "", "", http.StatusBadRequest, err
	}
	//set country
	country := arrPath[4]
	//set default scope to nil (total)
	scope := ""
	//get all parameters from URL and branch if an error occurred
	arrParam, err := url.ParseQuery(inpURL.RawQuery)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	//branch if any parameters exist
	if len(arrParam) > 0 {
		//branch if field 'scope' exist otherwise return an error
		if targetParameter, ok := arrParam["scope"]; ok {
			scope = targetParameter[0]
		} else {
			err := errors.New("url validation: wrong parameter")
			return "", "", http.StatusBadRequest, err
		}
	}
	return country, scope, http.StatusOK, nil
}
// ValidateCountry checks if country is empty.
func ValidateCountry(country string) error {
	if country == "" {
		err := errors.New("country validation: empty field")
		return err
	}
	return nil
}
// ConvertCountry converts any inputted name to uphold strict syntax.
func ConvertCountry(country string) string {
	return strings.Title(strings.ToLower(country))
}
// ValidateDates checks all possible date formatting mistakes.
func ValidateDates(dates string) error {
	//declare error variable
	var err error
	//split date by '-' for format checking
	arrDate := strings.Split(dates, "-")
	//branch if date doesn't have correct amount of elements
	if (len(arrDate) != 6) || (len(dates) != 21) {
		err = errors.New("date validation: not enough elements")
		return err
	}
	//branch if start date isn't using correct format (YYYY-MM-DD)
	if (len(arrDate[0]) != 4) || (len(arrDate[1]) != 2) || (len(arrDate[2]) != 2) {
		err = errors.New("date validation: start date doesn't follow required date format")
		return err
	}
	//branch if end date isn't using correct format (YYYY-MM-DD)
	if (len(arrDate[3]) != 4) || (len(arrDate[4]) != 2) || (len(arrDate[5]) != 2) {
		err = errors.New("date validation: end date doesn't follow required date format")
		return err
	}
	//branch if date elements aren't integers or larger than 0. 'hehe-01-00' == false
	for _, elemDate := range arrDate {
		elemDateNum, err := strconv.Atoi(elemDate)
		if err != nil || elemDateNum < 1 {
			err = errors.New("date validation: wrong type, should be numbers that are larger than 0")
			return err
		}
	}
	//branch if end date isn't larger or equal to start date
	if dates[:10] > dates[11:] {
		err = errors.New("date validation: start date is larger than end date")
	}
	return err
}
