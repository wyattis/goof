{{- define "qb" -}}

{{ template "types" }}

{{ range .Models }}
{{ template "select-builder" . }}
{{ template "select-builder-factory" . }}
{{ template "insert" . }}
{{ template "insert-factory" . }}
{{ template "update" . }}
{{ template "update-factory" . }}
{{ end }}


// Define our exported Select builders
var Select struct {
  {{ range .Models -}}
  // Build a select query for the model {{ .Name }}
  From{{ .Name }} {{ .BuilderName }}SelectBuilderFactory
  {{ end }}
} = struct {
  {{ range .Models -}}
  From{{ .Name }} {{ .BuilderName }}SelectBuilderFactory
  {{ end }}
}{}

// Define our export Insert builders
var Insert struct {
  {{ range .Models -}}
  // Build an insert query for the model {{ .Name }}
  Into{{ .Name }} {{ .BuilderName }}InsertFactory
  {{ end }}
} = struct {
  {{ range .Models -}}
  Into{{ .Name }} {{ .BuilderName }}InsertFactory
  {{ end }}
}{}

// Define our exported Update builders
var Update struct {
  {{ range .Models -}}
  // Build an update query for the model {{ .Name }}
  {{ .Name }} {{ .BuilderName }}UpdateFactory
  {{ end }}
} = struct {
  {{ range .Models -}}
  {{ .Name }} {{ .BuilderName }}UpdateFactory
  {{ end }}
}{}

{{ end }}


{{ define "update-factory" }}

type {{ .BuilderName }}UpdateFactory struct {}

func (f {{ .BuilderName }}UpdateFactory) Start() *{{ .BuilderName }}Update {
  return new{{ .BuilderName }}Update()
}

{{ end }}

{{ define "update" }}

{{ template "update-where" . -}}

var default{{ .Name }}UpdateFields = []string{
  {{ range .Fields -}}
  {{ if .ShouldUpdate }}"{{ .ColName }}",{{ end }}
  {{- end }}
}

func new{{ .BuilderName }}Update() *{{ .BuilderName }}Update {
  u := &{{ .BuilderName }}Update{ }
  u.Where = &{{ .BuilderName }}UpdateWhere{ parent: u }
  return u
}

type {{ .BuilderName }}Update struct {
  value {{ .ModelName }}
  setCols []string
  params []interface{}
  setZeros bool
  Where *{{ .BuilderName }}UpdateWhere
}

{{ range .Fields -}}
{{- if .ShouldUpdate}}
// Set the value of the field {{ .Name }}. This will be included in the update query.
func (b *{{ $.BuilderName }}Update) Set{{.Name}}(val {{ .TypeStr }}) *{{ $.BuilderName }}Update {
  b.params = append(b.params, val)
  b.setCols = append(b.setCols, "{{ .ColName }}")
  return b
}
{{ end -}}
{{- end }}

// Set all values on the model using the provided model. This will be included in the update query.
func (b *{{ .BuilderName }}Update) Set(model {{ .ModelName }}) *{{ .BuilderName }}Update {
  b.value = model
  return b
}

// Zero values will be included in the update query. This is useful for updating fields that are not nullable.
func (b *{{ .BuilderName }}Update) IncludeZeros() *{{ .BuilderName }}Update {
  b.setZeros = true
  return b
}

// Run the update query.
func (b {{ .BuilderName }}Update) Exec(db gsql.IExec) (sql.Result, error) {
  sql, params, err := b.ToSql()
  if err != nil {
    return nil, err
  }
  return db.Exec(sql, params...)
}

// Same as Exec, but with a context.
func (b {{ .BuilderName }}Update) ExecContext(ctx context.Context, db gsql.IExecContext) (sql.Result, error) {
  sql, params, err := b.ToSql()
  if err != nil {
    return nil, err
  }
  return db.ExecContext(ctx, sql, params...)
}

