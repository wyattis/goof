package schema

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

type TableDef struct {
	Schema       *SchemaDef
	OriginalName string
	Name         string
	WillCreate   bool
	IfNotExists  bool
	Columns      []*columnDef
	Indices      []*indexDef
}

type Table struct {
	tableDef *TableDef
}

// Create an index on this table using the given columns
func (t *Table) Index(cols ...string) *indexBuilder {
	i := &indexDef{
		Table:   t.tableDef,
		Columns: cols,
	}
	t.tableDef.Indices = append(t.tableDef.Indices, i)
	return &indexBuilder{i}
}

// Create a unique index on this table using the given columns
func (t *Table) Unique(cols ...string) *indexBuilder {
	i := &indexDef{
		Table:   t.tableDef,
		Unique:  true,
		Columns: cols,
	}
	t.tableDef.Indices = append(t.tableDef.Indices, i)
	return &indexBuilder{i}
}

// Create a (or modify) column on this table
func (t *Table) Column(name string, mods ...ColumnMod) *columnBuilder {
	c := &columnDef{
		table:        t,
		OriginalName: name,
		Name:         name,
		Kind:         TypeVarChar,
		KindLen:      255,
		IsNull:       false, // default to not null for all drivers
	}
	t.tableDef.Columns = append(t.tableDef.Columns, c)
	builder := columnBuilder{c}
	return builder.applyMods(mods...)
}

func (t *Table) Primary(name string, mods ...ColumnMod) *columnBuilder {
	m := []ColumnMod{Primary(), Integer()}
	return t.Column(name, append(m, mods...)...)
}

func (t *Table) BigInt(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{BigInt()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) String(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{VarChar(255)}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) DateTime(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Datetime()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) Timestamp(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Timestamp()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) Integer(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Integer()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) VarChar(name string, n int, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{VarChar(n)}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) Text(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Text()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) NVarChar(name string, n int, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{NVarChar(n)}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) Json(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Json()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) Enum(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Enum()}, mods...)
	return t.Column(name, mods...)
}

func (t *Table) Boolean(name string, mods ...ColumnMod) *columnBuilder {
	mods = append([]ColumnMod{Boolean()}, mods...)
	return t.Column(name, mods...)
}

func (t *TableDef) NumPrimary() int {
	res := 0
	for _, c := range t.Columns {
		if c.IsPrimary {
			res++
		}
	}
	return res
}

func (t *TableDef) Statements() (statements []string) {
	if t.WillCreate {
		statements = append(statements, t.createStatement())
	} else {
		statements = append(statements, t.alterStatement())
	}
	statements = append(statements, t.indexStatements()...)
	return
}

var defaultTypeMap = map[ColumnType]string{
	TypeVarChar:   "VARCHAR",
	TypeNVarChar:  "NVARCHAR",
	TypeText:      "TEXT",
	TypeJson:      "JSON",
	TypeInteger:   "INTEGER",
	TypeDateTime:  "DATETIME",
	TypeEnum:      "ENUM",
	TypeBoolean:   "BOOLEAN",
	TypeBigInt:    "BIGINT",
	TypeDecimal:   "DECIMAL",
	TypeTinyInt:   "TINYINT",
	TypeFloat:     "FLOAT",
	TypeDate:      "DATE",
	TypeTime:      "TIME",
	TypeTimestamp: "TIMESTAMP",
}

var sqliteTypeMap = map[ColumnType]string{}
var mysqlTypeMap = map[ColumnType]string{}
var postgresTypeMap = map[ColumnType]string{}

//go:embed templates/*
var templates embed.FS

func sqliteFuncMap() template.FuncMap {
	typeMap := map[ColumnType]string{}
	for k, v := range defaultTypeMap {
		typeMap[k] = v
	}
	for k, v := range sqliteTypeMap {
		typeMap[k] = v
	}
	return template.FuncMap{
		"GetType": func(kind ColumnType, num int) string {
			res := typeMap[kind]
			if res == "VARCHAR" || res == "NVARCHAR" {
				res += fmt.Sprintf("(%d)", num)
			}
			return res
		},
		"GetDefault": func(kind ColumnType, val interface{}) string {
			if val == nil {
				return ""
			}
			switch kind {
			case TypeVarChar, TypeNVarChar, TypeText, TypeJson, TypeEnum:
				return fmt.Sprintf(" DEFAULT '%s'", val)
			case TypeInteger, TypeBigInt, TypeDecimal, TypeTinyInt, TypeFloat:
				return fmt.Sprintf(" DEFAULT %v", val)
			case TypeBoolean:
				if val.(bool) {
					return " DEFAULT TRUE"
				} else {
					return " DEFAULT FALSE"
				}
			case TypeDateTime, TypeDate, TypeTime, TypeTimestamp:
				c, ok := val.(Constant)
				if !ok {
					return fmt.Sprintf(" DEFAULT '%s'", val)
				}
				return fmt.Sprintf(" DEFAULT %s", c.Constant(DriverTypeSqlite3))
			default:
				return ""
			}
		},
		"join": strings.Join,
	}
}

func (t *TableDef) loadTemplates() (tmp *template.Template) {
	dirFs, err := fs.Sub(templates, "templates")
	if err != nil {
		panic(err)
	}
	var funcMap template.FuncMap
	switch t.Schema.Driver {
	case DriverTypeMysql:
		// TODO
	case DriverTypePostgres:
		// TODO
	case DriverTypeSqlite3:
		funcMap = sqliteFuncMap()
	default:
		panic("unknown driver type")
	}
	tmp = template.New("table").Funcs(funcMap)
	tmpName := fmt.Sprintf("%s.tpl", t.Schema.Driver)
	tmp, err = tmp.ParseFS(dirFs, tmpName)
	if err != nil {
		panic(err)
	}
	return
}

func (t *TableDef) createStatement() (s string) {
	tmp := t.loadTemplates()
	res := bytes.Buffer{}
	if err := tmp.ExecuteTemplate(&res, "create_table", t); err != nil {
		panic(err)
	}
	s = res.String()
	return
}

func (t *TableDef) alterStatement() string {
	return ""
}

func (t *TableDef) indexStatements() (statements []string) {
	tmp := t.loadTemplates()
	for _, idx := range t.Indices {
		res := bytes.Buffer{}
		if err := tmp.ExecuteTemplate(&res, "create_index", idx); err != nil {
			panic(err)
		}
		statements = append(statements, res.String())
	}
	return
}
