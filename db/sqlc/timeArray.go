package db

import (
	"database/sql/driver"
	"errors"
	"time"
)

import (
	"github.com/lib/pq"
)

// ErrParseData signifies that an error occured while parsing SQL data
var ErrParseData = errors.New("Unable to parse SQL data")

// PgTime wraps a time.Time
type PgTime struct{ time.Time }

// Scan implements the sql.Scanner interface
func (t *PgTime) Scan(val interface{}) error {
	switch v := val.(type) {
	case time.Time:
		t.Time = v
		return nil
	case []uint8: // byte is the same as uint8: https://golang.org/pkg/builtin/#byte
		_t, err := pq.ParseTimestamp(nil, string(v))
		if err != nil {
			return ErrParseData
		}
		t.Time = _t
		return nil
	case string:
		_t, err := pq.ParseTimestamp(nil, v)
		if err != nil {
			return ErrParseData
		}
		t.Time = _t
		return nil
	}
	return ErrParseData
}

// Value implements the driver.Valuer interface
func (t *PgTime) Value() (driver.Value, error) { return pq.FormatTimestamp(t.Time), nil }

// PgTimeArray wraps a time.Time slice to be used as a Postgres array
// type PgTimeArray []time.Time
type PgTimeArray []PgTime

// type PgTimeArray []pq.NullTime

// Scan implements the sql.Scanner interface
func (a *PgTimeArray) Scan(src interface{}) error { return pq.GenericArray{a}.Scan(src) }

// Value implements the driver.Valuer interface
func (a *PgTimeArray) Value() (driver.Value, error) { return pq.GenericArray{a}.Value() }