// Convert the update statement to parameterized SQL and parameters
func (b {{ .BuilderName }}Update) ToSql() (sql string, params []interface{}, err error) {
  sql = "UPDATE {{ .TableName }} SET "
  
  // If Set<Field> was called, use those fields.
  if len(b.setCols) > 0 {
    sql += strings.Join(b.setCols, " = ?, ") + " = ?"
    params = b.params
  } else if b.setZeros {
    // Update all updateable fields, even if they are zero
    sql += `{{ range $i, $f := .UpdateFields -}}
    {{ if $i }}, {{ end -}}
    {{ $f.ColName }} = ?
    {{- end }}`
    params = []interface{}{
      {{ range .UpdateFields -}}
      b.value.{{ .Name }},
      {{- end }}
    }
  } else {
    // Only update our non-zero fields
    updates := []string{}
    {{ range $i, $f := .UpdateFields }}
    if {{ checkNonZero .Type "b.value." .Name }} {
      updates = append(updates, "{{ .ColName }} = ?")
      params = append(params, b.value.{{ .Name }})
    }
    {{- end }}
    sql += strings.Join(updates, ", ")
  }

  if len(b.Where.expressions) > 0 {
    sql += " WHERE "
    for i, e := range b.Where.expressions {
      if i > 0 {
        sql += " AND "
      }
      sql += e.expr
      params = append(params, e.params...)
    }
  } else {
    {{ range .PrimaryKeys -}}
    if !({{ checkNonZero .Type "b.value." .Name }}) {
      return "", nil, errors.New("Primary key of `{{ .Name }}` must be non-zero to update. Otherwise, use Where() to specify which rows to update")
    }
    {{- end }}
    sql += ` WHERE 
    {{- range $i, $f := .PrimaryKeys -}}
    {{ if $i }} AND{{ end }} {{ $f.ColName }} = ?
    {{- end -}}
    `
    params = append(params, {{ range $i, $f := .PrimaryKeys -}}
      {{ if $i }}, {{ end -}}
      b.value.{{ $f.Name }}
      {{- end -}}
    )
  }
  return
}

{{- end }}


{{ define "insert-factory" }}

type {{ .BuilderName }}InsertFactory struct {}

func (f {{ .BuilderName }}InsertFactory) Values(values ...{{.ModelName}}) *{{ .BuilderName }}Insert {
  return new{{ .BuilderName }}Insert().Values(values...)
}

{{ end }}

{{ define "insert" }}

func new{{ .BuilderName }}Insert() *{{ .BuilderName }}Insert {
  return &{{ .BuilderName }}Insert{ }
}

type {{ .BuilderName }}Insert struct {
  values []{{ .ModelName }}
}

// Define the values to insert with this query
func (b *{{ .BuilderName }}Insert) Values(values ...{{ .ModelName }}) *{{ .BuilderName }}Insert {
  b.values = append(b.values, values...)
  return b
}

// Insert a single model into the database.
func (b {{ .BuilderName }}Insert) Exec(db gsql.IExec) (sql.Result, error) {
  sql, params, err := b.ToSql()
  if err != nil {
    return nil, err
  }
  return db.Exec(sql, params...)
}

// Same as Exec, but with a context.
func (b {{ .BuilderName }}Insert) ExecContext(ctx context.Context, db gsql.IExecContext) (sql.Result, error) {
  sql, params, err := b.ToSql()
  if err != nil {
    return nil, err
  }
  return db.ExecContext(ctx, sql, params...)
}

// Returns the parameterized SQL and parameters for the query.
func (b {{ .BuilderName }}Insert) ToSql() (sql string, params []interface{}, err error) {
  sql = "INSERT INTO {{ .TableName }} ("
  var cols []string
  var vals []string
  {{ range .Fields -}}
  {{ if .ShouldInsert }}
  cols = append(cols, "{{ .ColName }}")
  vals = append(vals, "?")
  {{ end }}
  {{- end }}
  sql += strings.Join(cols, ", ") + ") VALUES "
  valsStr := "(" + strings.Join(vals, ", ") + ")"
  for i, v := range b.values {
    if i > 0 {
      sql += ", "
    }
    sql += valsStr
    params = append(params,
      {{ range .Fields -}}
      {{ if .ShouldInsert }}v.{{ .Name }},{{ end }}
      {{- end }}
    )
  }
  return
}

{{ end }}



{{ define "select-builder" }}

{{ template "select-where" .}}

