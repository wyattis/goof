package schema

type indexDef struct {
	Table       *TableDef
	Name        string
	Unique      bool
	IfNotExists bool
	Columns     []string
}

type indexBuilder struct {
	index *indexDef
}

func (t *indexBuilder) Unique() *indexBuilder {
	t.index.Unique = true
	return t
}

func (t *indexBuilder) Name(name string) *indexBuilder {
	t.index.Name = name
	return t
}

func (t *indexBuilder) IfNotExists() *indexBuilder {
	t.index.IfNotExists = true
	return t
}
