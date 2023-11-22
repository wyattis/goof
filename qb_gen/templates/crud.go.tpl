{{ define "crud" }}

var CRUD = struct {
  {{ range .Models }}
  {{ .Name }} {{ .BuilderName }}Crud
  {{ end }}
}{}

{{ range .Models }}
type {{ .BuilderName }}Crud struct {}

func (c {{ .BuilderName }}Crud) Create(ctx context.Context, db gsql.IExecContext, m *{{ .ModelName }}) (err error) {
  res, err := Insert.
    Into{{ .Name }}.
    Values(*m).
    ExecContext(ctx, db)
  if err != nil {
    return err
  }
  id, err := res.LastInsertId()
  if err != nil {
    return err
  }
  // TODO: How should we handle multiple primary keys
  m.Id = {{.PrimaryField.TypeStr}}(id)
  return
}

func (c {{ .BuilderName }}Crud) Get(ctx context.Context, db gsql.IQueryRowContext, id int64) (m *{{ .ModelName }}, err error) {
  m = &{{ .ModelName }}{}
  err = Select.
    From{{ .Name }}.
    Where().Id("=", id).
    GetScanContext(ctx, m, db)
  if err != nil {
    return nil, err
  }
  return
}

func (c {{ .BuilderName }}Crud) Update(ctx context.Context, db gsql.IExecContext, id int64, m *{{ .ModelName }}) (err error) {
  {{ with .PrimaryField }}
  if {{ checkNonZero .Type "m." .Name }} && int64(m.{{ .Name }}) != id {
    return errors.New("id mismatch during update")
  }
  {{ end }}
  _, err = Update.
    {{ .Name }}.
    Start().
    Set(*m).
    Where.
    Id("=", id).
    ExecContext(ctx, db)
  if err != nil {
    return
  }
  m.{{ .PrimaryField.Name }} = {{ .PrimaryField.TypeStr }}(id)
  return
}

func (c {{ .BuilderName }}Crud) Delete(ctx context.Context, db gsql.IExecContext, id int64) (err error) {
/* TODO: implement delete function */
  /* _, err = Delete.
    From{{ .Name }}.
    Where().
    Id("=", id).
    ExecContext(ctx, db) */
  return
}

{{ end }}
{{ end }}