package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Routes = map[string]*Route
type RequestExpectation = func(*http.Request) (err error)

type Modifier func(*Route)
type responseModifier func(http.ResponseWriter)

type Route struct {
	header              http.Header
	getBody             func() io.Reader
	requestExpectations []RequestExpectation
	modifiers           []responseModifier
	errors              []error
}

func (r Route) Copy() *Route {
	route := &Route{
		header:  r.header.Clone(),
		getBody: r.getBody,
	}
	for _, expectation := range r.requestExpectations {
		route.requestExpectations = append(route.requestExpectations, expectation)
	}
	for _, modifier := range r.modifiers {
		route.modifiers = append(route.modifiers, modifier)
	}
	return route
}

type Modifiers []Modifier

func Header(key, value string) Modifier {
	return func(r *Route) {
		r.header.Add(key, value)
	}
}

func ExpectHeader(key, value string) Modifier {
	return func(r *Route) {
		r.requestExpectations = append(r.requestExpectations, func(r *http.Request) (err error) {
			if r.Header.Get(key) != value {
				err = fmt.Errorf("expected header %q to be %q; got %q", key, value, r.Header.Get(key))
			}
			return
		})
	}
}

func ExpectContentType(value string) Modifier {
	return ExpectHeader("Content-Type", value)
}

func Status(code int) Modifier {
	return func(r *Route) {
		r.modifiers = append(r.modifiers, func(w http.ResponseWriter) {
			w.WriteHeader(code)
		})
	}
}

func Json(data any, mods ...Modifier) *Route {
	r := Route{
		header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		getBody: func() io.Reader {
			reader, writer := io.Pipe()
			go func() {
				defer writer.Close()
				if err := json.NewEncoder(writer).Encode(data); err != nil {
					panic(err)
				}
			}()
			return reader
		},
	}
	for _, mod := range mods {
		mod(&r)
	}
	return &r
}

func handler(routes map[string]map[string]*Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		route, ok := routes[method][path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("no route for %s %s", method, path)))
			return
		}
		for _, expectation := range route.requestExpectations {
			if err := expectation(r); err != nil {
				route.errors = append(route.errors, err)
			}
		}
		if len(route.errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			for _, err := range route.errors {
				fmt.Fprintln(w, err)
			}
			return
		}
		for _, modifier := range route.modifiers {
			modifier(w)
		}
		for key, values := range route.header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		if route.getBody != nil {
			if _, err := io.Copy(w, route.getBody()); err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func ModifyRoutes(routes Routes, mods ...Modifier) (newRoutes Routes) {
	newRoutes = make(Routes, len(routes))
	for pattern, route := range routes {
		route = route.Copy()
		newRoutes[pattern] = route
		for _, mod := range mods {
			mod(route)
		}
	}
	return newRoutes
}

func NewServer(routes Routes) *httptest.Server {
	methodResponses := map[string]map[string]*Route{
		"GET":    {},
		"POST":   {},
		"PUT":    {},
		"PATCH":  {},
		"DELETE": {},
	}
	for pattern, res := range routes {
		method, pattern, found := strings.Cut(pattern, " ")
		if !found {
			panic("invalid pattern")
		}
		pattern = strings.TrimSpace(pattern)
		// fmt.Println("registering route", method, pattern)
		methodResponses[method][pattern] = res
	}
	return httptest.NewServer(handler(methodResponses))
}
