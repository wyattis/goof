package conf

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/wyattis/goof/log"
	"github.com/wyattis/z/zreflect"
)

type envConfigurer struct {
	filepaths     []string
	prefix        string
	fileMustExist bool
	onlyFiles     bool
}

func (e *envConfigurer) Init(val interface{}) (err error) {
	return
}

func (e *envConfigurer) SetPrefix(prefix string) {
	e.prefix = prefix
}

func (e *envConfigurer) Apply(val interface{}, args ...string) (err error) {
	vals := make(map[string]string)
	if err := e.loadFiles(vals); err != nil {
		return err
	}
	if !e.onlyFiles {
		if err := e.loadEnv(vals); err != nil {
			return err
		}
	}

	it := zreflect.FieldIterator(val)
	for it.Next() {
		if !it.IsStructField() {
			continue
		}
		v := it.Value()
		field := it.Field()
		path := it.Path()
		key := it.Key()
		name := field.Tag.Get("env")
		if name == "" {
			if e.prefix != "" {
				path = append([]string{e.prefix}, path...)
			}
			name = strings.ToUpper(strings.Join(append(path, key), "_"))
		} else {
			name = strings.ToUpper(name)
		}
		if strVal, ok := vals[name]; ok {
			log.Trace().Strs("path", append(path, key)).Str("env", name).Str("value", strVal).Msg("Setting value from environment")
			if err = setValFromStr(v, field, strVal); err != nil {
				return
			}
		}
	}
	return
}

func (e *envConfigurer) loadFiles(vals map[string]string) (err error) {
	for _, path := range e.filepaths {
		log.Trace().Str("path", path).Msg("Loading environment file")
		f, err := os.Open(path)
		if err != nil {
			if !e.fileMustExist && os.IsNotExist(err) {
				continue
			}
			return err
		}
		defer f.Close()
		m, err := godotenv.Parse(f)
		log.Trace().Str("path", path).Err(err).Msg("Loaded environment file")
		if err != nil {
			return err
		}
		for k, v := range m {
			vals[k] = v
		}
	}
	return
}

// Load environment variables into the given map.
func (e *envConfigurer) loadEnv(vals map[string]string) (err error) {
	for _, v := range os.Environ() {
		key, val, found := strings.Cut(v, "=")
		if !found {
			continue
		}
		vals[key] = val
	}
	return
}
