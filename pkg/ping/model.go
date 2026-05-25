package ping

// resp represents the ping response.
//
// It has a single field, Timestamp, which is a string representing the timestamp.
type resp struct {
	Timestamp string `json:"timestamp"`
}
