package notification

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/api/cases"
	"main/api/policy"
	"net/http"
	"net/http/httptest"
	"time"
)

var Ticker *time.Ticker
var i int64
var cachedCases = map[string]cases.Cases {}
var cachedPolicies = map[string]policy.Policy {}
// Initialize signature (via init())
var SignatureKey = "X-SIGNATURE"
//var Mac hash.Hash
var Secret []byte

func Schedule() {
	done := make(chan bool)
	for {
		select {
		case <-done:
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

func callURL(notification Notification) {
	var output []byte 
	if notification.Information == "stringency" {
		url := "http://localhost:8080/corona/v1/policy/" + notification.Country
		var policy policy.Policy
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request in callURL")
			return
		}
		recorder := httptest.NewRecorder()
		policy.Handler(recorder, req)
		if cachedPolicy, ok := cachedPolicies[notification.ID]; ok {
			if cachedPolicy.Stringency != policy.Stringency {
				cachedPolicies[notification.ID] = policy
			} else if notification.Trigger == "ON_CHANGE" {
				return
			}
		} else if notification.Trigger == "ON_CHANGE" {
			cachedPolicies[notification.ID] = policy
		}
		output, err = json.Marshal(policy)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		url := "http://localhost:8080/corona/v1/country/" + notification.Country
		var cases cases.Cases
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating HTTP request in callURL")
			return
		}
		recorder := httptest.NewRecorder()
		cases.Handler(recorder, req)
		if cachedCase, ok := cachedCases[notification.ID]; ok {
			if cachedCase.Confirmed != cases.Confirmed && cachedCase.Recovered != cases.Recovered {
				cachedCases[notification.ID] = cases
			} else if notification.Trigger == "ON_CHANGE" {
				return
			}
		} else if notification.Trigger == "ON_CHANGE" {
			cachedCases[notification.ID] = cases
		}
		output, err = json.Marshal(cases)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	req, err := http.NewRequest(http.MethodPost, notification.URL, bytes.NewBuffer(output))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Hash content
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(output))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Convert to string & add to header
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in HTTP request: " + err.Error())
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Something is wrong with invocation response: " + err.Error())
	}
}
