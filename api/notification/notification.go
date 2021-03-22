package notification

import (
	"encoding/json"
	"main/api/cases"
	"main/debug"
	"main/fun"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type Notification struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Timeout     int    `json:"timeout"`
	Information string `json:"information"`
	Country     string `json:"country"`
	Trigger     string `json:"trigger"`
}

var notifications = map[string]Notification {}

func (notification *Notification) update(notificationInput NotificationInput) {
	notification.ID = fun.RandString(10)
	notification.URL = notificationInput.URL
	notification.Timeout = notificationInput.Timeout
	notification.Information = notificationInput.Field
	notification.Country = notificationInput.Country
	notification.Trigger = notificationInput.Trigger
	notifications[notification.ID] = *notification
}

func (notification *Notification) POST(w http.ResponseWriter, r *http.Request) {
	var notificationInput NotificationInput
	//read input from client and branch if an error occurred
	err := json.NewDecoder(r.Body).Decode(&notificationInput)
	if err != nil {
		debug.ErrorMessag.Update(
			http.StatusBadRequest, 
			"Notification.POST() -> Parsing data from client",
			err.Error(),
			"Wrong JSON format sent.",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	var cases cases.Cases
	req, err := http.NewRequest("GET", "http://localhost:8080/corona/v1/country/" + notificationInput.Country, nil)
	if err != nil {
		debug.ErrorMessag.Update(
			http.StatusInternalServerError, 
			"Notification.POST() -> Checking if country name is valid",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	recorder := httptest.NewRecorder()
	cases.Handler(recorder, req)
	if recorder.Code != http.StatusOK {
		debug.ErrorMessag.Update(
			http.StatusNotFound, 
			"Notification.POST() -> Checking if country name is valid",
			"country validation: country doesn't exist in our database",
			"Not valid country name. Example 'Norway'",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	_, err = url.ParseRequestURI(notificationInput.URL)
	if err != nil {
		debug.ErrorMessag.Update(
			http.StatusNotFound,
			"Notification.POST() -> Checking if URL is valid",
			err.Error(),
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	notification.update(notificationInput)
	//create feedback message to send to client
	var feedback Feedback
	feedback.update(
		http.StatusCreated, 
		"Webhook successfully created for '" + notification.URL + "'",
		notification.ID,
	)
	err = feedback.print(w)
	if err != nil {
		debug.ErrorMessag.Update(
			http.StatusInternalServerError, 
			"Notification.POST() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessag.Print(w)
		return
	}
}

func (notification *Notification) GET(w http.ResponseWriter, r *http.Request) {
	arrPath := strings.Split(r.URL.Path, "/")
	//branch if there aren't enough elements in URL and return error
	if len(arrPath) != 5 {
		debug.ErrorMessag.Update(
			http.StatusBadRequest, 
			"Notification.GET() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	id := arrPath[4]
	if id == "" {
		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		var output []Notification
		for _, element := range notifications {
			output = append(output, element)
		}
		err := json.NewEncoder(w).Encode(&output)
		if err != nil {
			debug.ErrorMessag.Update(
				http.StatusInternalServerError, 
				"Notification.GET() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessag.Print(w)
			return
		}
	} else {
		if _, ok := notifications[id]; !ok {
			debug.ErrorMessag.Update(
				http.StatusNotFound, 
				"Notification.GET() -> Checking if ID exist",
				"ID validation: can't find ID",
				"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
			)
			debug.ErrorMessag.Print(w)
			return
		}
		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		output := notifications[id]
		err := json.NewEncoder(w).Encode(&output)
		if err != nil {
			debug.ErrorMessag.Update(
				http.StatusInternalServerError, 
				"Notification.GET() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessag.Print(w)
			return
		}
	}
}

func (notification *Notification) DELETE(w http.ResponseWriter, r *http.Request) {
	arrPath := strings.Split(r.URL.Path, "/")
	//branch if there aren't enough elements in URL and return error
	if len(arrPath) != 5 {
		debug.ErrorMessag.Update(
			http.StatusBadRequest, 
			"Notification.DELETE() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	id := arrPath[4]
	if _, ok := notifications[id]; !ok {
		debug.ErrorMessag.Update(
			http.StatusNotFound, 
			"Notification.DELETE() -> Checking if ID exist",
			"ID validation: can't find ID",
			"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessag.Print(w)
		return
	}
	delete(notifications, id)
	//create feedback message to send to client
	var feedback Feedback
	feedback.update(
		http.StatusOK, 
		"Webhook successfully deleted",
		id,
	)
	err := feedback.print(w)
	if err != nil {
		debug.ErrorMessag.Update(
			http.StatusInternalServerError, 
			"Notification.DELETE() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessag.Print(w)
		return
	}
}
