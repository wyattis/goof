// Simple struct-based configuration library
// Supports config from environment variables, .env files, command line flags, config files (yaml, json, toml, etc) and
// default tags.
package conf

import (
	"errors"
	"flag"
	"reflect"
	"strconv"
	"time"

	"github.com/wyattis/z/ztime"
)

type ConfigSettable interface {
	SetConfig(val interface{}) error
}

type Configurer interface {
	Init(val interface{}) error
	Apply(val interface{}, args ...string) error
}

type configOption = func(c *zconfig) error

type zconfig struct {
	configurers   []Configurer
	isInitialized bool
}

// Create a configurer. Configurer precedence is the order in which they are provided.
func New(opts ...configOption) *zconfig {
	z := &zconfig{}
	for _, opt := range opts {
		if err := opt(z); err != nil {
			panic(err)
		}
	}
	return z
}

func (c *zconfig) Init(val interface{}) (err error) {
	c.isInitialized = true
	for _, opt := range c.configurers {
		if err = opt.Init(val); err != nil {
			return
		}
	}
	return
}

// Actually update the struct given the values provided
func (c *zconfig) Apply(val interface{}, args ...string) (err error) {
	if !c.isInitialized {
		if err = c.Init(val); err != nil {
			return
		}
	}
	if len(c.configurers) == 0 {
		if err = Auto()(c); err != nil {
			return
		}
	}

	// Apply configurers in reverse order
	revConfigurers := make([]Configurer, len(c.configurers))
	for i := 0; i < len(c.configurers); i++ {
		revConfigurers[i] = c.configurers[len(c.configurers)-i-1]
	}
	for _, opt := range revConfigurers {
		if err = opt.Apply(val, args...); err != nil {
			return
		}
	}
	return
}

// Helper function to quickly load a struct using the configurers provided
func Load(val interface{}, opts ...configOption) (err error) {
	return New(opts...).Apply(val)
}

// Configure a struct using any command line flags, environment variables,
// variables in .env and finally, values set using default tags
func Auto() configOption {
	return func(c *zconfig) error {
		Flag()(c)
		Env()(c)
		Defaults()(c)
		return nil
	}
}

func Defaults() configOption {
	return func(c *zconfig) error {
		c.configurers = append(c.configurers, &defaultConfigurer{})
		return nil
	}
}

// This will load configuration values from both environment variables and any .env files provided. By default, if a
// provided path does not exist it will be skipped. Keys are set using the struct field name, but can be overridden
// using the `env` struct tag. If the `env` struct tag is not set, the key is the struct field name converted to
// uppercase with words separated by underscores. Additional underscores are added for nested structs.
func Env(paths ...string) configOption {
	return func(c *zconfig) error {
		c.configurers = append(c.configurers, &envConfigurer{})
		return nil
	}
}

// Load configuration values from any .env files provided. By default, if a provided path does not exist it will be
// skipped. The keys use the same format and struct tag as the Env configurer.
func EnvFiles(paths ...string) configOption {
	return func(c *zconfig) error {
		c.configurers = append(c.configurers, &envConfigurer{filepaths: paths, onlyFiles: true})
		return nil
	}
}

// Load configuration values from command line flags. The flag name can be overridden using the `flag` struct tag. The
// default flag name is the struct field name converted to kebab-case. Additional dashes are added for nested structs.
func Flag(opts ...flagOpt) configOption {
	return func(c *zconfig) error {
		configurer := &FlagConfigurer{}
		for _, opt := range opts {
			if err := opt(configurer); err != nil {
				return err
			}
		}
		c.configurers = append(c.configurers, configurer)
		return nil
	}
}

type flagOpt = func(c *FlagConfigurer) error

// Provide a custom flag set. This is useful if you want to parse the flags yourself or are using an external library.
func UseFlagSet(fs *flag.FlagSet) flagOpt {
	return func(c *FlagConfigurer) error {
		c.flagSet = fs
		return nil
	}
}

// Tell the Flag configurer not to parse the flag set. This is useful if you want to parse the flags yourself or are
// using an external library.
func SkipParse() flagOpt {
	return func(c *FlagConfigurer) error {
		c.dontParseFlagSet = true
		return nil
	}
}

var errDontDescend = errors.New("don't descend into field")

type fieldHandler = func(path []string, key string, field reflect.StructField, value reflect.Value) (err error)

func forEachField(v any, cb fieldHandler) (err error) {
	return recursiveIterFields(nil, reflect.Indirect(reflect.ValueOf(v)), cb)
}

func recursiveIterFields(path []string, val reflect.Value, cb fieldHandler) (err error) {
	k := val.Kind()
	if k == reflect.Ptr {
		val = val.Elem()
	}
	isIterable := k == reflect.Struct || k == reflect.Map
	if !isIterable {
		return
	}
	typ := val.Type()
	if k == reflect.Map {
		for _, key := range val.MapKeys() {
			value := val.MapIndex(key)
			if err = cb(path, key.String(), nil, value); err != nil {
				if err == errDontDescend {
					err = nil
					continue
				}
				return
			}
			k := value.Type().Kind()
			if k == reflect.Struct || k == reflect.Map {
				if err = recursiveIterFields(append(path, key.String()), reflect.Indirect(value), cb); err != nil {
					return
				}
			}
		}
	} else {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			value := val.Field(i)
			if err = cb(path, field.Name, field, value); err != nil {
				if err == errDontDescend {
					err = nil
					continue
				}
				return
			}
			k := field.Type.Kind()
			if k == reflect.Struct || k == reflect.Map {
				if err = recursiveIterFields(append(path, field.Name), reflect.Indirect(value), cb); err != nil {
					return
				}
			}
		}
	}
	return
}

// This parses a string value and applies it to the reflected Value (reflect.Value)
func setValFromStr(v reflect.Value, fieldType reflect.StructField, strVal string) (err error) {
	k := v.Kind()
	t := v.Type()

	isPointer := k == reflect.Ptr
	if isPointer {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		k = v.Kind()
		t = v.Type()
	}

	if t == reflect.TypeOf(time.Time{}) {
		return setTimeVal(v, fieldType.Tag, strVal)
	} else if t == reflect.TypeOf(time.Duration(0)) {
		return setDurationVal(v, strVal)
	}

	switch k {
	case reflect.String:
		v.SetString(strVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return err
		}
		v.SetFloat(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(strVal)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return errors.New("unsupported type")
	}
	return
}

func setTimeVal(dest reflect.Value, tag reflect.StructTag, strVal string) (err error) {
	timeFormat := tag.Get("time-format")
	var t time.Time
	if timeFormat == "" {
		t, err = ztime.Parse(strVal)
	} else {
		t, err = ztime.Parse(strVal, timeFormat)
	}
	if err != nil {
		return err
	}
	dest.Set(reflect.ValueOf(t))
	return nil
}

func setDurationVal(dest reflect.Value, strVal string) (err error) {
	d, err := time.ParseDuration(strVal)
	if err != nil {
		return err
	}
	dest.Set(reflect.ValueOf(d))
	return nil
}
