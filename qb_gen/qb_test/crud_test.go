package qb_test

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/cmp"
	"github.com/wyattis/goof/crud"
	"github.com/wyattis/goof/gsql"
	"github.com/wyattis/goof/gsql/driver"
	"github.com/wyattis/goof/qb_gen/qb"
	"github.com/wyattis/goof/qb_gen/test_models"
	"github.com/wyattis/goof/route/route_gin"
	"github.com/wyattis/goof/test"
	"github.com/wyattis/z/zset"
)

func openTestDb(t *testing.T) *sql.DB {
	db, err := gsql.Open(driver.Config{Driver: driver.TypeSqlite3, Database: ":memory:"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCrudRoutes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	routes := crud.Routes[test_models.User](&sql.DB{}, qb.CRUD.User, crud.Config{
		Name: "user",
		Mode: crud.ModeAll,
	})
	route_gin.Mount(router, routes)

	expected := zset.New("GET /user/:id", "POST /user", "PUT /user/:id", "DELETE /user/:id", "GET /user")
	registered := zset.New[string]()
	for _, r := range router.Routes() {
		registered.Add(r.Method + " " + r.Path)
	}

	if !expected.Equal(*registered) {
		t.Errorf("expected %v, got %v", expected, registered)
	}
}

func TestCrudServer(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	db := openTestDb(t)
	router := gin.New()
	routes := crud.Routes[test_models.User](db, qb.CRUD.User, crud.Config{
		Name: "user",
		Mode: crud.ModeAll,
	})
	router.Use(func(c *gin.Context) {
		c.Next()
	})
	route_gin.Mount(router, routes)

	s := httptest.NewServer(router)
	defer s.Close()

	var post = test_models.User{
		Name: "test",
	}
	if err := test.PostJson(s, "/user", post, &post); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if post.Id != 1 {
		t.Errorf("expected id 1, got %v", post.Id)
	}
	if post.Name != "test" {
		t.Errorf("expected name test, got %v", post.Name)
	}

	var get test_models.User
	if err := test.GetJson(s, "/user/1", &get); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if get.Id != 1 {
		t.Errorf("expected id 1, got %v", get.Id)
	}
	if get.Name != "test" {
		t.Errorf("expected name test, got %v", get.Name)
	}

	var put = test_models.User{
		Name: "test2",
	}
	if err := test.PutJson(s, "/user/1", put, &put); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if put.Id != 1 {
		t.Errorf("expected id 1, got %v", put.Id)
	}
	if put.Name != "test2" {
		t.Errorf("expected name test2, got %v", put.Name)
	}

	var second test_models.User
	if err := test.PostJson(s, "/user", test_models.User{Name: "test3"}, &second); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var list []test_models.User
	if err := test.GetJson(s, "/page/user", &list); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 users, got %v", len(list))
	}
	expectedList := []test_models.User{put, second}
	if !cmp.DeepEqual(list, expectedList) {
		t.Errorf("expected users %v, got %v", expectedList, list)
	}

	if err := test.Delete(s, "/user/1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
