package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/wyattis/goof/schema"
)

var (
	ErrSchemaVersionHigherThanTarget = fmt.Errorf("schema version is higher than target version")
	ErrSchemaVersionLowerThanTarget  = fmt.Errorf("schema version is lower than target version")
	ErrNoMigrationForVersion         = fmt.Errorf("no migration for version")
	ErrDatabaseIsDirty               = fmt.Errorf("database is dirty")
)

type SchemaMutator func(s *schema.Schema)

type Migration struct {
	Id   uint
	Up   SchemaMutator
	Down SchemaMutator
}

var Migrations = []Migration{}

func Add(migration Migration) {
	Migrations = append(Migrations, migration)
}

func Begin(db *sql.DB, fn func(tx *sql.Tx) (err error)) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	hasCommitted := false
	defer func() {
		if !hasCommitted {
			fmt.Println("rolling back")
			err2 := tx.Rollback()
			if err2 != nil {
				fmt.Println("rollback err", err2)
			}
		}
	}()
	if err = fn(tx); err != nil {
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	hasCommitted = true
	return
}

func initializeSchema(db *sql.DB, driverType schema.DriverType, name string) (err error) {
	return Begin(db, func(tx *sql.Tx) (err error) {
		s := schema.New(driverType, name)
		s.CreateIfNotExists("schema_migrations", func(t *schema.Table) {
			t.Primary("id")
			t.VarChar("hash", 255).Null()
			t.Boolean("dirty")
			t.Timestamp("started_at").Default(schema.NOW{})
			t.Timestamp("finished_at").Null()
		})
		return s.Schema.Run(tx)
	})
}

func currentVersion(db *sql.DB, driverType schema.DriverType, name string) (version uint, err error) {
	if err = initializeSchema(db, driverType, name); err != nil {
		return
	}
	q := "SELECT id FROM `schema_migrations` ORDER BY id DESC LIMIT 1"
	err = db.QueryRow(q).Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func hasMatchingVersion(migrations []Migration, version uint) bool {
	hasMatchingVersion := false
	for _, m := range migrations {
		if m.Id == version {
			hasMatchingVersion = true
			break
		}
	}
	return hasMatchingVersion
}

func databaseIsClean(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT count(*) FROM `schema_migrations` where dirty").Scan(&count)
	return count == 0 && (err == nil || strings.Contains(err.Error(), "no such table"))
}

func validateMigration(migrations []Migration, db *sql.DB, driver schema.DriverType, name string, version uint) (schemaVersion uint, err error) {
	sort.Slice(Migrations, func(i, j int) bool {
		return Migrations[i].Id < Migrations[j].Id
	})
	if !hasMatchingVersion(migrations, version) {
		err = ErrNoMigrationForVersion
		return
	}
	if !databaseIsClean(db) {
		err = ErrDatabaseIsDirty
		return
	}
	return currentVersion(db, driver, name)
}

func migrateUpTo(migrations []Migration, db *sql.DB, driver schema.DriverType, name string, version uint) (err error) {
	schemaVersion, err := validateMigration(migrations, db, driver, name, version)
	if err != nil {
		return
	}
	if schemaVersion > version {
		err = ErrSchemaVersionHigherThanTarget
		return
	}
	for _, m := range migrations {
		if m.Id > schemaVersion && m.Id <= version {
			err = Begin(db, func(tx *sql.Tx) (err error) {
				// mark current migration as dirty before we start
				q := "INSERT INTO `schema_migrations` (`id`, `dirty`) VALUES (?, ?)"
				_, err = tx.Exec(q, m.Id, true)
				if err != nil {
					return
				}
				s := schema.New(driver, name)
				m.Up(s)
				if err = s.Schema.Run(tx); err != nil {
					return
				}
				q = fmt.Sprintf("UPDATE `schema_migrations` SET `dirty` = ?, `finished_at` = %s WHERE `id` = ?", schema.NOW{}.Constant(driver))
				_, err = tx.Exec(q, false, m.Id)
				return
			})
			if err != nil {
				return
			}
		}
	}
	return
}

func migrateDownTo(migrations []Migration, db *sql.DB, driver schema.DriverType, name string, version uint) (err error) {
	schemaVersion, err := validateMigration(migrations, db, driver, name, version)
	if err != nil {
		return
	}
	if schemaVersion < version {
		err = ErrSchemaVersionLowerThanTarget
		return
	}
	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Id <= schemaVersion && m.Id > version {
			err = Begin(db, func(tx *sql.Tx) (err error) {
				// mark current migration as dirty before we start
				q := "UPDATE `schema_migrations` SET `dirty` = ? WHERE `id` = ?"
				_, err = tx.Exec(q, true, m.Id)
				if err != nil {
					return
				}
				s := schema.New(driver, name)
				m.Down(s)
				if err = s.Schema.Run(tx); err != nil {
					return
				}
				q = "DELETE FROM `schema_migrations` WHERE `id` = ?"
				_, err = tx.Exec(q, m.Id)
				return
			})
			if err != nil {
				return
			}
		}
	}
	return
}

// Migrate up to the provided version. Throws an errors if there isn't a migration matching the provided version, if the
// schema version is higher than the provided version or if the database is dirty (failed a previous migration).
func UpTo(db *sql.DB, driverType schema.DriverType, name string, version uint) (err error) {
	return migrateUpTo(Migrations, db, driverType, name, version)
}
