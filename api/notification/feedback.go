package notification

import (
	"encoding/json"
	"net/http"
)

// Feedback structure stores information about successful method request.
//
// Functionality: update, print
type Feedback struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ID		   string `json:"id"`
}

// update sets new data in structure.
func (feedback *Feedback) update(status int, message string, id string) {
	feedback.StatusCode = status
	feedback.Message = message
	feedback.ID = id
}

// print sends structure to client.
func (feedback Feedback) print(w http.ResponseWriter) error {
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(feedback.StatusCode)
	//send output to user and branch if an error occured
	err := json.NewEncoder(w).Encode(feedback)
	if err != nil {
		return err
	}
	return nil
}
