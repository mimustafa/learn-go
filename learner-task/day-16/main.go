package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Observation struct {
	ID, PatientID, Code, Unit string
	Value                     float64
	ObservedAt, CreatedAt     time.Time
}

type Tx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}

type Beginner interface {
	BeginTx(context.Context, *sql.TxOptions) (Tx, error)
}

type sqlBeginner struct{ db *sql.DB }

func configurePool(db *sql.DB, open, idle int, maxLifetime, maxIdle time.Duration) error {
	if open < 1 || idle < 0 || idle > open || maxLifetime <= 0 || maxIdle <= 0 {
		return errors.New("invalid database pool configuration")
	}
	db.SetMaxOpenConns(open)
	db.SetMaxIdleConns(idle)
	db.SetConnMaxLifetime(maxLifetime)
	db.SetConnMaxIdleTime(maxIdle)
	return nil
}

func printPoolStats(db *sql.DB) {
	stats := db.Stats()
	fmt.Printf("pool open=%d in_use=%d idle=%d wait_count=%d\n",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
}

func (db sqlBeginner) BeginTx(ctx context.Context, options *sql.TxOptions) (Tx, error) {
	tx, err := db.db.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

const insert = `INSERT INTO observations
	(id, patient_id, code, value, unit, observed_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

func importBatch(ctx context.Context, db Beginner, observations []Observation) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin import: %w", err)
	}
	// Roll back every early exit; after Commit, this deferred call is harmless.
	defer tx.Rollback()
	for _, item := range observations {
		if _, err := tx.ExecContext(ctx, insert,
			item.ID, item.PatientID, item.Code, item.Value, item.Unit,
			item.ObservedAt, item.CreatedAt); err != nil {
			return fmt.Errorf("insert observation %q: %w", item.ID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit import: %w", err)
	}
	return nil
}

type fakeTx struct {
	failAt, writes int
	committed      bool
	rolledBack     bool
}

func (tx *fakeTx) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	if tx.writes == tx.failAt {
		return nil, errors.New("injected insert failure")
	}
	tx.writes++
	return nil, nil
}
func (tx *fakeTx) Commit() error   { tx.committed = true; return nil }
func (tx *fakeTx) Rollback() error { tx.rolledBack = true; return nil }

type fakeDB struct{ tx *fakeTx }

func (db fakeDB) BeginTx(context.Context, *sql.TxOptions) (Tx, error) { return db.tx, nil }

func main() {
	now := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	batch := []Observation{
		{ID: "o1", PatientID: "p1", Code: "heart-rate", Value: 72, Unit: "bpm", ObservedAt: now, CreatedAt: now},
		{ID: "o2", PatientID: "p1", Code: "steps", Value: 100, Unit: "count", ObservedAt: now, CreatedAt: now},
	}
	success := &fakeTx{failAt: -1}
	if err := importBatch(context.Background(), fakeDB{success}, batch); err != nil || !success.committed {
		panic("successful transaction did not commit")
	}
	failure := &fakeTx{failAt: 1}
	if err := importBatch(context.Background(), fakeDB{failure}, batch); err == nil || !failure.rolledBack || failure.committed {
		panic("failed transaction did not roll back")
	}
	fmt.Printf("success committed=%v; injected failure rolled back=%v\n", success.committed, failure.rolledBack)
	fmt.Println("pool policy example: max-open=5 max-idle=2 lifetime=30m idle-time=5m")
}
