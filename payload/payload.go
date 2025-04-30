package payload

import (
	"time"
)

type Payload struct {
	SentTime time.Time `json:"sent_time"`
	Contents string    `json:"contents"`
}