// The list of fields that will be selected by default
var {{ .BuilderName }}DefaultSelectFields = []string{
  {{ range .Fields -}}
  {{ if .ShouldScan }}"{{ .ColName }}",{{ end }}
  {{- end }}
}

// The list of fields that will be inserted by default
var {{ .BuilderName }}DefaultInsertFields = []string{
  {{ range .Fields -}}
  {{ if .ShouldInsert }}"{{ .ColName }}",{{ end }}
  {{- end }}
}

// The list of fields that will be updated by default
var {{ .BuilderName }}DefaultUpdateFields = []string{
  {{ range .Fields -}}
  {{ if .ShouldUpdate }}"{{ .ColName }}",{{ end }}
  {{- end }}
}

// Our select query builder for the model {{ .ModelName }}
type {{ .BuilderName }}SelectBuilder struct {
  columns []string
  from string
  orderBy []order
  limit int
  offset int

  Where *{{ .BuilderName }}SelectWhere
}

func new{{ .BuilderName }}SelectBuilder() *{{ .BuilderName }}SelectBuilder {
  b := &{{ .BuilderName }}SelectBuilder{ }
  b.Where = &{{ .BuilderName }}SelectWhere{ parent: b }
  return b
}

// Define the columns to select with this query
func (b *{{ .BuilderName }}SelectBuilder) Columns(columns ...string) *{{ .BuilderName }}SelectBuilder {
  b.columns = columns
  return b
}

// Define the table to select from. If not provided, this will be inferred from the TableName() method on the model if 
// it exists, otherwise it will be inferred from the model name.
func (b *{{ .BuilderName }}SelectBuilder) From(from string) *{{ .BuilderName }}SelectBuilder {
  b.from = from
  return b
}

// Order by the given column in ascending order.
func (b *{{ .BuilderName }}SelectBuilder) OrderBy(column string) *{{ .BuilderName }}SelectBuilder {
  b.orderBy = append(b.orderBy, order{column, false})
  return b
}

// Order by the given column in descending order.
func (b *{{ .BuilderName }}SelectBuilder) OrderByDesc(column string) *{{ .BuilderName }}SelectBuilder {
  b.orderBy = append(b.orderBy, order{column, true})
  return b
}

// Limit the number of results returned by the query.
func (b *{{ .BuilderName }}SelectBuilder) Limit(limit int) *{{ .BuilderName }}SelectBuilder {
  b.limit = limit
  return b
}

// Offset the results returned by the query.
func (b *{{ .BuilderName }}SelectBuilder) Offset(offset int) *{{ .BuilderName }}SelectBuilder {
  b.offset = offset
  return b
}

// Returns a single model matching the query. If the query returns multiple rows, only the first will be returned.
func (b *{{ .BuilderName }}SelectBuilder) Get(db gsql.IQueryRow) (model *{{ .ModelName }}, err error) {
  model = &{{ .ModelName }}{}
  err = b.GetScan(model, db)
  return 
}

// Same as Get, but with a context.
func (b *{{ .BuilderName }}SelectBuilder) GetContext(ctx context.Context, db gsql.IQueryRowContext) (model *{{ .ModelName }}, err error) {
  model = &{{ .ModelName }}{}
  if err = b.GetScanContext(ctx, model, db); err != nil {
    return nil, err
  }
  return
}

// Same as Get, but scans the result into the provided model.
func (b *{{ .BuilderName }}SelectBuilder) GetScan(model *{{ .ModelName }}, db gsql.IQueryRow) (err error) {
  sql, params, err := b.ToSql()
  if err != nil {
    return
  }
  err = db.QueryRow(sql, params...).Scan({{ .ScanFields "model" }})
  return
}

// Same as GetScan, but with a context.
func (b *{{ .BuilderName }}SelectBuilder) GetScanContext(ctx context.Context, model *{{ .ModelName }}, db gsql.IQueryRowContext) (err error) {
  sql, params, err := b.ToSql()
  if err != nil {
    return
  }
  err = db.QueryRowContext(ctx, sql, params...).Scan({{ .ScanFields "model" }})
  return
}

