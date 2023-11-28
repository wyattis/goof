{{ define "crud" }}

var CRUD = struct {
  {{ range .Models }}
  {{ .Name }} {{ .BuilderName }}Crud
  {{ end }}
}{}

{{ range .Models }}
type Uri{{ .Name }}Crud struct {
  Id int64 `uri:"{{.PrimaryField.UriName}}" binding:"required,min=1"`
}

type {{ .BuilderName }}Crud struct {}

func (c {{.BuilderName}}Crud) URI() any {
  return Uri{{.Name}}Crud{}
}

func (c {{ .BuilderName }}Crud) Create(ctx context.Context, db gsql.IExecQueryRowContext, m *{{ .ModelName }}) (err error) {
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
  err = Select.From{{ .Name }}.Where().Id("=", {{.PrimaryField.TypeStr}}(id)).GetScanContext(ctx, db, m)
  return
}

func (c {{ .BuilderName }}Crud) Get(ctx context.Context, db gsql.IQueryRowContext, id int64) (m *{{ .ModelName }}, err error) {
  m = &{{ .ModelName }}{}
  err = Select.
    From{{ .Name }}.
    Where().Id("=", id).
    GetScanContext(ctx, db, m)
  if err != nil {
    return nil, err
  }
  return
}

func (c {{ .BuilderName }}Crud) GetPage(ctx context.Context, db gsql.IQueryContext, page, size int, orderBy string, desc bool) (list []{{ .ModelName }}, err error) {
  q := Select.
    From{{ .Name }}.
    Offset(int(page * size)).
    Limit(int(size))
  if orderBy != "" {
    if desc {
      q = q.OrderByDesc(orderBy)
    } else {
      q = q.OrderBy(orderBy)
    }
  }
  return q.ExecContext(ctx, db)
}

func (c {{ .BuilderName }}Crud) Update(ctx context.Context, db gsql.IExecContext, id int64, m *{{ .ModelName }}) (err error) {
  {{ with .PrimaryField }}
  if !({{ isZero .Type "m." .Name }}) && int64(m.{{ .Name }}) != id {
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

// CRUD ROUTES

// Get route expects a pattern containing ":{{.PrimaryField.UriName}}"
func (c {{.BuilderName}}Crud) GetRoute(pattern string, db gsql.IQueryRowContext) route.IRoute {
  if !strings.Contains(pattern, ":{{.PrimaryField.UriName}}") {
    panic("{{.Name}}.GetRoute expects a pattern containing ':{{.PrimaryField.UriName}}'")
  }
  return goof.ToJson[*{{ .ModelName }}](pattern, func(ctx *gin.Context) (res *{{ .ModelName }}, status int, err error) {
    var uri Uri{{.Name}}Crud
    if err = ctx.ShouldBindUri(&uri); err != nil {
      return
    }
    res, err = c.Get(ctx, db, uri.Id)
    return
  }).Get().Name("Return a {{ .ModelName }} matching :{{.PrimaryField.UriName}}")
}

// Any pattern can be used with CreateRoute
func (c {{.BuilderName}}Crud) CreateRoute(pattern string, db gsql.IExecQueryRowContext) route.IRoute {
  return goof.Json[{{ .ModelName }}](pattern, func(ctx *gin.Context, payload {{ .ModelName }}) (res {{ .ModelName }}, status int, err error) {
    // TODO(future): How could we handle multiple primary keys
    payload.{{ .PrimaryField.Name }} = {{ .PrimaryField.TypeStr }}(0) 
    if err = c.Create(ctx, db, &payload); err != nil {
      return
    }
    res = payload
    return
  }).Post().Name("Create a {{ .ModelName }}")
}

// Update route expects a pattern containing ":{{.PrimaryField.UriName}}"
func (c {{.BuilderName}}Crud) UpdateRoute(pattern string, db gsql.IExecContext) route.IRoute {
  if !strings.Contains(pattern, ":{{.PrimaryField.UriName}}") {
    panic("{{.Name}}.UpdateRoute expects a pattern containing ':{{.PrimaryField.UriName}}'")
  }
  return goof.Json[{{ .ModelName }}](pattern, func(ctx *gin.Context, payload {{ .ModelName }}) (res {{ .ModelName }}, status int, err error) {
    var uri Uri{{.Name}}Crud
    if err = ctx.ShouldBindUri(&uri); err != nil {
      return
    }
    // TODO(future): How could we handle multiple primary keys
    payload.{{ .PrimaryField.Name }} = {{ .PrimaryField.TypeStr }}(0) 
    if err = c.Update(ctx, db, uri.Id, &payload); err != nil {
      return
    }
    res = payload
    return
  }).Put().Name("Update a {{ .ModelName }} matching :{{.PrimaryField.UriName}}")
}

// Delete route expects a pattern containing ":{{.PrimaryField.UriName}}"
func (c {{.BuilderName}}Crud) DeleteRoute(pattern string, db gsql.IExecContext) route.IRoute {
  if !strings.Contains(pattern, ":{{.PrimaryField.UriName}}") {
    panic("{{.Name}}.DeleteRoute expects a pattern containing ':{{.PrimaryField.UriName}}'")
  }
  return goof.Status(pattern, func(ctx *gin.Context) (status int, err error) {
    var uri Uri{{.Name}}Crud
    if err = ctx.ShouldBindUri(&uri); err != nil {
      return
    }
    err = c.Delete(ctx, db, uri.Id)
    return
  }).Delete().Name("Delete a {{ .ModelName }} matching :{{.PrimaryField.UriName}}")
}

// List route can use any pattern
func (c {{.BuilderName}}Crud) ListRoute(pattern string, db gsql.IQueryContext) route.IRoute {
  return goof.ToJson[[]{{ .ModelName }}](pattern, func(ctx *gin.Context) (res []{{ .ModelName }}, status int, err error) {
    var q goof.PageQuery
    if err = ctx.ShouldBindQuery(&q); err != nil {
      return
    }
    res, err = c.GetPage(ctx, db, q.Page, q.Size, q.OrderBy, q.Desc)
    return
  }).Get().Name("Return a page of []{{ .ModelName }}")
}


{{ end }}
{{ end }}