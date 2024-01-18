package conf

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/wyattis/goof/log"
	"github.com/wyattis/z/zreflect"
)

type defaultConfigurer struct{}

func (d *defaultConfigurer) Init(val interface{}) (err error) {
	return
}

// Iterate over the value and set any fields that have a default tag
func (d *defaultConfigurer) Apply(val interface{}, args ...string) (err error) {
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}
	it := zreflect.FieldIterator(val)
	for it.Next() {
		v := it.Value()
		fmt.Println(it.Type(), it.Path(), it.Key(), v.String(), it.IsStruct())
		if !it.IsStruct() {
			continue
		}
		field := it.Field()
		defaultTag := field.Tag.Get("default")
		if defaultTag == "" {
			continue
		}
		if v.CanInterface() {
			f := v
			if f.CanAddr() {
				f = f.Addr()
			}
			if set, ok := f.Interface().(ConfigSettable); ok {
				if err = set.SetConfig(defaultTag); err != nil {
					return
				}
				it.DontDescend()
			}
		}
		log.Trace().Str("field", field.Name).Str("value", defaultTag).
			Bool("canSet", v.CanSet()).
			Bool("canAddr", v.CanAddr()).
			Bool("canInterface", v.CanInterface()).
			Str("type", v.Type().String()).
			Msg("Setting default value")
		k := v.Kind()
		newVal, err := getValFromStr(v, field, defaultTag)
		if err != nil {
			return err
		}
		if k == reflect.Ptr {
			it.Set(newVal)
		} else {
			it.Set(newVal.Elem())
		}

		// if err = setValFromStr(v, field, defaultTag); err != nil {
		// 	return
		// }
		if field.Type == reflect.TypeOf(time.Time{}) {
			it.DontDescend()
		}
	}
	return
}
