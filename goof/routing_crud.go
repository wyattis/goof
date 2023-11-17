package goof

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/wyattis/goof/log"
	"github.com/wyattis/z/zslice/zstrings"
	"github.com/wyattis/z/zstring"
)

type field struct {
	name      string
	visible   bool
	updatable bool
	dbName    string
	jsonName  string
}

type CrudOpts struct {
	Table string

	Get    bool
	List   bool
	Create bool
	Update bool
	Delete bool

	PathName        string
	VisibleFields   []string
	UpdatableFields []string
}

// Create a CRUD controller. DO NOT USE
func CRUD[T any](db *sqlx.DB, model T, opts *CrudOpts) *crud[T] {
	if opts == nil {
		opts = &CrudOpts{
			Get:    true,
			List:   true,
			Create: true,
			Update: true,
			Delete: true,
		}
	}

	if opts.Table == "" {
		opts.Table = getTableName[T](model)
	}

	if opts.PathName == "" {
		opts.PathName = getRouteName[T](model)
	}

	if len(opts.VisibleFields) == 0 {
		opts.VisibleFields = getJsonColumns(model)
	}

	if len(opts.UpdatableFields) == 0 {
		opts.UpdatableFields = getDbColumns(model)
	}

	var fields []field

	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ff := field{
			name:     f.Name,
			dbName:   f.Tag.Get("db"),
			jsonName: f.Tag.Get("json"),
		}
		ff.visible = zstrings.Contains(opts.VisibleFields, ff.jsonName)
		ff.updatable = zstrings.Contains(opts.UpdatableFields, ff.dbName)
		fields = append(fields, ff)
	}

	return &crud[T]{
		opts:   *opts,
		model:  model,
		fields: fields,
		db:     db,
	}
}

type crud[T any] struct {
	opts   CrudOpts
	model  T
	fields []field
	db     *sqlx.DB
}

func (c *crud[T]) Routes() (routes []IRoute) {
	if c.opts.Get {
		routes = append(routes, c.getRoute().Routes()...)
	}
	if c.opts.List {
		routes = append(routes, c.listRoute().Routes()...)
	}
	if c.opts.Create {
		routes = append(routes, c.createRoute().Routes()...)
	}
	// TODO: update and delete
	return
}

// Standard GET route for a CRUD controller
func (c *crud[T]) getRoute() Routable {
	pattern := fmt.Sprintf("/%s/:id", c.opts.PathName)
	cols := strings.Join(c.visibleColumns(), ",")
	q := fmt.Sprintf("SELECT %s FROM `%s` WHERE id = ?", cols, c.opts.Table)
	log.Debug().Str("sql", q).Msg("get route sql")
	getStmt, err := c.db.Preparex(q)
	if err != nil {
		panic(err)
	}
	return Json(pattern, func(ctx *gin.Context, payload T) (res T, status int, err error) {
		id := ctx.Param("id")
		if id == "" {
			err = fmt.Errorf("id is required")
			status = http.StatusBadRequest
			return
		}
		err = getStmt.Get(&res, id)
		return
	})
}

// Standard LIST route for a CRUD controller
// TODO: this should be deprecated in favor of pagination
func (c *crud[T]) listRoute() Routable {
	pattern := fmt.Sprintf("/%s", c.opts.PathName)
	cols := strings.Join(c.visibleColumns(), ",")
	q := fmt.Sprintf("SELECT %s FROM `%s`", cols, c.opts.Table)
	log.Debug().Str("sql", q).Msg("list route sql")
	listStmt, err := c.db.Preparex(q)
	if err != nil {
		panic(err)
	}
	return Json(pattern, func(ctx *gin.Context, payload T) (res []T, status int, err error) {
		err = listStmt.Select(&res)
		return
	})
}

// Standard CREATE route for a CRUD controller
func (c *crud[T]) createRoute() Routable {
	pattern := fmt.Sprintf("/%s", c.opts.PathName)
	cols := strings.Join(c.updatableColumns(), ",")
	namedPlaceholder := strings.Join(c.updatableColumns(), ",")
	q := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", c.opts.Table, cols, namedPlaceholder)
	createStmt, err := c.db.PrepareNamed(q)
	log.Debug().Str("sql", q).Msg("create route sql")
	if err != nil {
		panic(err)
	}
	return Json(pattern, func(ctx *gin.Context, payload T) (res T, status int, err error) {
		err = createStmt.Get(&res, payload)
		return
	})
}

func (c *crud[T]) visibleColumns() (cols []string) {
	for _, f := range c.fields {
		if f.visible {
			cols = append(cols, f.dbName)
		}
	}
	return
}

func (c *crud[T]) updatableColumns() (cols []string) {
	for _, f := range c.fields {
		if f.updatable {
			cols = append(cols, f.dbName)
		}
	}
	return
}

// Get the name of a type formatted for a path
func getRouteName[T any](m T) string {
	t := reflect.TypeOf(m)
	structName := t.Name()
	return zstring.CamelToSnake(structName, "-", 2)
}

// Get the name of a type formatted for a database table
func getTableName[T any](m T) string {
	t := reflect.TypeOf(m)
	structName := t.Name()
	return zstring.CamelToSnake(structName, "_", 2)
}

// Get all struct tags of a given type
func getStructTags[T any](m T, tag string) []string {
	t := reflect.TypeOf(m)
	var tags []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags = append(tags, field.Tag.Get(tag))
	}
	return tags
}

// Get all db columns
func getDbColumns[T any](m T) (columns []string) {
	return getStructTags[T](m, "db")
}

// Get all json columns
func getJsonColumns[T any](m T) (columns []string) {
	return getStructTags[T](m, "json")
}

// Format column names as sqlx named placeholders
func getNamedPlaceholders(columns []string) (placeholders []string) {
	for _, col := range columns {
		placeholders = append(placeholders, fmt.Sprintf(":%s", col))
	}
	return
}
