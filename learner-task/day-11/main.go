package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Observation struct {
	ID         string  `json:"id"`
	PatientID  string  `json:"patient_id"`
	Code       string  `json:"code"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"`
	ObservedAt string  `json:"observed_at"`
	Source     string  `json:"source,omitempty"`
}

type JSONStore struct{ Path string }

func (s JSONStore) Load() ([]Observation, error) {
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		// This example treats a missing file as a first run, not as corruption.
		return []Observation{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read observations: %w", err)
	}
	var observations []Observation
	if err := json.Unmarshal(data, &observations); err != nil {
		return nil, fmt.Errorf("decode observations: %w", err)
	}
	return observations, nil
}

func (s JSONStore) Save(observations []Observation) error {
	// A sibling temp file keeps the rename on the same filesystem.
	dir := filepath.Dir(s.Path)
	file, err := os.CreateTemp(dir, "."+filepath.Base(s.Path)+"-*")
	if err != nil {
		return fmt.Errorf("create temporary store: %w", err)
	}
	temp := file.Name()
	renamed := false
	defer func() {
		if !renamed {
			_ = os.Remove(temp)
		}
	}()
	if err := file.Chmod(0600); err != nil {
		_ = file.Close()
		return err
	}
	if err := json.NewEncoder(file).Encode(observations); err != nil {
		_ = file.Close()
		return fmt.Errorf("encode observations: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync observations: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close observations: %w", err)
	}
	if err := os.Rename(temp, s.Path); err != nil {
		return fmt.Errorf("replace observation store: %w", err)
	}
	// Keep the temp-file cleanup from removing the file after a successful rename.
	renamed = true
	return nil
}

func main() {
	dir, err := os.MkdirTemp(".", "day11-reference-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	store := JSONStore{Path: filepath.Join(dir, "observations.json")}
	if _, err := store.Load(); err != nil {
		panic(err)
	}
	legacy := `[{"id":"obs-1","patient_id":"patient-1","code":"heart-rate","value":72,"unit":"bpm","observed_at":"2026-10-31T12:00:00Z"}]`
	if err := os.WriteFile(store.Path, []byte(legacy), 0600); err != nil {
		panic(err)
	}
	items, err := store.Load()
	if err != nil || len(items) != 1 || items[0].Source != "" {
		panic("legacy optional field did not decode")
	}
	items[0].Source = "manual"
	if err := store.Save(items); err != nil {
		panic(err)
	}
	loaded, err := store.Load()
	if err != nil || loaded[0].Source != "manual" {
		panic("atomic round trip failed")
	}
	if err := os.WriteFile(store.Path, []byte("{"), 0600); err != nil {
		panic(err)
	}
	if _, err := store.Load(); err == nil {
		panic("malformed file was accepted")
	}
	fmt.Println("JSON store: missing file, legacy decode, atomic save, and malformed-file checks passed")
}
