package goof

import (
	"reflect"
	"testing"
)

type testMode struct {
	A string `json:"a_json" db:"a"`
	B int    `json:"b_json" db:"b"`
}

func TestGetDbColumns(t *testing.T) {
	columns := getDbColumns(testMode{})
	if !reflect.DeepEqual(columns, []string{"a", "b"}) {
		t.Errorf("expected %v, got %v", []string{"a", "b"}, columns)
	}
}

func TestGetJsonColumns(t *testing.T) {
	columns := getJsonColumns(testMode{})
	expected := []string{"a_json", "b_json"}
	if !reflect.DeepEqual(columns, expected) {
		t.Errorf("expected %v, got %v", expected, columns)
	}
}

func TestGetRouteName(t *testing.T) {
	name := getRouteName(testMode{})
	expected := "test-mode"
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

func TestGetTableName(t *testing.T) {
	name := getTableName(testMode{})
	expected := "test_mode"
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}
