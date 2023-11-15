// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package driver

import (
	"errors"
	"fmt"
)

const (
	// EnvironmentLocal is a Environment of type local.
	EnvironmentLocal Environment = "local"
	// EnvironmentCloudRun is a Environment of type cloud_run.
	EnvironmentCloudRun Environment = "cloud_run"
)

var ErrInvalidEnvironment = errors.New("not a valid Environment")

// String implements the Stringer interface.
func (x Environment) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Environment) IsValid() bool {
	_, err := ParseEnvironment(string(x))
	return err == nil
}

var _EnvironmentValue = map[string]Environment{
	"local":     EnvironmentLocal,
	"cloud_run": EnvironmentCloudRun,
}

// ParseEnvironment attempts to convert a string to a Environment.
func ParseEnvironment(name string) (Environment, error) {
	if x, ok := _EnvironmentValue[name]; ok {
		return x, nil
	}
	return Environment(""), fmt.Errorf("%s is %w", name, ErrInvalidEnvironment)
}

// MarshalText implements the text marshaller method.
func (x Environment) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Environment) UnmarshalText(text []byte) error {
	tmp, err := ParseEnvironment(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *Environment) Set(val string) error {
	v, err := ParseEnvironment(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *Environment) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *Environment) Type() string {
	return "Environment"
}

const (
	// TypePostgres is a Type of type postgres.
	TypePostgres Type = "postgres"
	// TypeMysql is a Type of type mysql.
	TypeMysql Type = "mysql"
	// TypeSqlite3 is a Type of type sqlite3.
	TypeSqlite3 Type = "sqlite3"
)

var ErrInvalidType = errors.New("not a valid Type")

// String implements the Stringer interface.
func (x Type) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Type) IsValid() bool {
	_, err := ParseType(string(x))
	return err == nil
}

var _TypeValue = map[string]Type{
	"postgres": TypePostgres,
	"mysql":    TypeMysql,
	"sqlite3":  TypeSqlite3,
}

// ParseType attempts to convert a string to a Type.
func ParseType(name string) (Type, error) {
	if x, ok := _TypeValue[name]; ok {
		return x, nil
	}
	return Type(""), fmt.Errorf("%s is %w", name, ErrInvalidType)
}

// MarshalText implements the text marshaller method.
func (x Type) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Type) UnmarshalText(text []byte) error {
	tmp, err := ParseType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *Type) Set(val string) error {
	v, err := ParseType(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *Type) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *Type) Type() string {
	return "Type"
}
