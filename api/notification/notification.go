package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Notification struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

func (notification *Notification) POST(w http.ResponseWriter, r *http.Request) {
	// Expects incoming body in terms of WebhookRegistration struct
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
	}
	notifications = append(notifications, *notification)
	// Note: Approach does not guarantee persistence or permanence of resource id (for CRUD)
	//fmt.Fprintln(w, len(webhooks)-1)
	fmt.Println("Webhook " + notification.URL + " has been registered.")
	http.Error(w, strconv.Itoa(len(notifications)-1), http.StatusCreated)
}

func (notification *Notification) GET(w http.ResponseWriter, r *http.Request) {
	// For now just return all webhooks, don't respond to specific resource requests
	err := json.NewEncoder(w).Encode(notifications)
	if err != nil {
		http.Error(w, "Something went wrong: " + err.Error(), http.StatusInternalServerError)
	}
}
