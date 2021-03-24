package fun

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// ParseURL parses URL to country and possible scope.
func ParseURL(inpURL *url.URL) (string, string, int, error) {
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(inpURL.Path, "/")
	if len(arrPath) != 5 {
		err := errors.New("url validation: either too many or too few arguments in url path")
		return "", "", http.StatusBadRequest, err
	}
	country := arrPath[4]
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
