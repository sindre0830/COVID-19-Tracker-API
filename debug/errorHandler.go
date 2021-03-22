package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorMessag is a global variable.
var ErrorMessage Debug
// Debug structure keeps all error data.
//
// Functionality: Update, Print
type Debug struct {
	StatusCode 		 int    `json:"status_code"`
	Location   		 string `json:"location"`
	RawError   		 string `json:"raw_error"`
	PossibleReason   string `json:"possible_reason"`
}
// Update adds new information to error msg.
func (debug *Debug) Update(status int, loc string, err string, reason string) {
	debug.StatusCode = status
	debug.Location = loc
	debug.RawError = err
	//update reason if status code shows client error
	if status == http.StatusBadRequest || status == http.StatusNotFound || status == http.StatusMethodNotAllowed {
		debug.PossibleReason = reason
	} else {
		debug.PossibleReason = "Unknown"
	}
}
// Print prints error msg to user and terminal.
func (debug *Debug) Print(w http.ResponseWriter) {
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(debug.StatusCode)
	//send error output to user
	err := json.NewEncoder(w).Encode(debug)
	//branch if something went wrong with output
	if err != nil {
		fmt.Println("ERROR encoding JSON in Debug.Print()", err)
		return
	}
	//send error output to console
	fmt.Printf("\nError {\n\tstatus_code:\t\t%v,\n\tlocation:\t\t%s,\n\traw_error:\t\t%s,\n\tpossible_reason:\t%s\n}\n", debug.StatusCode, debug.Location, debug.RawError, debug.PossibleReason)
}
