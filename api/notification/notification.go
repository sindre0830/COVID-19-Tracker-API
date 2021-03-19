package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var webhooks []Notification

type Notification struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			// Expects incoming body in terms of WebhookRegistration struct
			webhook := Notification{}
			err := json.NewDecoder(r.Body).Decode(&webhook)
			if err != nil {
				http.Error(w, "Something went wrong: " + err.Error(), http.StatusBadRequest)
			}
			webhooks = append(webhooks, webhook)
			// Note: Approach does not guarantee persistence or permanence of resource id (for CRUD)
			//fmt.Fprintln(w, len(webhooks)-1)
			fmt.Println("Webhook " + webhook.URL + " has been registered.")
			http.Error(w, strconv.Itoa(len(webhooks)-1), http.StatusCreated)
		case http.MethodGet:
			// For now just return all webhooks, don't respond to specific resource requests
			err := json.NewEncoder(w).Encode(webhooks)
			if err != nil {
				http.Error(w, "Something went wrong: " + err.Error(), http.StatusInternalServerError)
			}
		case http.MethodDelete:
			//
		default: http.Error(w, "Invalid method " + r.Method, http.StatusBadRequest)
	}
}