// Returns the parameterized SQL and parameters for the query.
func (b {{ .BuilderName }}SelectBuilder) ToSql() (sql string, params []interface{}, err error) {
  sql = "SELECT "
  if len(b.columns) == 0 {
    sql += strings.Join({{ .BuilderName }}DefaultSelectFields, ", ")
  } else {
    sql += strings.Join(b.columns, ", ")
  }
  if b.from == "" {
    b.from = "{{ .TableName }}"
  }
  sql += " FROM " + b.from
  if len(b.Where.expressions) > 0 {
    sql += " WHERE "
    for i, e := range b.Where.expressions {
      if i > 0 {
        sql += " AND "
      }
      sql += e.expr
      params = append(params, e.params...)
    }
  }
  if len(b.orderBy) > 0 {
    sql += " ORDER BY "
    for i, o := range b.orderBy {
      if i > 0 {
        sql += ", "
      }
      sql += o.column
      if o.desc {
        sql += " DESC"
      }
    }
  }
  if b.limit > 0 {
    sql += fmt.Sprintf(" LIMIT %d", b.limit)
  }
  if b.offset > 0 {
    sql += fmt.Sprintf(" OFFSET %d", b.offset)
  }
  return
}

{{ end }}

{{ define "update-where" }}
type {{ .BuilderName }}UpdateWhere struct {
  expressions []expression
  parent *{{ .BuilderName }}Update
}

{{ range .Fields }}
func (w *{{ $.BuilderName }}UpdateWhere) {{ .Name }}(expr string, params ...any) *{{ $.BuilderName }}UpdateWhere {
  expr = fmt.Sprintf("{{ .ColName }} " + expr)
  if !strings.Contains(expr, "?") {
    expr += " ?"
  }
  w.expressions = append(w.expressions, expression{expr, params})
  return w
}
{{ end }}

func (w {{ $.BuilderName }}UpdateWhere) Exec(db gsql.IExec) (sql.Result, error) {
  return w.parent.Exec(db)
}

func (w {{ $.BuilderName }}UpdateWhere) ExecContext(ctx context.Context, db gsql.IExecContext) (sql.Result, error) {
  return w.parent.ExecContext(ctx, db)
}

func (w {{ $.BuilderName }}UpdateWhere) ToSql() (sql string, params []interface{}, err error) {
  return w.parent.ToSql()
}

{{ end }}

{{ define "select-where" }}

type {{ .BuilderName }}SelectWhere struct {
  expressions []expression
  parent *{{ .BuilderName }}SelectBuilder
}

{{ range .Fields }}
func (w *{{ $.BuilderName }}SelectWhere) {{ .Name }}(expr string, params ...any) *{{ $.BuilderName }}SelectWhere {
  expr = fmt.Sprintf("{{ .ColName }} " + expr)
  if !strings.Contains(expr, "?") {
    expr += " ?"
  }
  w.expressions = append(w.expressions, expression{expr, params})
  return w
}
{{ end }}

func (w {{ .BuilderName }}SelectWhere) OrderBy(column string) *{{ .BuilderName }}SelectBuilder {
  return w.parent.OrderBy(column)
}

func (w {{ .BuilderName }}SelectWhere) OrderByDesc(column string) *{{ .BuilderName }}SelectBuilder {
  return w.parent.OrderByDesc(column)
}

func (w {{ .BuilderName }}SelectWhere) Limit(limit int) *{{ .BuilderName }}SelectBuilder {
  return w.parent.Limit(limit)
}

func (w {{ .BuilderName }}SelectWhere) Offset(offset int) *{{ .BuilderName }}SelectBuilder {
  return w.parent.Offset(offset)
}

func (w {{ .BuilderName }}SelectWhere) ToSql() (sql string, params []interface{}, err error) {
  return w.parent.ToSql()
}

// Returns a single model matching the query. If the query returns multiple rows, only the first will be returned.
func (w {{ .BuilderName }}SelectWhere) Get(db gsql.IQueryRow) (model *{{ .ModelName }}, err error) {
  return w.parent.Get(db)
}

// Same as Get, but with a context.
func (w {{ .BuilderName }}SelectWhere) GetContext(ctx context.Context, db gsql.IQueryRowContext) (model *{{ .ModelName }}, err error) {
  return w.parent.GetContext(ctx, db)
}

