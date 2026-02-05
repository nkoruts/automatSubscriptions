package subscription

import "errors"

var ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
var ErrSubscriptionNotFound = errors.New("subscription not found")
var ErrUnregisteredUserDevice = errors.New("deviceId is not registered")
