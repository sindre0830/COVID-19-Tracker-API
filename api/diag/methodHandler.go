package diag

import (
	"main/debug"
	"net/http"
)

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var diagnosis Diagnosis
			diagnosis.Handler(w, r)
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
