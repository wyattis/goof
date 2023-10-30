package schema

import (
	"database/sql"
	"fmt"
)

type TableMutator func(t *Table)

func New(driver DriverType, name string) *Schema {
	return &Schema{
		Schema: &SchemaDef{
			Driver: driver,
			Name:   name,
		},
	}
}

type Statement struct {
	Sql    string
	Params []interface{}
}

type SchemaDef struct {
	Driver          DriverType
	Name            string
	Tables          []*TableDef
	Execs           []Statement
	DroppingTables  []string
	DroppingForeign []string
	DroppingIndices []string
	DropCreated     bool
}

type Schema struct {
	Schema *SchemaDef
}

func (s *Schema) Exec(statement string, params ...interface{}) {
	s.Schema.Execs = append(s.Schema.Execs, Statement{Sql: statement, Params: params})
}

func (s *Schema) Create(table string, fn TableMutator) {
	builder := Table{
		tableDef: &TableDef{
			Schema:     s.Schema,
			Name:       table,
			WillCreate: true,
		},
	}
	s.Schema.Tables = append(s.Schema.Tables, builder.tableDef)
	fn(&builder)
	return
}

func (s *Schema) CreateIfNotExists(table string, fn TableMutator) {
	s.Create(table, func(t *Table) {
		t.tableDef.IfNotExists = true
		fn(t)
	})
}

func (s *Schema) Table(table string, fn TableMutator) {
	builder := Table{
		tableDef: &TableDef{
			Schema: s.Schema,
			Name:   table,
		},
	}
	s.Schema.Tables = append(s.Schema.Tables, builder.tableDef)
	fn(&builder)
}

func (s *Schema) Drop(table string) {
	s.Schema.DroppingTables = append(s.Schema.DroppingTables, table)
}

func (s *Schema) DropForeign(table string, foreign string) {
	s.Schema.DroppingForeign = append(s.Schema.DroppingForeign, fmt.Sprintf("fk_%s_%s", table, foreign))
}

func (s *Schema) DropIndex(name string) {
	s.Schema.DroppingIndices = append(s.Schema.DroppingIndices, name)
}

func (s *Schema) DropCreated() {
	s.Schema.DropCreated = true
}

func (s *SchemaDef) Statements() (statements []string) {
	for _, table := range s.Tables {
		statements = append(statements, table.Statements()...)
	}
	statements = append(statements, s.DropStatements()...)
	return
}

func (s *SchemaDef) DropStatements() (statements []string) {
	if s.DropCreated {
		for _, table := range s.Tables {
			statements = append(statements, fmt.Sprintf("DROP TABLE `%s`", table.Name))
			// TODO: drop indices and foreign keys
		}
	} else {
		for _, index := range s.DroppingIndices {
			statements = append(statements, fmt.Sprintf("DROP INDEX `%s`", index))
		}
		for _, foreign := range s.DroppingForeign {
			statements = append(statements, fmt.Sprintf("DROP FOREIGN KEY `%s`", foreign))
		}
		for _, table := range s.DroppingTables {
			statements = append(statements, fmt.Sprintf("DROP TABLE `%s`", table))
		}
	}
	return
}

func (s *SchemaDef) Run(tx *sql.Tx) (err error) {
	statements := s.Statements()
	for _, statement := range statements {
		_, err = tx.Exec(statement)
		if err != nil {
			return
		}
	}
	return
}
