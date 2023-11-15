package migrate

import (
	"os"
	"testing"

	"github.com/wyattis/goof/schema"
	"github.com/wyattis/goof/sql/driver"

	_ "github.com/mattn/go-sqlite3"
)

var userCommentMigrations = []Migration{
	{
		Version: 1,
		Up: func(s *schema.Schema) {
			s.Create("user", func(t *schema.Table) {
				t.Primary("id")
				t.VarChar("name", 255).Null()
				t.Integer("age").Null()
				t.Boolean("active").Default(true)
				t.Timestamp("created_at").Default(schema.NOW{})
				t.Timestamp("updated_at").Null()
			})
		},
		Down: func(s *schema.Schema) {
			s.Drop("user")
		},
	},
	{
		Version: 2,
		Up: func(s *schema.Schema) {
			s.Create("comment", func(t *schema.Table) {
				t.Primary("id")
				t.Integer("user_id").References("user", "id")
				t.Text("body").Null()
				t.Timestamp("created_at").Default(schema.NOW{})
				t.Timestamp("updated_at").Null()
			})
		},
		Down: func(s *schema.Schema) {
			// s.DropForeign("user_id")
			s.Drop("comment")
		},
	},
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSqliteUp(t *testing.T) {
	// sqliteDbFile = "test.db"
	db, err := setupSqlite()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := MigrateUpTo(userCommentMigrations, db, driver.TypeSqlite3, "test", 1); err != nil {
		t.Fatalf("Failed the first migration to version 1: %s", err)
	}
	if err := MigrateUpTo(userCommentMigrations, db, driver.TypeSqlite3, "test", 1); err != nil {
		t.Fatalf("Failed the second call to version 1. Migration should be idempotent: %s", err)
	}
	if err := MigrateUpTo(userCommentMigrations, db, driver.TypeSqlite3, "test", 2); err != nil {
		t.Fatalf("Failed the first call to version 2: %s", err)
	}
	if err := MigrateUpTo(userCommentMigrations, db, driver.TypeSqlite3, "test", 2); err != nil {
		t.Fatalf("Failed the second call to version 2. Migration should be idempotent: %s", err)
	}
}

func TestSqliteDown(t *testing.T) {
	// sqliteDbFile = "test.db"
	db, err := setupSqlite()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := MigrateUpTo(userCommentMigrations, db, driver.TypeSqlite3, "test", 2); err != nil {
		t.Error("Failed to migrate up", err)
	}
	if err := migrateDownTo(userCommentMigrations, db, driver.TypeSqlite3, "test", 1); err != nil {
		t.Error("Failed to migrate down", err)
	}
}
