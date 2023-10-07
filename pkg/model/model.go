package model

type Table struct {
	Name         string            `yaml:"name,flow"`
	Fields       []*Field          `yaml:"fields,flow"`
	Fks          []*ForeignKey     `yaml:"fks,flow"`
	Pk           *PrimaryKey       `yaml:"pk,flow"`
	FieldsByName map[string]*Field `yaml:"fieldsByName,flow"`
}

type Field struct {
	Name         string `yaml:"name,flow"`
	Type         string `yaml:"type,flow"`
	IsPrimaryKey bool   `yaml:"isPrimaryKey,flow"`
	IsNullable   bool   `yaml:"isNotNull,flow"`
}

type ForeignKey struct {
	SourceTable       string   `yaml:"sourceTable,flow"`
	SourceFields      []string `yaml:"sourceFields,flow"`
	DestinationTable  string   `yaml:"destinationTable,flow"`
	DestinationFields []string `yaml:"destinationFields,flow"`
}

type PrimaryKey struct {
	TableName  string   `yaml:"tableName,flow"`
	FieldNames []string `yaml:"fieldNames,flow"`
}

// EntityLink is describing how tables are linked together.
// - Left will always contain the first table in lexicographic order
// - Right will always contain the first table in lexicographic order
// eg: with tables `authors` (left) and `posts` (right)
// we are going to build the link like this: `authors <-> posts`, with
// - Left: authors (source) -> posts (destination)
// - Right: authors (destination) <- posts (source)
type EntityLink struct {
	Left  *Link `yaml:"left,flow"`
	Right *Link `yaml:"right,flow"`
}

type Link struct {
	SourceName      string `yaml:"sourceName,flow"`
	DestinationName string `yaml:"destinationName,flow"`
	IsNullable      bool   `yaml:"isNullable,flow"`
}
