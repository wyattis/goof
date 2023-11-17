package cmp

import (
	"database/sql/driver"
	"reflect"
)

// A deep equal implentation that uses the driver.Valuer interface for equality if it's available.
func DeepEqual(a, b any) bool {
	av, aIsV := a.(driver.Valuer)
	bv, bIsV := b.(driver.Valuer)
	if aIsV && bIsV {
		aVal, _ := av.Value()
		bVal, _ := bv.Value()
		return reflect.DeepEqual(aVal, bVal)
	}

	at := reflect.TypeOf(a)
	switch at.Kind() {
	case reflect.Struct:
		return deepEqualStruct(a, b)
	case reflect.Slice:
		return deepEqualSlice(a, b)
	case reflect.Map:
		return deepEqualMap(a, b)
	default:
		return reflect.DeepEqual(a, b)
	}

}

func deepEqualMap(a, b any) bool {
	am := reflect.ValueOf(a)
	bm := reflect.ValueOf(b)

	if am.Len() != bm.Len() {
		return false
	}

	for _, key := range am.MapKeys() {
		av := am.MapIndex(key)
		bv := bm.MapIndex(key)
		if !DeepEqual(av.Interface(), bv.Interface()) {
			return false
		}
	}
	return true
}

func deepEqualSlice(a, b any) bool {
	as := reflect.ValueOf(a)
	bs := reflect.ValueOf(b)

	if as.Len() != bs.Len() {
		return false
	}

	for i := 0; i < as.Len(); i++ {
		if !DeepEqual(as.Index(i).Interface(), bs.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func deepEqualStruct(a, b any) bool {
	at := reflect.TypeOf(a)

	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < at.NumField(); i++ {
		af := at.Field(i)
		if af.PkgPath != "" {
			continue // skip unexported fields
		}
		avf := av.Field(i)
		bvf := bv.Field(i)
		if !DeepEqual(avf.Interface(), bvf.Interface()) {
			return false
		}
	}
	return true
}
