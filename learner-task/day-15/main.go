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

const (
	// Placeholders keep values separate from SQL syntax.
	insertObservation = `INSERT INTO observations
		(id, patient_id, code, value, unit, observed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	selectObservation = `SELECT id, patient_id, code, value, unit, observed_at, created_at
		FROM observations WHERE id = $1`
	updateObservation = `UPDATE observations SET value = $1, unit = $2 WHERE id = $3`
)

var errNotFound = errors.New("observation not found")

type rowQueryer interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func getObservation(ctx context.Context, db rowQueryer, id string) (Observation, error) {
	var item Observation
	err := db.QueryRowContext(ctx, selectObservation, id).Scan(
		&item.ID, &item.PatientID, &item.Code, &item.Value, &item.Unit, &item.ObservedAt, &item.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		// Translate the driver-level result into a store-level not-found error.
		return Observation{}, errNotFound
	}
	if err != nil {
		return Observation{}, fmt.Errorf("scan observation: %w", err)
	}
	return item, nil
}

const migrationAddSource = `ALTER TABLE observations ADD COLUMN source text NULL;`
const migrationRequireSource = `ALTER TABLE observations ALTER COLUMN source SET DEFAULT 'manual';
UPDATE observations SET source = 'manual' WHERE source IS NULL;
ALTER TABLE observations ALTER COLUMN source SET NOT NULL;`

func main() {
	item := Observation{
		ID: "obs-1", PatientID: "patient-1", Code: "heart-rate", Value: 72,
		Unit: "bpm", ObservedAt: time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2026, 10, 31, 12, 0, 1, 0, time.UTC),
	}
	args := []any{item.ID, item.PatientID, item.Code, item.Value, item.Unit, item.ObservedAt, item.CreatedAt}
	fmt.Println("parameterized insert:", insertObservation)
	fmt.Printf("bound args: %#v\n", args)
	fmt.Println("select:", selectObservation)
	fmt.Println("update:", updateObservation)
	fmt.Printf("migration stages: %q then %q\n", migrationAddSource, migrationRequireSource)
	fmt.Println("No database connection is opened by this stdlib-only reference.")
}
