package henrymedscodechallenge

import "time"

type Availability struct {
	ProviderID string    `json:"provider_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
}

type Status string

const (
	Pending   Status = "pending"
	Confirmed Status = "confirmed"
)

type Appointment struct {
	ID         string    `json:"id"`
	ProviderID string    `json:"provider_id"`
	ClientID   string    `json:"client_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Status     Status    `json:"status"`
	Expires    time.Time `json:"expires"`
}
