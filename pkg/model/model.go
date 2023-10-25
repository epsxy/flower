package model

import (
	"sort"
	"strings"
)

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

func (l *EntityLink) Id() string {
	if l.Left != nil {
		return GenLinkId(l.Left.DestinationName, l.Left.SourceName)
	}
	if l.Right != nil {
		return GenLinkId(l.Right.DestinationName, l.Right.SourceName)
	}
	return ""
}

func GenLinkId(tableA, tableB string) string {
	idArr := []string{tableA, tableB}
	sort.Strings(idArr)
	return strings.Join(idArr, "_")
}
