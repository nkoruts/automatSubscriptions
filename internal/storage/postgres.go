package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nkoruts/automatSubscriptions/internal/subscription"
)

type PostgresStorage struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresStorage{ctx: ctx, db: db}, nil
}

func (p *PostgresStorage) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS subscriptions (
   		key UUID PRIMARY KEY,
    	owner TEXT NOT NULL,
    	device_id TEXT,
    	created_at TIMESTAMP NOT NULL,
    	expired_at TIMESTAMP NOT NULL
	);
	`
	_, err := p.db.Exec(p.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) AddSubscription(owner string, days int) error {
	sub := subscription.NewSubscription(owner, days)

	query := `
		INSERT INTO subscriptions (key, owner, device_id, created_at, expired_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := p.db.Exec(p.ctx, query,
		sub.Key,
		sub.Owner,
		nil,
		sub.CreatedAt,
		sub.ExpiredAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) DeleteSubscription(key string) error {
	query := `DELETE FROM subscriptions WHERE key = $1`
	tag, err := p.db.Exec(p.ctx, query, key)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return subscription.ErrSubscriptionNotFound
	}
	return nil
}

func (p *PostgresStorage) UpdateSubscription(key, deviceId string) error {
	query := `UPDATE subscriptions SET device_id = $1 WHERE key = $2`
	tag, err := p.db.Exec(p.ctx, query, deviceId, key)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return subscription.ErrSubscriptionNotFound
	}
	return nil
}

func (p *PostgresStorage) CheckSubscription(key, deviceId string) (bool, error) {
	query := `SELECT device_id, expired_at FROM subscriptions WHERE key = $1`

	var dbDeviceID *string
	var expiredAt time.Time

	err := p.db.QueryRow(p.ctx, query, key).Scan(&dbDeviceID, &expiredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, subscription.ErrSubscriptionNotFound
		}
		return false, err
	}

	if time.Now().After(expiredAt) {
		return false, subscription.ErrSubscriptionNotFound
	}

	if dbDeviceID == nil {
		return false, nil
	}

	if *dbDeviceID != deviceId {
		return false, subscription.ErrUnregisteredUserDevice
	}

	return true, nil
}

func (p *PostgresStorage) GetList() map[string]subscription.Subscription {
	query := `SELECT key, owner, device_id, created_at, expired_at FROM subscriptions`

	rows, err := p.db.Query(p.ctx, query)
	if err != nil {
		return map[string]subscription.Subscription{}
	}
	defer rows.Close()

	result := make(map[string]subscription.Subscription)

	for rows.Next() {
		var sub subscription.Subscription
		err := rows.Scan(&sub.Key, &sub.Owner, &sub.DeviceId, &sub.CreatedAt, &sub.ExpiredAt)
		if err != nil {
			continue
		}
		result[sub.Key] = sub
	}
	return result
}
