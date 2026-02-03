package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Key       string
	Owner     string
	CreatedAt time.Time
	ExpiredAt time.Time
}

func NewSubscription(owner string, days int) Subscription {
	expirationTime := time.Now().AddDate(0, 0, days).Round(time.Hour)
	return Subscription{
		Key:       uuid.NewString(),
		Owner:     owner,
		CreatedAt: time.Now(),
		ExpiredAt: expirationTime,
	}
}
