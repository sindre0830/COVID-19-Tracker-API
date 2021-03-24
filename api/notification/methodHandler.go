package notification

import (
	"main/debug"
	"net/http"
)

// MethodHandler handles the method of a http request.
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
