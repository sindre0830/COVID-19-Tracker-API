package notification

// NotificationInput structure stores client input for a webhook.
type NotificationInput struct {
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}
