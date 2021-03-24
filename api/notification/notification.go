package notification

import (
	"encoding/json"
	"main/api/cases"
	"main/api/policy"
	"main/debug"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// Notifications stores all webhooks in memory.
//
// This is useful for the schedule loop.
var Notifications = map[string]Notification {}

// Notification structure stores information of a valid webhook.
//
// Functionality: POST, GET, DELETE
type Notification struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Timeout     int64  `json:"timeout"`
	Information string `json:"information"`
	Country     string `json:"country"`
	Trigger     string `json:"trigger"`
}

// POST handles a POST request from the http request.
func (notification *Notification) POST(w http.ResponseWriter, r *http.Request) {
	//read input from client and branch if an error occurred
	var notificationInput NotificationInput
	err := json.NewDecoder(r.Body).Decode(&notificationInput)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"Notification.POST() -> Parsing data from client",
			err.Error(),
			"Wrong JSON format sent.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if URL is valid (very simple check) and branch if an error occurred
	_, err = url.ParseRequestURI(notificationInput.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Notification.POST() -> Checking if URL is valid",
			err.Error(),
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if timeout is valid and return and error if it isn't
	if notificationInput.Timeout < 1 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Notification.POST() -> Checking if timeout value is valid",
			"timeout validation: value less than 1",
			"Not valid timeout value. Example '3600'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if trigger is valid and return and error if it isn't
	notificationInput.Trigger = strings.ToUpper(notificationInput.Trigger)
	if notificationInput.Trigger != "ON_CHANGE" && notificationInput.Trigger != "ON_TIMEOUT" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Notification.POST() -> Checking if trigger value is valid",
			"trigger validation: trigger is not 'ON_CHANGE' or 'ON_TIMEOUT'",
			"Not valid trigger. Example 'ON_TIMEOUT'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if field is valid and return and error if it isn't
	notificationInput.Field = strings.ToLower(notificationInput.Field)
	if notificationInput.Field != "confirmed" && notificationInput.Field != "stringency" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Notification.POST() -> Checking if field is valid",
			"field validation: field is not 'confirmed' or 'stringency'",
			"Not valid field. Example 'stringency'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if country name is valid and branch if an error occurred
	if notificationInput.Field == "confirmed" {
		var cases cases.Cases
		req, err := http.NewRequest("GET", "http://localhost:8080/corona/v1/country/" + notificationInput.Country, nil)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError, 
				"Notification.POST() -> Cases.Handler() -> Checking if country name is valid",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		recorder := httptest.NewRecorder()
		cases.Handler(recorder, req)
		if recorder.Code != http.StatusOK {
			debug.ErrorMessage.Update(
				http.StatusNotFound, 
				"Notification.POST() -> Cases.Handler() -> Checking if country name is valid",
				"country validation: country doesn't exist in our database",
				"Not valid country name. Example 'Norway'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		var policy policy.Policy
		req, err := http.NewRequest("GET", "http://localhost:8080/corona/v1/policy/" + notificationInput.Country, nil)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError, 
				"Notification.POST() -> Policy.Handler() -> Checking if country name is valid",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		recorder := httptest.NewRecorder()
		policy.Handler(recorder, req)
		if recorder.Code != http.StatusOK {
			debug.ErrorMessage.Update(
				http.StatusNotFound, 
				"Notification.POST() -> Policy.Handler() -> Checking if country name is valid",
				"country validation: country doesn't exist in our database",
				"Not valid country name. Example 'Norway'",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
	//set data in structure
	notification.URL = notificationInput.URL
	notification.Timeout = notificationInput.Timeout
	notification.Information = notificationInput.Field
	notification.Country = notificationInput.Country
	notification.Trigger = notificationInput.Trigger
	//add data to database and branch if an error occurred
	err = DB.Add(notification)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Notification.POST() -> Database.Add() -> Adding data to database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//create feedback message to send to client and branch if an error occurred
	var feedback Feedback
	feedback.update(
		http.StatusCreated, 
		"Webhook successfully created for '" + notification.URL + "'",
		notification.ID,
	)
	err = feedback.print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Notification.POST() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}

// GET handles a GET request from the http request.
func (notification *Notification) GET(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 5 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"Notification.GET() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set id and check if it's specified by client
	id := arrPath[4]
	if id == "" {
		//check if there are any webhooks created
		var output []Notification
		if len(Notifications) == 0 {
			http.Error(w, "", http.StatusNoContent)
			return
		}
		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//add webhooks to slice
		for _, element := range Notifications {
			output = append(output, element)
		}
		//send output to user and branch if an error occured
		err := json.NewEncoder(w).Encode(&output)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError, 
				"Notification.GET() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		//check if id is valid and return an error if it isn't
		if _, ok := Notifications[id]; !ok {
			debug.ErrorMessage.Update(
				http.StatusNotFound, 
				"Notification.GET() -> Checking if ID exist",
				"ID validation: can't find ID",
				"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//send output to user and branch if an error occured
		output := Notifications[id]
		err := json.NewEncoder(w).Encode(&output)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError, 
				"Notification.GET() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
}

// DELETE handles a DELETE request from the http request.
func (notification *Notification) DELETE(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 5 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"Notification.DELETE() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//delete webhook by id and branch if an error occurred
	id := arrPath[4]
	err := DB.Delete(id)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"Notification.DELETE() -> Database.Delete() -> Deleting data from database",
			err.Error(),
			"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//create feedback message to send to client and branch if an error occurred
	var feedback Feedback
	feedback.update(
		http.StatusOK, 
		"Webhook successfully deleted",
		id,
	)
	err = feedback.print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Notification.DELETE() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}
