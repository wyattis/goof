package qb_test

import (
	"testing"

	"github.com/wyattis/goof/cmp"
	"github.com/wyattis/goof/qb_gen/qb"
	"github.com/wyattis/goof/qb_gen/test_models"
)

type IToSql interface {
	ToSql() (sql string, params []any, err error)
}

type testStatement struct {
	query  IToSql
	sql    string
	params []any
}

var selectStatements = []testStatement{
	{
		query:  qb.Select.FromUser,
		sql:    "SELECT id, name FROM user",
		params: []any{},
	},
	{
		query:  qb.Select.FromUser.Where().Id("= ?", 1),
		sql:    "SELECT id, name FROM user WHERE id = ?",
		params: []any{1},
	},
	{
		query:  qb.Select.FromUser.Where().Id("=", 1).OrderBy("name"),
		sql:    "SELECT id, name FROM user WHERE id = ? ORDER BY name",
		params: []any{1},
	},
	{
		query:  qb.Select.FromUser.Where().Id("=", 1).OrderBy("name").Limit(1),
		sql:    "SELECT id, name FROM user WHERE id = ? ORDER BY name LIMIT 1",
		params: []any{1},
	},
	{
		query:  qb.Select.FromUser.Where().Id("=", 1).OrderBy("name").Limit(1).Offset(1),
		sql:    "SELECT id, name FROM user WHERE id = ? ORDER BY name LIMIT 1 OFFSET 1",
		params: []any{1},
	},
	{
		query:  qb.Select.FromUser.Where().Id("=", 1).OrderByDesc("name").Limit(1).Offset(1),
		sql:    "SELECT id, name FROM user WHERE id = ? ORDER BY name DESC LIMIT 1 OFFSET 1",
		params: []any{1},
	},
	{
		query:  qb.Select.FromUser.Where().Id("=", 1).OrderBy("name").Limit(1).Offset(1).OrderByDesc("id"),
		sql:    "SELECT id, name FROM user WHERE id = ? ORDER BY name, id DESC LIMIT 1 OFFSET 1",
		params: []any{1},
	},
	{
		query:  qb.Select.FromUser.OrderBy("name").Where.Id("=", 1),
		sql:    "SELECT id, name FROM user WHERE id = ? ORDER BY name",
		params: []any{1},
	},
}

func TestSelectStatements(t *testing.T) {
	for _, s := range selectStatements {
		sql, params, err := s.query.ToSql()
		if err != nil {
			t.Fatal(err)
		}
		if sql != s.sql {
			t.Fatalf("Expected sql '%s', got '%s'", s.sql, sql)
		}
		if !cmp.DeepEqual(params, s.params) {
			t.Fatalf("Expected params '%v', got '%v'", s.params, params)
		}
	}
}

var insertStatements = []testStatement{
	{
		query:  qb.Insert.IntoUser.Values(test_models.User{Name: "Alice"}),
		sql:    "INSERT INTO user (name) VALUES (?)",
		params: []any{"Alice"},
	},
	{
		query:  qb.Insert.IntoUser.Values(test_models.User{Name: "Alice"}).Values(test_models.User{Name: "Bob"}),
		sql:    "INSERT INTO user (name) VALUES (?), (?)",
		params: []any{"Alice", "Bob"},
	},
	{
		query:  qb.Insert.IntoUser.Values(test_models.User{Name: "Alice"}, test_models.User{Name: "Bob"}),
		sql:    "INSERT INTO user (name) VALUES (?), (?)",
		params: []any{"Alice", "Bob"},
	},
}

func TestInsertStatements(t *testing.T) {
	for _, s := range insertStatements {
		sql, params, err := s.query.ToSql()
		if err != nil {
			t.Fatal(err)
		}
		if sql != s.sql {
			t.Fatalf("Expected sql '%s', got '%s'", s.sql, sql)
		}
		if !cmp.DeepEqual(params, s.params) {
			t.Fatalf("Expected params '%v', got '%v'", s.params, params)
		}
	}
}

var updateStatements = []testStatement{
	{
		query:  qb.Update.User.Start().SetName("Alice").Where.Id("=", 1),
		sql:    "UPDATE user SET name = ? WHERE id = ?",
		params: []any{"Alice", 1},
	},
	{
		query:  qb.Update.User.Start().SetName("Alice").Where.Name("=", "Bob"),
		sql:    "UPDATE user SET name = ? WHERE name = ?",
		params: []any{"Alice", "Bob"},
	},
	{
		query:  qb.Update.User.Start().Set(test_models.User{Id: 10, Name: "Alice"}),
		sql:    "UPDATE user SET name = ? WHERE id = ?",
		params: []any{"Alice", 10},
	},
	{
		query:  qb.Update.User.Start().Set(test_models.User{Name: "Alice"}).Where.Id("=", 1),
		sql:    "UPDATE user SET name = ? WHERE id = ?",
		params: []any{"Alice", 1},
	},
	{
		query:  qb.Update.User.Start().Set(test_models.User{Id: 1, Name: "Alice"}).Where.Id("=", 2),
		sql:    "UPDATE user SET name = ? WHERE id = ?",
		params: []any{"Alice", 2},
	},
}

func TestUpdateStatements(t *testing.T) {
	for _, s := range updateStatements {
		sql, params, err := s.query.ToSql()
		if err != nil {
			t.Fatal(err)
		}
		if sql != s.sql {
			t.Fatalf("Expected sql '%s', got '%s'", s.sql, sql)
		}
		if !cmp.DeepEqual(params, s.params) {
			t.Fatalf("Expected params '%v', got '%v'", s.params, params)
		}
	}
}
