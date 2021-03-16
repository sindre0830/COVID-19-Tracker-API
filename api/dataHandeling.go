package api

import (
	"io/ioutil"
	"net/http"
	"time"
)

// requestData gets raw data from API's
func requestData(url string) ([]byte, int, error) {
	//create new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	//branch if there is an error
	if err != nil {
		return nil, req.Response.StatusCode, err
	}
	//timeout after 2 seconds
	apiClient := http.Client{
		Timeout: time.Second * 2,
	}
	//get response
	res, err := apiClient.Do(req)
	//branch if there is an error
	if err != nil {
		return nil, res.StatusCode, err
	}
	//read output
	output, err := ioutil.ReadAll(res.Body)
	//branch if there is an error
	if err != nil {
		return nil, 0, err
	}
	return output, 0, nil
}
