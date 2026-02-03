package server

import (
	"encoding/json"
	"errors"
	"time"
)

type SubscriptionDTO struct {
	Owner string
	Days  int
}

func (s *SubscriptionDTO) ValidateForCreate() error {
	if s.Owner == "" {
		return errors.New("owner is empty")
	}
	if s.Days <= 0 {
		return errors.New("days is incorrent")
	}

	return nil
}

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
