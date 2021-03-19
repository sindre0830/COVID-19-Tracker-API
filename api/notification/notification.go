package notification

import (
	"encoding/json"
	"main/debug"
	"net/http"
)

type Notification struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

func (notification *Notification) POST(w http.ResponseWriter, r *http.Request) {
	//read input from client and branch if an error occurred
	err := json.NewDecoder(r.Body).Decode(&notification)
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
	notifications = append(notifications, *notification)
	//create feedback message to send to client
	var feedback Feedback
	feedback.update(
		http.StatusCreated, 
		"Webhook successfully created for '" + notification.URL + "'.",
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
	// For now just return all webhooks, don't respond to specific resource requests
	err := json.NewEncoder(w).Encode(notifications)
	if err != nil {
		http.Error(w, "Something went wrong: " + err.Error(), http.StatusInternalServerError)
	}
}
