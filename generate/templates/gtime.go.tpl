{{define "gtime.go"}}
// Package gtime provides database wrappers for all of the formats provided in the time package that will format and 
// parse the time.Time type to and from the database using that format.

package gtime

import (
  "database/sql/driver"
  "fmt"
  "time"
)

// parseLayout parses a time.Time from a string or []byte using the given layout
func parseLayout(src any, layout string) (t time.Time, err error) {
	switch src := src.(type) {
	case time.Time:
		t = src
	case string:
		t, err = time.Parse(layout, src)
		if err != nil {
			return
		}
	case []byte:
		t, err = time.Parse(layout, string(src))
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unsupported type: %T", src)
	}
	return
}

type formattable interface {
	Format(layout string) string
}

{{ range .Formats}}
const Time{{.Name}}Format = {{.Format}} // The format used for {{.Name}}

// Time{{.Name}} is a wrapper around time.Time that implements the sql.Scanner and driver.Valuer with the 
// format `Time{{.Name}}Format`. If the value is zero when `Value()` is called, it will use the current time
type Time{{.Name}} struct {
  time.Time
}

// Implements the `sql.Scanner` interface
func (t *Time{{.Name}}) Scan(src interface{}) (err error) {
	t.Time, err = parseLayout(src, Time{{.Name}}Format)
	return
}

// Implements the `driver.Valuer` interface
func (t Time{{.Name}}) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t Time{{.Name}}) String() string {
	return t.Time.Format(Time{{.Name}}Format)
}

// Uses the `Format` method to compare the two values instead of the default
func (t Time{{.Name}}) Equal(other formattable) bool {
	return t.String() == other.Format(Time{{.Name}}Format)
}

{{end}}

{{end}}
