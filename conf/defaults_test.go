package conf

import (
	"reflect"
	"testing"
	"time"

	"github.com/wyattis/z/ztime"
)

type stringDefaults struct {
	One    string  `default:"one"`
	Two    string  `default:"two"`
	Ptr    *string `default:"ptr"`
	Nested struct {
		Three string `default:"three"`
	}
}

func ptrTo[T any](t T) *T {
	return &t
}

func TestStringDefaults(t *testing.T) {
	val := stringDefaults{}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := stringDefaults{
		One: "one",
		Two: "two",
		Ptr: ptrTo("ptr"),
		Nested: struct {
			Three string `default:"three"`
		}{Three: "three"},
	}

	if !reflect.DeepEqual(expected, val) {
		t.Errorf("Expected %+v, but got %+v", expected, val)
	}

}

type intDefaults struct {
	Int    int   `default:"1"`
	Int8   int8  `default:"8"`
	Int16  int16 `default:"16"`
	Int32  int32 `default:"32"`
	Int64  int64 `default:"64"`
	IntPtr *int  `default:"95"`
	Max    int   `default:"2147483647"`
	Nested struct {
		Three int `default:"3"`
	}
}

func TestIntDefaults(t *testing.T) {
	val := intDefaults{}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := intDefaults{
		Int:    1,
		Int8:   8,
		Int16:  16,
		Int32:  32,
		Int64:  64,
		IntPtr: ptrTo(95),
		Max:    2147483647,
		Nested: struct {
			Three int `default:"3"`
		}{Three: 3},
	}
	if val != expected {
		t.Errorf("Expected %+v, but got %+v", expected, val)
	}
}

type uintDefaults struct {
	Uint    uint   `default:"1"`
	Uint8   uint8  `default:"8"`
	Uint16  uint16 `default:"16"`
	Uint32  uint32 `default:"32"`
	Uint64  uint64 `default:"64"`
	UintPtr *uint  `default:"95"`
	Max     uint   `default:"4294967295"`
	Nested  struct {
		Three uint `default:"3"`
	}
}

func TestUintDefaults(t *testing.T) {
	val := uintDefaults{}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := uintDefaults{
		Uint:    1,
		Uint8:   8,
		Uint16:  16,
		Uint32:  32,
		Uint64:  64,
		UintPtr: ptrTo(uint(95)),
		Max:     4294967295,
		Nested: struct {
			Three uint `default:"3"`
		}{Three: 3},
	}
	if val != expected {
		t.Errorf("Expected %+v, but got %+v", expected, val)
	}
}

type floatDefaults struct {
	Float32  float32  `default:"32.32"`
	Float64  float64  `default:"64.64"`
	FloatPtr *float64 `default:"95.95"`
	Nested   struct {
		Three float32 `default:"3.3"`
	}
}

func TestFloatDefaults(t *testing.T) {
	val := floatDefaults{}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := floatDefaults{
		Float32:  32.32,
		Float64:  64.64,
		FloatPtr: ptrTo(95.95),
		Nested: struct {
			Three float32 `default:"3.3"`
		}{Three: 3.3},
	}
	if val != expected {
		t.Errorf("Expected %+v, but got %+v", expected, val)
	}
}

type boolDefaults struct {
	True    bool  `default:"true"`
	False   bool  `default:"false"`
	BoolPtr *bool `default:"true"`
}

func TestBoolDefaults(t *testing.T) {
	val := boolDefaults{}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := boolDefaults{
		True:    true,
		False:   false,
		BoolPtr: ptrTo(true),
	}
	if !reflect.DeepEqual(val, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, val)
	}
}

type timeDefaults struct {
	Now     time.Time  `default:"now"`
	Epoch   time.Time  `default:"1970-01-01T00:00:00Z"`
	TimePtr *time.Time `default:"2006-01-02T15:04:05Z"`
	Nested  struct {
		Three time.Time `default:"2006-01-02T15:04:05Z"`
	}
}

func TestTimeDefaults(t *testing.T) {
	val := timeDefaults{}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := timeDefaults{
		Now:     time.Now(),
		Epoch:   time.Unix(0, 0),
		TimePtr: ptrTo(ztime.MustParse("2006-01-02T15:04:05Z", "2006-01-02T15:04:05Z")),
		Nested: struct {
			Three time.Time `default:"2006-01-02T15:04:05Z"`
		}{Three: ztime.MustParse("2006-01-02T15:04:05Z", "2006-01-02T15:04:05Z")},
	}

	if !ztime.EqualWithin(val.Now, expected.Now, time.Second) {
		t.Errorf("Expected %+v, but got %+v", expected.Now, val.Now)
	}
	if !expected.Epoch.Equal(val.Epoch) {
		t.Errorf("Expected %+v, but got %+v", expected.Epoch, val.Epoch)
	}
	if !expected.Nested.Three.Equal(val.Nested.Three) {
		t.Errorf("Expected %+v, but got %+v", expected.Nested.Three, val.Nested.Three)
	}
}

type defaultsMap map[string]any
type defaultsMapSub struct {
	One string `default:"one"`
}

func TestStrMapDefaults(t *testing.T) {
	val := defaultsMap{
		"one": "one",
		"sub": defaultsMapSub{},
	}
	c := defaultConfigurer{}
	if err := c.Apply(&val); err != nil {
		t.Error(err)
	}
	expected := defaultsMap{
		"one": "one",
		"sub": defaultsMapSub{One: "one"},
	}

	if !reflect.DeepEqual(expected, val) {
		t.Errorf("Expected %+v, but got %+v", expected, val)
	}
}
