package main

import (
	"errors"
	"testing"
)

func TestObservationValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   Observation
		wantErr bool
	}{
		{name: "valid", input: Observation{Label: "temperature", Value: 36.8, Unit: "°C"}},
		{name: "missing label", input: Observation{Value: 36.8, Unit: "°C"}, wantErr: true},
		{name: "missing unit", input: Observation{Label: "temperature", Value: 36.8}, wantErr: true},
		{name: "negative", input: Observation{Label: "temperature", Value: -1, Unit: "°C"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.input.Validate()
			if (err != nil) != test.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %t", err, test.wantErr)
			}
		})
	}
}

type fakeWriter struct {
	saved []Observation
	err   error
}

func (f *fakeWriter) Save(observation Observation) error {
	f.saved = append(f.saved, observation)
	return f.err
}

func TestRecordSavesValidObservation(t *testing.T) {
	writer := &fakeWriter{}
	observation := Observation{Label: "temperature", Value: 36.8, Unit: "°C"}

	if err := record(writer, observation); err != nil {
		t.Fatalf("record() error = %v", err)
	}
	if len(writer.saved) != 1 || writer.saved[0] != observation {
		t.Fatalf("saved = %#v, want one observation %#v", writer.saved, observation)
	}
}

func TestRecordPropagatesWriterError(t *testing.T) {
	wantErr := errors.New("storage unavailable")
	writer := &fakeWriter{err: wantErr}
	observation := Observation{Label: "temperature", Value: 36.8, Unit: "°C"}

	err := record(writer, observation)
	if !errors.Is(err, wantErr) {
		t.Fatalf("record() error = %v, want wrapped %v", err, wantErr)
	}
}
