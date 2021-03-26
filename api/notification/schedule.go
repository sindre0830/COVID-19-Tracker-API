package notification

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/api/cases"
	"main/api/policy"
	"main/dict"
	"net/http"
	"net/http/httptest"
	"time"
)

// cachedCases keeps Cases data from last call if the trigger is 'ON_CHANGE'.
var cachedCases = map[string]cases.Cases {}

// cachedPolicies keeps Policy data from last call if the trigger is 'ON_CHANGE'.
var cachedPolicies = map[string]policy.Policy {}

// secret
var secret = []byte{43, 123, 65, 232, 4, 42, 35, 234, 21, 122, 214}

// Schedule iterates every second and calls webhooks on timeout.
//
// Source: https://gobyexample.com/tickers
func Schedule() {
	//ticks every second.
	var ticker = time.NewTicker(time.Second * 1)
	var i int64
	for ;; <- ticker.C {
		i++
		//check every webhook and call on timeout
		for _, elem := range Notifications {
			if i % elem.Timeout == 0 {
				go callURL(elem)
			}
		}
	}
}

// callURL sends data to webhooks depending on their trigger.
func callURL(notification Notification) {
	//check which field is requested and get data
	var output []byte 
	if notification.Information == "stringency" {
		url := dict.URL + "/corona/v1/policy/" + notification.Country
		var policy policy.Policy
		//create new GET request and branch if an error occurred
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf(
				"%v {\n\tError when creating HTTP request to Policy.Handler().\n\tRaw error: %v\n}\n", 
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
			return
		}
		//call the policy handler and branch if the status code is not OK
		//this stops timed out request being sent to the webhook
		recorder := httptest.NewRecorder()
		policy.Handler(recorder, req)
		if recorder.Result().StatusCode != http.StatusOK {
			fmt.Printf(
				"%v {\n\tError when creating HTTP request to Policy.Handler().\n\tStatus code: %v\n}\n", 
				time.Now().Format("2006-01-02 15:04:05"), recorder.Result().StatusCode,
			)
			return
		}
		//check if the trigger is 'ON_CHANGE' and if it already exists in memory
		//and check if the values has changed. if it has, update the map with the new data
		if cachedPolicy, ok := cachedPolicies[notification.ID]; ok {
			if cachedPolicy.Stringency != policy.Stringency {
				cachedPolicies[notification.ID] = policy
			} else if notification.Trigger == "ON_CHANGE" {
				return
			}
		} else if notification.Trigger == "ON_CHANGE" {
			cachedPolicies[notification.ID] = policy
		}
		//convert from structure to bytes and branch if an error occurred
		output, err = json.Marshal(policy)
		if err != nil {
			fmt.Printf(
				"%v {\n\tError when parsing Policy structure.\n\tRaw error: %v\n}\n", 
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
			return
		}
	} else {
		url := dict.URL + "/corona/v1/country/" + notification.Country
		var cases cases.Cases
		//create new GET request and branch if an error occurred
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf(
				"%v {\n\tError when creating HTTP request to Cases.Handler().\n\tRaw error: %v\n}\n", 
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
			return
		}
		//call the cases handler and branch if the status code is not OK
		//this stops timed out request being sent to the webhook
		recorder := httptest.NewRecorder()
		cases.Handler(recorder, req)
		if recorder.Result().StatusCode != http.StatusOK {
			fmt.Printf(
				"%v {\n\tError when creating HTTP request to Cases.Handler().\n\tStatus code: %v\n}\n", 
				time.Now().Format("2006-01-02 15:04:05"), recorder.Result().StatusCode,
			)
			return
		}
		//check if the trigger is 'ON_CHANGE' and if it already exists in memory
		//and check if the values has changed. if it has, update the map with the new data
		if cachedCase, ok := cachedCases[notification.ID]; ok {
			if cachedCase.Confirmed != cases.Confirmed && cachedCase.Recovered != cases.Recovered {
				cachedCases[notification.ID] = cases
			} else if notification.Trigger == "ON_CHANGE" {
				return
			}
		} else if notification.Trigger == "ON_CHANGE" {
			cachedCases[notification.ID] = cases
		}
		//convert from structure to bytes and branch if an error occurred
		output, err = json.Marshal(cases)
		if err != nil {
			fmt.Printf(
				"%v {\n\tError when parsing Cases structure.\n\tRaw error: %v\n}\n", 
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
			return
		}
	}
	//create new POST request and branch if an error occurred
	req, err := http.NewRequest(http.MethodPost, notification.URL, bytes.NewBuffer(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when creating new POST request.\n\tRaw error: %v\n}\n", 
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		return
	}
	//hash structure and branch if an error occurred
	mac := hmac.New(sha256.New, secret)
	_, err = mac.Write([]byte(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when hashing content before POST request.\n\tRaw error: %v\n}\n", 
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		return
	}
	//convert hashed structure to string and add to header
	req.Header.Add("Signature", hex.EncodeToString(mac.Sum(nil)))
	//update header to JSON
	req.Header.Set("Content-Type", "application/json")
	//send request to client and branch if an error occured
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when sending HTTP content to webhook.\n\tRaw error: %v\n}\n", 
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		return
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		fmt.Printf(
			"%v {\n\tWebhook URL is not valid. Deleting webhook...\n\tStatus code: %v\n}\n", 
			time.Now().Format("2006-01-02 15:04:05"), res.StatusCode,
		)
		DB.Delete(notification.ID)
		return
	}
}
