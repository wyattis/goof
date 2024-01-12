package conf

import (
	"errors"
	"reflect"
	"time"

	"github.com/wyattis/goof/log"
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
	return forEachField(val, func(path []string, key string, field reflect.StructField, v reflect.Value) (err error) {
		defaultTag := field.Tag.Get("default")
		if v.CanInterface() {
			f := v
			if f.CanAddr() {
				f = f.Addr()
			}
			if set, ok := f.Interface().(ConfigSettable); ok {
				if err = set.SetConfig(defaultTag); err != nil {
					return
				}
				return errDontDescend
			}
		}
		if defaultTag != "" {
			log.Trace().Str("field", field.Name).Str("value", defaultTag).Msg("Setting default value")
			if err = setValFromStr(v, field, defaultTag); err != nil {
				return
			}
		} else if field.Type == reflect.TypeOf(time.Time{}) {
			// Don't descend into time.Time fields
			return errDontDescend
		}
		return
	})
}
