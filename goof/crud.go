package goof

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/gsql"
)

type CRUD[T any] interface {
	Create(ctx context.Context, db gsql.IExecContext, v *T) (err error)
	Get(ctx context.Context, db gsql.IQueryRowContext, id int64) (v *T, err error)
	Update(ctx context.Context, db gsql.IExecContext, id int64, v *T) (err error)
	Delete(ctx context.Context, db gsql.IExecContext, id int64) (err error)
}

type routeGroup struct {
	routes []Routable
}

func (r *routeGroup) Routes() (res []IRoute) {
	for _, route := range r.routes {
		res = append(res, route.Routes()...)
	}
	return
}

// Create a group of routes
func Group(routes ...Routable) Routable {
	return &routeGroup{
		routes: routes,
	}
}

// Basically, this defines the configuration API that we are committing to because the CRUD interface can be changed
// simultaneously by qb_gen and goof.
type CrudConfig struct {
	Name   string
	All    bool
	Get    bool
	Create bool
	Update bool
	Delete bool
	List   bool
}

type idUri struct {
	Id int64 `uri:"id" binding:"required"`
}

// Create a group of routes for a CRUD interface
func CrudRoutes[T any](db gsql.IDBContext, crud CRUD[T], config CrudConfig) Routable {
	name := config.Name
	// TODO: how to customize the uri?

	routes := []Routable{}
	if config.All || config.Get {
		routes = append(routes,
			ToJson(fmt.Sprintf("/%s/:id", name), func(c *gin.Context) (res T, status int, err error) {
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
			}).Get(),
		)
	}

	if config.All || config.Create {
		routes = append(routes,
			Json(fmt.Sprintf("/%s", name), func(c *gin.Context, payload T) (res T, status int, err error) {
				if err = crud.Create(c, db, &payload); err != nil {
					return
				}
				res = payload
				return
			}).Post(),
		)
	}

	if config.All || config.Update {
		routes = append(routes,
			Json(fmt.Sprintf("/%s/:id", name), func(c *gin.Context, payload T) (res T, status int, err error) {
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
			}).Put(),
		)
	}

	if config.All || config.Delete {
		routes = append(routes,
			Status(fmt.Sprintf("/%s/:id", name), func(c *gin.Context) (status int, err error) {
				var u idUri
				if err := c.ShouldBindUri(&u); err != nil {
					return http.StatusBadRequest, err
				}
				err = crud.Delete(c, db, u.Id)
				return
			}).Delete(),
		)
	}

	return Group(routes...)
}
