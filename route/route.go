package route

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wyattis/z/zset"
)

type IRoute interface {
	Route() Route
}

type HandlerFunc func(*gin.Context) (err error)

func newRoute(pattern string) *Route {
	return &Route{Pattern: pattern, Methods: zset.New[string]()}
}

type Route struct {
	Uses        []Middleware
	Methods     *zset.Set[string]
	Pattern     string
	Name        string
	Description string
	Handler     HandlerFunc
	Children    []IRoute
	IsGroup     bool
}

func (r Route) IsAny() bool {
	return r.Methods.Equal(*anyMethods)
}

func NewRouteBuilder(pattern string) *RouteBuilder {
	return &RouteBuilder{route: newRoute(pattern)}
}

type RouteBuilder struct {
	route *Route
}

func (r *RouteBuilder) Route() Route {
	return *r.route
}

func (r *RouteBuilder) Any() *RouteBuilder {
	r.route.Methods = anyMethods.Clone()
	return r
}

func (r *RouteBuilder) Get() *RouteBuilder {
	r.route.Methods.Add(http.MethodGet)
	return r
}

func (r *RouteBuilder) Post() *RouteBuilder {
	r.route.Methods.Add(http.MethodPost)
	return r
}

func (r *RouteBuilder) Put() *RouteBuilder {
	r.route.Methods.Add(http.MethodPut)
	return r
}

func (r *RouteBuilder) Delete() *RouteBuilder {
	r.route.Methods.Add(http.MethodDelete)
	return r
}

func (r *RouteBuilder) Patch() *RouteBuilder {
	r.route.Methods.Add(http.MethodPatch)
	return r
}

func (r *RouteBuilder) Head() *RouteBuilder {
	r.route.Methods.Add(http.MethodHead)
	return r
}

func (r *RouteBuilder) Options() *RouteBuilder {
	r.route.Methods.Add(http.MethodOptions)
	return r
}

func (r *RouteBuilder) Trace() *RouteBuilder {
	r.route.Methods.Add(http.MethodTrace)
	return r
}

func (r *RouteBuilder) Use(middlewares ...Middleware) *RouteBuilder {
	r.route.Uses = append(r.route.Uses, middlewares...)
	return r
}

func (r *RouteBuilder) Handle(handler HandlerFunc) *RouteBuilder {
	r.route.Handler = handler
	return r
}

func (r *RouteBuilder) Name(name string) *RouteBuilder {
	r.route.Name = name
	return r
}

var anyMethods = zset.New(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodTrace)

func R(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern)
}

func Any(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Any()
}

func Get(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Get()
}

func Post(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Post()
}

func Put(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Put()
}

func Delete(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Delete()
}

func Patch(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Patch()
}

func Head(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Head()
}

func Options(pattern string) *RouteBuilder {
	return NewRouteBuilder(pattern).Options()
}

type group struct {
	prefix     string
	name       string
	middleware []Middleware
	children   []IRoute
}

type groupBuilder struct {
	group *group
}

func (g *groupBuilder) Use(middlewares ...Middleware) *groupBuilder {
	g.group.middleware = append(g.group.middleware, middlewares...)
	return g
}

func (g *groupBuilder) Routes(children ...IRoute) *groupBuilder {
	g.group.children = append(g.group.children, children...)
	return g
}

func (g *groupBuilder) Name(name string) *groupBuilder {
	g.group.name = name
	return g
}

func (g *groupBuilder) Route() Route {
	return Route{
		Name:     g.group.name,
		Pattern:  g.group.prefix,
		Uses:     g.group.middleware,
		Children: g.group.children,
		IsGroup:  true,
	}
}

// Create a new group of routes with the given prefix and middleware
func Group(prefix string, children ...IRoute) *groupBuilder {
	return &groupBuilder{group: &group{prefix: prefix, children: children}}
}

// Print a list of routes to the given writer. TODO: this should handle grouping of routes
func PrintRoutes(writer io.Writer, routes ...IRoute) (err error) {
	return printRouteGroup(writer, 0, routes...)
}

func writeCols(writer io.Writer, linePrefix string, cols [][]string) (err error) {
	colWidths := make([]int, len(cols[0]))
	for _, col := range cols {
		for i, val := range col {
			if len(val) > colWidths[i] {
				colWidths[i] = len(val)
			}
		}
	}
	for _, col := range cols {
		for i, val := range col {
			if _, err = fmt.Fprintf(writer, "%s%-*s  ", linePrefix, colWidths[i], val); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(writer); err != nil {
			return
		}
	}
	return
}

func printRouteGroup(writer io.Writer, indent int, routes ...IRoute) (err error) {
	indentStr := strings.Repeat("  ", indent)
	cols := [][]string{}
	for _, r := range routes {
		route := r.Route()
		uses := make([]string, len(route.Uses))
		for i, use := range route.Uses {
			uses[i] = use.Name()
		}
		usesStr := strings.Join(uses, "->")
		if usesStr != "" {
			usesStr = usesStr + "->"
		}
		usesStr += route.Name
		if route.IsGroup {
			if err = writeCols(writer, indentStr, [][]string{{route.Pattern, "", usesStr, "{"}}); err != nil {
				return
			}
			if err = printRouteGroup(writer, indent+1, route.Children...); err != nil {
				return
			}
			if _, err = fmt.Fprintf(writer, "%s%s\n", indentStr, "}"); err != nil {
				return
			}
			continue
		} else {
			methods := route.Methods.Items()
			sort.Strings(methods)
			methodStr := strings.Join(methods, ",")
			cols = append(cols, []string{route.Pattern, methodStr, usesStr})
		}
	}
	if len(cols) > 0 {
		return writeCols(writer, indentStr, cols)
	}
	return
}

// TODO: In theory it's possible to generate a TypeScript API, right?
func TSRouter(routes ...IRoute) {
	// for _, route := range routes {
	// 	routes := route.Routes()
	// 	for _, r := range routes {
	// 		println(r.Methods(), r.Pattern(), reflect.TypeOf(route))
	// 	}
	// }
}
