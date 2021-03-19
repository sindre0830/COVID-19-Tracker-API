package notification

type notificationGetOne struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Timeout     int    `json:"timeout"`
	Information string `json:"information"`
	Country     string `json:"country"`
	Trigger     string `json:"trigger"`
}
