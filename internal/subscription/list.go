package subscription

import (
	"maps"
	"sync"
)

type List struct {
	subscriptions map[string]Subscription
	mtx           sync.RWMutex
}

func NewList() *List {
	return &List{
		subscriptions: make(map[string]Subscription),
	}
}

func (l *List) GetList() map[string]Subscription {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tmp := maps.Clone(l.subscriptions)
	return tmp
}

func (l *List) AddSubscription(owner string, days int) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	subscription := NewSubscription(owner, days)

	if _, ok := l.subscriptions[subscription.Key]; ok {
		return ErrSubscriptionAlreadyExists
	}
	l.subscriptions[subscription.Key] = subscription

	return nil
}

func (l *List) DeleteSubscription(key string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.subscriptions[key]; !ok {
		return ErrSubscriptionNotFound
	}

	delete(l.subscriptions, key)
	return nil
}

func (l *List) UpdateSubscription(key, deviceId string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	subscription, ok := l.subscriptions[key]
	if !ok {
		return ErrSubscriptionNotFound
	}

	idCopy := deviceId
	subscription.DeviceId = &idCopy

	l.subscriptions[key] = subscription
	return nil
}
