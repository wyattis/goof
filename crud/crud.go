package crud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/goof"
	"github.com/wyattis/goof/gsql"
	"github.com/wyattis/goof/route"
)

type CRUD[T any] interface {
	Create(ctx context.Context, db gsql.IExecContext, v *T) (err error)
	Get(ctx context.Context, db gsql.IQueryRowContext, id int64) (v *T, err error)
	GetPage(ctx context.Context, db gsql.IQueryContext, page, pageSize int64, orderBy string, desc bool) (v []T, err error)
	Update(ctx context.Context, db gsql.IExecContext, id int64, v *T) (err error)
	Delete(ctx context.Context, db gsql.IExecContext, id int64) (err error)
}

const ModeGet uint8 = 1 << 0
const ModeCreate uint8 = 1 << 1
const ModeUpdate uint8 = 1 << 2
const ModeDelete uint8 = 1 << 3
const ModeList uint8 = 1 << 4

const ModeGetCreate = ModeGet | ModeCreate
const ModeAll = ModeGet | ModeCreate | ModeUpdate | ModeDelete | ModeList
const ModeRead = ModeGet | ModeList
const ModeWrite = ModeCreate | ModeUpdate | ModeDelete
const ModeCreateDelete = ModeCreate | ModeDelete
const ModeAny = ModeAll

// Check if the mode matches the mask
func ModeMatches(mode, mask uint8) bool {
	return mode&mask != 0
}

// Basically, this defines the configuration API that we are committing to because the CRUD interface can be changed
// simultaneously by qb_gen and goof.
type Config struct {
	Name            string
	Mode            uint8
	DefaultPageSize int64
}

type idUri struct {
	Id int64 `uri:"id" binding:"required"`
}

type pageQuery struct {
	Page    int64  `form:"page"`
	Size    int64  `form:"size"`
	OrderBy string `form:"orderBy"`
	Desc    bool   `form:"desc"`
}

// Create a set of routes for the given crud interface
func Routes[T any](db gsql.IDBContext, crud CRUD[T], config Config) route.IRoute {
	if config.Name == "" {
		panic(fmt.Sprintf("must define a name for the CRUD routes: %T", crud))
	}
	name := config.Name
	if config.DefaultPageSize == 0 {
		config.DefaultPageSize = 25
	}

	if config.Mode == 0 {
		config.Mode = ModeAll
	}

	routes := []route.IRoute{}
	// if mode allows Get
	if ModeMatches(config.Mode, ModeGet) {
		routes = append(routes, Get[T](fmt.Sprintf("/%s/:id", name), db, crud))
	}

	if ModeMatches(config.Mode, ModeCreate) {
		routes = append(routes, Create[T](fmt.Sprintf("/%s", name), db, crud))
	}

	if ModeMatches(config.Mode, ModeUpdate) {
		routes = append(routes, Update[T](fmt.Sprintf("/%s/:id", name), db, crud))
	}

	if ModeMatches(config.Mode, ModeDelete) {
		routes = append(routes, Delete[T](fmt.Sprintf("/%s/:id", name), db, crud))
	}

	if ModeMatches(config.Mode, ModeList) {
		routes = append(routes, Page[T](fmt.Sprintf("/%s", name), db, crud))
	}
	return route.Group("", routes...)
}

// Register a CREATE route for the given crud interface
func Create[T any](pattern string, db gsql.IDBContext, crud CRUD[T]) route.IRoute {
	return goof.Json(pattern, func(c *gin.Context, payload T) (res T, status int, err error) {
		if err = crud.Create(c, db, &payload); err != nil {
			return
		}
		res = payload
		return
	}).Post()
}

// Register an UPDATE route for the given crud interface
func Update[T any](pattern string, db gsql.IDBContext, crud CRUD[T]) route.IRoute {
	return goof.Json(pattern, func(c *gin.Context, payload T) (res T, status int, err error) {
		var u idUri
		if err := c.ShouldBindUri(&u); err != nil {
			return res, http.StatusBadRequest, err
		}
		// This will fail if the id is present in the payload and it doesn't match the provided id
		if err = crud.Update(c, db, u.Id, &payload); err != nil {
			return
		}
		res = payload
		return
	}).Put()
}

// Register a DELETE route for the given crud interface
func Delete[T any](pattern string, db gsql.IDBContext, crud CRUD[T]) route.IRoute {
	return goof.Status(pattern, func(c *gin.Context) (status int, err error) {
		var u idUri
		if err := c.ShouldBindUri(&u); err != nil {
			return http.StatusBadRequest, err
		}
		err = crud.Delete(c, db, u.Id)
		return
	}).Delete()
}

// Register a GET route for the given crud interface
func Get[T any](pattern string, db gsql.IDBContext, crud CRUD[T]) route.IRoute {
	return goof.ToJson(pattern, func(c *gin.Context) (res T, status int, err error) {
		var u idUri
		if err := c.ShouldBindUri(&u); err != nil {
			return res, http.StatusBadRequest, err
		}
		r, err := crud.Get(c, db, u.Id)
		if err != nil {
			return res, http.StatusInternalServerError, err
		}
		res = *r
		return
	}).Get()
}

func Page[T any](pattern string, db gsql.IDBContext, crud CRUD[T]) route.IRoute {
	return goof.ToJson(pattern, func(c *gin.Context) (res []T, status int, err error) {
		var query pageQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			return res, http.StatusBadRequest, err
		}
		if query.Size == 0 {
			query.Size = 25
		}
		res, err = crud.GetPage(c, db, query.Page, query.Size, query.OrderBy, query.Desc)
		return
	}).Get()
}