// Same as Get, but scans the result into the provided model.
func (w {{ .BuilderName }}SelectWhere) GetScan(model *{{ .ModelName }}, db gsql.IQueryRow) (err error) {
  return w.parent.GetScan(model, db)
}

// Same as GetScan, but with a context.
func (w {{ .BuilderName }}SelectWhere) GetScanContext(ctx context.Context, model *{{ .ModelName }}, db gsql.IQueryRowContext) (err error) {
  return w.parent.GetScanContext(ctx, model, db)
}

{{ end }}



{{ define "select-builder-factory" }}
type {{ .BuilderName }}SelectBuilderFactory struct {}

func (f {{ .BuilderName }}SelectBuilderFactory) Columns(columns ...string) *{{ .BuilderName }}SelectBuilder {
  return new{{ .BuilderName }}SelectBuilder().Columns(columns...)
}

// Define the table to select from. If not provided, this will be inferred from the TableName() method on the model if 
// it exists, otherwise it will be inferred from the model name.
func (b *{{ .BuilderName }}SelectBuilderFactory) From(from string) *{{ .BuilderName }}SelectBuilder {
  return new{{ .BuilderName }}SelectBuilder().From(from)
}

// Apply any where conditions for the query. Example usage: Select.From{{ .Name }}.Where("id = ?", id)
func (b *{{ .BuilderName }}SelectBuilderFactory) Where() *{{ .BuilderName }}SelectWhere {
  return new{{ .BuilderName }}SelectBuilder().Where
}

// Order by the given column in ascending order.
func (b *{{ .BuilderName }}SelectBuilderFactory) OrderBy(column string) *{{ .BuilderName }}SelectBuilder {
  return new{{ .BuilderName }}SelectBuilder().OrderBy(column)
}

// Order by the given column in descending order.
func (b *{{ .BuilderName }}SelectBuilderFactory) OrderByDesc(column string) *{{ .BuilderName }}SelectBuilder {
  return new{{ .BuilderName }}SelectBuilder().OrderByDesc(column)
}

// Limit the number of results returned by the query.
func (b *{{ .BuilderName }}SelectBuilderFactory) Limit(limit int) *{{ .BuilderName }}SelectBuilder {
  return new{{ .BuilderName }}SelectBuilder().Limit(limit)
}

// Offset the results returned by the query.
func (b *{{ .BuilderName }}SelectBuilderFactory) Offset(offset int) *{{ .BuilderName }}SelectBuilder {
  return new{{ .BuilderName }}SelectBuilder().Offset(offset)
}

// Returns the parameterized SQL and parameters for the query.
func (b {{ .BuilderName }}SelectBuilderFactory) ToSql() (sql string, params []interface{}, err error) {
  return new{{ .BuilderName }}SelectBuilder().ToSql()
}

// Returns a single model matching the query. If the query returns multiple rows, only the first will be returned.
func (b {{ .BuilderName }}SelectBuilderFactory) Get(db gsql.IQueryRow) (model *{{ .ModelName }}, err error) {
  return new{{ .BuilderName }}SelectBuilder().Get(db)
}

// Same as Get, but with a context.
func (b {{ .BuilderName }}SelectBuilderFactory) GetContext(ctx context.Context, db gsql.IQueryRowContext) (model *{{ .ModelName }}, err error) {
  return new{{ .BuilderName }}SelectBuilder().GetContext(ctx, db)
}

// Same as Get, but scans the result into the provided model.
func (b {{ .BuilderName }}SelectBuilderFactory) GetScan(model *{{ .ModelName }}, db gsql.IQueryRow) (err error) {
  return new{{ .BuilderName }}SelectBuilder().GetScan(model, db)
}

// Same as GetScan, but with a context.
func (b {{ .BuilderName }}SelectBuilderFactory) GetScanContext(ctx context.Context, model *{{ .ModelName }}, db gsql.IQueryRowContext) (err error) {
  return new{{ .BuilderName }}SelectBuilder().GetScanContext(ctx, model, db)
}
{{ end }}




{{ define "types" }}
type expression struct {
	expr   string
	params []any
}

type order struct {
  column string
  desc   bool
}

{{ end }}