package notification

type NotificationInput struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}
