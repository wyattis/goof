package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:generate go run ./generate.go

//go:embed templates/*
var templates embed.FS

func generateTime() (err error) {
	type format struct {
		Name   string
		Format string
	}
	data := map[string][]format{
		"Formats": {
			{"ANSIC", "time.ANSIC"},
			{"UnixDate", "time.UnixDate"},
			{"RubyDate", "time.RubyDate"},
			{"RFC822", "time.RFC822"},
			{"RFC822Z", "time.RFC822Z"},
			{"RFC850", "time.RFC850"},
			{"RFC1123", "time.RFC1123"},
			{"RFC1123Z", "time.RFC1123Z"},
			{"RFC3339", "time.RFC3339"},
			{"RFC3339Nano", "time.RFC3339Nano"},
			{"Kitchen", "time.Kitchen"},
			{"Stamp", "time.Stamp"},
			{"StampMilli", "time.StampMilli"},
			{"StampMicro", "time.StampMicro"},
			{"StampNano", "time.StampNano"},
			{"DateTime", `"2006-01-02 15:04:05"`},
			{"DateOnly", `"2006-01-02"`},
			{"TimeOnly", `"15:04:05"`},
		},
	}
	t, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		return
	}

	if err = executeToFile(t, "gtime.go", data, "../gtime/gtime.go"); err != nil {
		return
	}
	return executeToFile(t, "gtime_test.go", data, "../gtime/gtime_test.go")
}

func executeToFile(set *template.Template, templateName string, data interface{}, filename string) (err error) {
	t := set.Lookup(templateName)
	if t == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, data); err != nil {
		return
	}

	b := []byte("// Code generated by go generate; DO NOT EDIT.\n")
	b = append(b, buf.Bytes()...)
	return os.WriteFile(filename, b, 0644)
}

func Run() (err error) {
	return generateTime()
}

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}
