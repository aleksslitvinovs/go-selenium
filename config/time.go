package config

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
	return json.Marshal(t.String())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse duration")
	}

	t.Duration = d

	return nil
}
