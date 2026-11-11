package main

import (
	"errors"
	"testing"
)

func TestObservationValidate(t *testing.T) {
	// Table tests keep related inputs and expected outcomes together.
	tests := []struct {
		name    string
		input   Observation
		wantErr bool
	}{
		{name: "valid", input: Observation{Label: "temperature", Value: 36.8, Unit: "°C"}},
		{name: "missing label", input: Observation{Value: 36.8, Unit: "°C"}, wantErr: true},
		{name: "blank label", input: Observation{Label: " ", Value: 1, Unit: "x"}, wantErr: true},
		{name: "missing unit", input: Observation{Label: "temperature", Value: 36.8}, wantErr: true},
		{name: "negative value", input: Observation{Label: "x", Value: -1, Unit: "x"}, wantErr: true},
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

func TestRecordSavesExactlyOnce(t *testing.T) {
	writer := &fakeWriter{}
	observation := Observation{Label: "temperature", Value: 36.8, Unit: "°C"}
	if err := record(writer, observation); err != nil {
		t.Fatalf("record() error = %v", err)
	}
	if len(writer.saved) != 1 || writer.saved[0] != observation {
		t.Fatalf("saved = %#v, want exactly one %#v", writer.saved, observation)
	}
}

func TestRecordWrapsWriterError(t *testing.T) {
	wantErr := errors.New("storage unavailable")
	writer := &fakeWriter{err: wantErr}
	observation := Observation{Label: "temperature", Value: 36.8, Unit: "°C"}
	err := record(writer, observation)
	if !errors.Is(err, wantErr) {
		t.Fatalf("record() error = %v, want wrapped %v", err, wantErr)
	}
}
