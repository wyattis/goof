package conf

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"github.com/wyattis/z/ztime"
)

type flagTestStruct struct {
	String    string
	Int       int
	Int8      int8
	Int16     int16
	Int32     int32
	Int64     int64
	UseBool   bool
	CreatedAt time.Time
	Nested    struct {
		String string
	}
}

func TestZflag(t *testing.T) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	s := flagTestStruct{}
	c := FlagConfigurer{flagSet: set}
	if err := c.Init(&s); err != nil {
		t.Error(err)
	}
	c.flagSet.Usage()
	args := []string{"-string", "hello", "-int", "1", "-int8", "2", "-int16", "3", "-int32", "4", "-int64", "5", "-use-bool", "-nested-string", "world", "-created-at", "2017-01-01T00:00:00Z"}
	if err := c.Apply(&s, args...); err != nil {
		t.Error(err)
	}
	expected := flagTestStruct{
		String:    "hello",
		Int:       1,
		Int8:      2,
		Int16:     3,
		Int32:     4,
		Int64:     5,
		UseBool:   true,
		CreatedAt: ztime.MustParse("2017-01-01T00:00:00Z", "2006-01-02T15:04:05Z"),
		Nested: struct {
			String string
		}{
			String: "world",
		},
	}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v, but got %v", expected, s)
	}
}
