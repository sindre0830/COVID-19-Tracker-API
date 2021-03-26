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

// Ticker ticks every second.
var Ticker *time.Ticker

// cachedCases keeps Cases data from last call if the trigger is 'ON_CHANGE'.
var cachedCases = map[string]cases.Cases {}

// cachedPolicies keeps Policy data from last call if the trigger is 'ON_CHANGE'.
var cachedPolicies = map[string]policy.Policy {}

// SignatureKey
var SignatureKey string

// Secret
var Secret []byte

// Schedule iterates every second and calls webhooks on timeout.
//
// Source: https://gobyexample.com/tickers
func Schedule() {
	done := make(chan bool)
	var i int64
	for {
		select {
			case <- done:
				return
			case <- Ticker.C:
				i++
				for _, elem := range Notifications {
					if i % elem.Timeout == 0 {
						go callURL(elem)
					}
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
			fmt.Printf("\nError when creating HTTP request to Policy.Handler().\nRaw error: %v\n", err.Error())
			return
		}
		//call the policy handler and branch if the status code is not OK
		//this stops timed out request being sent to the webhook
		recorder := httptest.NewRecorder()
		policy.Handler(recorder, req)
		if recorder.Result().StatusCode != http.StatusOK {
			fmt.Printf("\nError when creating HTTP request to Policy.Handler().\nStatus code: %v\n", recorder.Result().StatusCode)
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
			fmt.Printf("\nError when parsing Policy structure.\nRaw error: %v\n", err.Error())
			return
		}
	} else {
		url := dict.URL + "/corona/v1/country/" + notification.Country
		var cases cases.Cases
		//create new GET request and branch if an error occurred
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("\nError when creating HTTP request to Cases.Handler().\nRaw error: %v\n", err.Error())
			return
		}
		//call the cases handler and branch if the status code is not OK
		//this stops timed out request being sent to the webhook
		recorder := httptest.NewRecorder()
		cases.Handler(recorder, req)
		if recorder.Result().StatusCode != http.StatusOK {
			fmt.Printf("\nError when creating HTTP request to Cases.Handler().\nStatus code: %v\n", recorder.Result().StatusCode)
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
			fmt.Printf("\nError when parsing Cases structure.\nRaw error: %v\n", err.Error())
			return
		}
	}
	//create new POST request and branch if an error occurred
	req, err := http.NewRequest(http.MethodPost, notification.URL, bytes.NewBuffer(output))
	if err != nil {
		fmt.Printf("\nError when creating new POST request.\nRaw error: %v\n", err.Error())
		return
	}
	//hash structure and branch if an error occurred
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(output))
	if err != nil {
		fmt.Printf("\nError when hashing content before POST request.\nRaw error: %v\n", err.Error())
		return
	}
	//convert hashed structure to string and add to header
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))
	//update header to JSON
	req.Header.Set("Content-Type", "application/json")
	//send request to client and branch if an error occured
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nError when sending HTTP content to webhook.\nRaw error: %v\n", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		fmt.Printf("\nWebhook URL is not valid. Deleting webhook...\nStatus code: %v\n", res.StatusCode)
		DB.Delete(notification.ID)
		return
	}
}
