package ping

import "time"

// Resp represents the ping response.
//
// It contains the server's current timestamp, indicating when the response
// was generated.
type Resp struct {
	// Timestamp is the exact time when the ping response was created.
	Timestamp time.Time `json:"timestamp"`
}
