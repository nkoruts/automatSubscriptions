package server

import (
	"encoding/json"
	"errors"
	"time"
)

// SubscriptionDTO
type SubscriptionDTO struct {
	Owner string
	Days  int
}

func (s *SubscriptionDTO) ValidateRequest() error {
	if s.Owner == "" {
		return errors.New("owner is empty")
	}
	if s.Days <= 0 {
		return errors.New("days is incorrent")
	}

	return nil
}

// SuccessDTO
type SuccessDTO struct {
	Success bool
}

func (s *SuccessDTO) ToString() string {
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ErrorDTO
type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// UpdateDTO
type UpdateDTO struct {
	DeviceId string
}

// CheckDTO
type CheckDTO struct {
	Key      string `json:"key"`
	DeviceID string `json:"deviceId"`
}

func (c *CheckDTO) ValidateRequest() error {
	if c.Key == "" {
		return errors.New("incorrect key")
	}
	return nil
}

type CheckResponse struct {
	Active bool `json:"active"`
}
