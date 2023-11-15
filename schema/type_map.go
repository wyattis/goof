package schema

type typeMap map[ColumnType]string

func (t *typeMap) copy() typeMap {
	res := typeMap{}
	for k, v := range *t {
		res[k] = v
	}
	return res
}

var mysqlTypeMap = typeMap{
	TypeVarChar:   "VARCHAR",
	TypeNVarChar:  "NVARCHAR",
	TypeText:      "TEXT",
	TypeJson:      "JSON",
	TypeDateTime:  "DATETIME",
	TypeEnum:      "ENUM",
	TypeBoolean:   "BOOLEAN",
	TypeInteger:   "INTEGER",
	TypeTinyInt:   "TINYINT",
	TypeSmallInt:  "SMALLINT",
	TypeMediumInt: "MEDIUMINT",
	TypeBigInt:    "BIGINT",
	TypeDecimal:   "DECIMAL",
	TypeNumeric:   "NUMERIC",
	TypeFloat:     "FLOAT",
	TypeDate:      "DATE",
	TypeTime:      "TIME",
	TypeTimestamp: "TIMESTAMP",
	TypeBit:       "BIT",
	TypeBinary:    "BINARY",
}

var sqliteTypeMap = typeMap{
	TypeVarChar:   "TEXT",
	TypeNVarChar:  "TEXT",
	TypeText:      "TEXT",
	TypeJson:      "TEXT",
	TypeDateTime:  "TEXT",
	TypeEnum:      "TEXT",
	TypeDate:      "TEXT",
	TypeTime:      "TEXT",
	TypeTimestamp: "TEXT",
	TypeBit:       "INTEGER",
	TypeBoolean:   "INTEGER",
	TypeInteger:   "INTEGER",
	TypeTinyInt:   "INTEGER",
	TypeSmallInt:  "INTEGER",
	TypeMediumInt: "INTEGER",
	TypeBigInt:    "INTEGER",
	TypeDecimal:   "REAL",
	TypeNumeric:   "REAL",
	TypeFloat:     "REAL",
	TypeDouble:    "REAL",
	TypeBinary:    "BLOB",
	TypeVarBinary: "BLOB",
}
