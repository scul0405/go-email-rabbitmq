package email

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	EmailID     uuid.UUID `json:"emailId"`
	To          []string  `json:"to"`
	From        string    `json:"from,omitempty"`
	Body        string    `json:"body"`
	Subject     string    `json:"subject"`
	ContentType string    `json:"contentType,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
