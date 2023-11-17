package cmp

import (
	"errors"
	"fmt"
	"reflect"
)

// DO NOT USE. Ensure that all of the fields in expected match the corresponding fields in actual. Ignores fields in actual that
// aren't in expected.
func SubsetOf(superset, subset any) (err error) {

	// If the expected value is nil, then we don't care what the actual value is.
	if superset == nil {
		return errors.New("expected is nil")
	}

	// If the actual value is nil, then we can't compare it to the expected value.
	if subset == nil {
		return errors.New("actual is nil")
	}

	// If the expected value is a pointer, then we need to dereference it.
	if reflect.TypeOf(superset).Kind() == reflect.Ptr {
		superset = reflect.ValueOf(superset).Elem().Interface()
	}

	// If the actual value is a pointer, then we need to dereference it.
	if reflect.TypeOf(subset).Kind() == reflect.Ptr {
		subset = reflect.ValueOf(subset).Elem().Interface()
	}

	expectedKind := reflect.TypeOf(superset).Kind()
	// If the expected value is a struct, then we need to compare each of its fields to the corresponding field in the
	// actual value.
	if expectedKind == reflect.Struct {
		for i := 0; i < reflect.TypeOf(superset).NumField(); i++ {
			field := reflect.TypeOf(superset).Field(i)
			if err = SubsetOf(reflect.ValueOf(superset).FieldByName(field.Name).Interface(), reflect.ValueOf(subset).FieldByName(field.Name).Interface()); err != nil {
				return
			}
		}
		return
	}

	// If the expected value is a map, then we need to compare each of its keys and values to the corresponding key and
	// value in the actual value.
	if expectedKind == reflect.Map {
		for _, key := range reflect.ValueOf(superset).MapKeys() {
			if err = SubsetOf(reflect.ValueOf(superset).MapIndex(key).Interface(), reflect.ValueOf(subset).MapIndex(key).Interface()); err != nil {
				return
			}
		}
	}

	// If the expected value is a slice, then we need to compare each of its elements to the corresponding element in
	// the actual value.
	if expectedKind == reflect.Slice {
		for i := 0; i < reflect.ValueOf(superset).Len(); i++ {
			if err = SubsetOf(reflect.ValueOf(superset).Index(i).Interface(), reflect.ValueOf(subset).Index(i).Interface()); err != nil {
				return
			}
		}
	}

	// If the expected value is a string, then we need to compare it to the actual value.
	if expectedKind == reflect.String {
		val := reflect.ValueOf(superset).String()
		if val == "" {
			return
		}
		if val != reflect.ValueOf(subset).String() {
			return fmt.Errorf("expected string '%s'; got '%s'", reflect.ValueOf(superset).String(), reflect.ValueOf(subset).String())
		}
	}

	// If the expected value is a bool, then we need to compare it to the actual value.
	if expectedKind == reflect.Bool {
		if reflect.ValueOf(superset).Bool() != reflect.ValueOf(subset).Bool() {
			return fmt.Errorf("expected bool '%t'; got '%t'", reflect.ValueOf(superset).Bool(), reflect.ValueOf(subset).Bool())
		}
	}

	// If the expected value is an int, then we need to compare it to the actual value.
	if expectedKind == reflect.Int || expectedKind == reflect.Int8 || expectedKind == reflect.Int16 || expectedKind == reflect.Int32 || expectedKind == reflect.Int64 {
		if reflect.ValueOf(superset).Int() != reflect.ValueOf(subset).Int() {
			return fmt.Errorf("expected int '%d'; got '%d'", reflect.ValueOf(superset).Int(), reflect.ValueOf(subset).Int())
		}
	}

	// If the expected value is a uint, then we need to compare it to the actual value.
	if expectedKind == reflect.Uint || expectedKind == reflect.Uint8 || expectedKind == reflect.Uint16 || expectedKind == reflect.Uint32 || expectedKind == reflect.Uint64 {
		if reflect.ValueOf(superset).Uint() != reflect.ValueOf(subset).Uint() {
			return fmt.Errorf("expected uint '%d'; got '%d'", reflect.ValueOf(superset).Uint(), reflect.ValueOf(subset).Uint())
		}
	}

	// If the expected value is a float, then we need to compare it to the actual value.
	if expectedKind == reflect.Float64 || expectedKind == reflect.Float32 {
		if reflect.ValueOf(superset).Float() != reflect.ValueOf(subset).Float() {
			return fmt.Errorf("expected float '%f'; got '%f'", reflect.ValueOf(superset).Float(), reflect.ValueOf(subset).Float())
		}
	}

	// If the expected value is a function, then we need to compare it to the actual value.
	if expectedKind == reflect.Func {
		if reflect.ValueOf(superset).Pointer() != reflect.ValueOf(subset).Pointer() {
			return errors.New("expected func" + reflect.ValueOf(superset).String() + "; got " + reflect.ValueOf(subset).String())
		}
	}

	// If the expected value is a channel, then we need to compare it to the actual value.
	if expectedKind == reflect.Chan {
		if reflect.ValueOf(superset).Pointer() != reflect.ValueOf(subset).Pointer() {
			return errors.New("expected chan" + reflect.ValueOf(superset).String() + "; got " + reflect.ValueOf(subset).String())
		}
	}

	// If the expected value is an interface, then we need to compare it to the actual value.
	if expectedKind == reflect.Interface {
		if reflect.ValueOf(superset).Pointer() != reflect.ValueOf(subset).Pointer() {
			return errors.New("expected interface" + reflect.ValueOf(superset).String() + "; got " + reflect.ValueOf(subset).String())
		}
	}

	return
}
