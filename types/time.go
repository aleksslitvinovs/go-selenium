package types

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Time is a wrapper around time.Duration with improved JSON (un)marshalling.
type Time struct {
	time.Duration
}

func (t *Time) String() string {
	return t.Duration.String()
}

// MarshalJSON marshals Time.
func (t *Time) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(t.String())
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to marshal duration")
	}

	return data, nil
}

// UnmarshalJSON unmarshals Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errors.Wrap(err, "failed to unmarshal duration")
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse duration")
	}

	t.Duration = d

	return nil
}
