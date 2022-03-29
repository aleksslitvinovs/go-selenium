package types

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Time struct {
	time.Duration
}

func (t *Time) String() string {
	return t.Duration.String()
}

func (t *Time) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(t.String())
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to marshal duration")
	}

	return data, nil
}

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
