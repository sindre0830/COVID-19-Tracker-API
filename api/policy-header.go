package api

// Policy stores data about COVID policies based on user input.
//
// Functionality: Handler, get, getTotal, getHistory, update
type Policy struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float64 `json:"stringency"`
	Trend	   float64 `json:"trend"`
	Update 	   string  `json:"update"`
}
