package schema

type ColumnType uint

const (
	TypeVarChar ColumnType = iota
	TypeNVarChar
	TypeText
	TypeJson
	TypeInteger
	TypeBigInt
	TypeMediumInt
	TypeSmallInt
	TypeTinyInt
	TypeDecimal
	TypeFloat
	TypeNumeric
	TypeDouble
	TypeBoolean
	TypeDate
	TypeDateTime
	TypeTime
	TypeTimestamp
	TypeEnum
	TypeBit
	TypeBinary
	TypeVarBinary
	TypeBlob
)

type columnRef struct {
	Table    string
	Column   string
	Alias    string
	OnUpdate FkAction
	OnDelete FkAction
}

type columnDef struct {
	table           *Table
	OriginalName    string
	Name            string
	Kind            ColumnType
	KindLen         int
	IsUnique        bool
	IsNull          bool
	IsPrimary       bool
	IsAutoincrement bool
	Comment         string
	EnumValues      []interface{}
	ReferenceTo     *columnRef
	DefaultVal      interface{}
}

func (c *columnDef) SoloPrimary() bool {
	return c.IsPrimary && c.table.tableDef.NumPrimary() == 1
}

type columnBuilder struct {
	Column *columnDef
}

func (c *columnBuilder) applyMods(mods ...ColumnMod) *columnBuilder {
	for _, mod := range mods {
		mod(c.Column)
	}
	return c
}

func (c *columnBuilder) Autoincrement() *columnBuilder {
	return c.applyMods(Autoincrement())
}

func (c *columnBuilder) OnUpdate(action FkAction) *columnBuilder {
	if c.Column.ReferenceTo == nil {
		panic("cannot set OnUpdate on column without reference")
	}
	c.Column.ReferenceTo.OnUpdate = action
	return c
}

func (c *columnBuilder) OnDelete(action FkAction) *columnBuilder {
	if c.Column.ReferenceTo == nil {
		panic("cannot set OnUpdate on column without reference")
	}
	c.Column.ReferenceTo.OnUpdate = action
	return c
}

func (c *columnBuilder) Name(name string) *columnBuilder {
	c.Column.Name = name
	return c
}

func (c *columnBuilder) Primary() *columnBuilder {
	return c.applyMods(Primary())
}

func (c *columnBuilder) Index(name string) *columnBuilder {
	return c.applyMods(Index(name))
}

func (c *columnBuilder) Unique() *columnBuilder {
	return c.applyMods(Unique())
}

func (c *columnBuilder) Null() *columnBuilder {
	return c.applyMods(Null())
}

func (c *columnBuilder) NotNull() *columnBuilder {
	return c.applyMods(NotNull())
}

func (c *columnBuilder) References(table, column string) *columnBuilder {
	return c.applyMods(References(table, column))
}

func (c *columnBuilder) Values(values ...interface{}) *columnBuilder {
	return c.applyMods(Values(values...))
}

func (c *columnBuilder) Default(value interface{}) *columnBuilder {
	return c.applyMods(Default(value))
}

func (c *columnBuilder) Comment(comment string) *columnBuilder {
	return c.applyMods(Comment(comment))
}

func (c *columnBuilder) Type(t ColumnType) *columnBuilder {
	return c.applyMods(Type(t))
}

type ColumnMod func(*columnDef)

func VarChar(n int) ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeVarChar
		c.KindLen = n
	}
}

func NVarChar(n int) ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeNVarChar
		c.KindLen = n
	}
}

func Text() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeText
	}
}

func Json() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeJson
	}
}

func Integer() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeInteger
	}
}

func Primary() ColumnMod {
	return func(c *columnDef) {
		c.IsPrimary = true
	}
}

func Comment(comment string) ColumnMod {
	return func(c *columnDef) {
		c.Comment = comment
	}
}

func Index(name string) ColumnMod {
	return func(c *columnDef) {
		c.table.Index(c.Name).Name(name)
	}
}

func Unique() ColumnMod {
	return func(c *columnDef) {
		c.IsUnique = true
	}
}

func Null() ColumnMod {
	return func(c *columnDef) {
		c.IsNull = true
	}
}

func NotNull() ColumnMod {
	return func(c *columnDef) {
		c.IsNull = false
	}
}

func Autoincrement() ColumnMod {
	return func(c *columnDef) {
		c.IsAutoincrement = true
	}
}

func Type(t ColumnType) ColumnMod {
	return func(c *columnDef) {
		c.Kind = t
	}
}

func Binary() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeBinary
	}
}

func VarBinary() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeVarBinary
	}
}

func Blob() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeBlob
	}
}

func Datetime() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeDateTime
	}
}

func Boolean() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeBoolean
	}
}

func Date() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeDate
	}
}

func Time() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeTime
	}
}

func Timestamp() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeTimestamp
	}
}

func TinyInt() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeTinyInt
	}
}

func BigInt() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeBigInt
	}
}

func Float() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeFloat
	}
}

func Decimal() ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeDecimal
	}
}

func Enum(values ...interface{}) ColumnMod {
	return func(c *columnDef) {
		c.Kind = TypeEnum
		Values(values...)(c)
	}
}

func Values(values ...interface{}) ColumnMod {
	return func(c *columnDef) {
		c.EnumValues = values
	}
}

func References(table, col string) ColumnMod {
	return func(c *columnDef) {
		c.ReferenceTo = &columnRef{
			Table:  table,
			Column: col,
		}
	}
}

func Default(value interface{}) ColumnMod {
	return func(c *columnDef) {
		c.DefaultVal = value
	}
}
