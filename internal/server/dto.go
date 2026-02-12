package server

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// SubscriptionDTO
type SubscriptionDTO struct {
	Owner string `json:"owner"`
	Days  int    `json:"days"`
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
	Success bool `json:"success"`
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
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
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
	DeviceId string `json:"deviceId"`
}

func (u *UpdateDTO) ValidateRequest() error {
	if u.DeviceId == "" {
		return errors.New("empty deviceId")
	}
	return nil
}

// CheckDTO
type CheckDTO struct {
	Key      string `json:"key"`
	DeviceID string `json:"deviceId"`
}

func (c *CheckDTO) ValidateRequest() error {
	_, err := uuid.Parse(c.Key)
	if err != nil {
		return err
	}
	if c.DeviceID == "" {
		return errors.New("empty deviceId")
	}
	return nil
}

type CheckResponse struct {
	Active bool `json:"active"`
}
