package notification

import (
	"encoding/json"
	"net/http"
)


type Feedback struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ID		   string `json:"id"`
}

func (feedback *Feedback) update(status int, message string, id string) {
	feedback.StatusCode = status
	feedback.Message = message
	feedback.ID = id
}

func (feedback *Feedback) print(w http.ResponseWriter) error {
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(feedback.StatusCode)
	//send error output to user and branch if an error occurred
	err := json.NewEncoder(w).Encode(feedback)
	if err != nil {
		return err
	}
	return nil
}
