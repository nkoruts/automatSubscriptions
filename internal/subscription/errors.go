package subscription

import "errors"

var ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
var ErrSubscriptionNotFound = errors.New("subscription not found")
