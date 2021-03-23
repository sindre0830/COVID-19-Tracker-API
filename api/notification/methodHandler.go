package notification

import (
	"main/debug"
	"net/http"
)

// Initialize signature (via init())
var SignatureKey = "X-SIGNATURE"
//var Mac hash.Hash
var Secret []byte

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			var notification Notification
			notification.POST(w, r)
		case http.MethodGet:
			var notification Notification
			notification.GET(w, r)
		case http.MethodDelete:
			var notification Notification
			notification.DELETE(w, r)
		default:
			debug.ErrorMessage.Update(
				http.StatusMethodNotAllowed, 
				"MethodHandler() -> Validating method",
				"method validation: wrong method",
				"Method not implemented.",
			)
			debug.ErrorMessage.Print(w)
	}
}
