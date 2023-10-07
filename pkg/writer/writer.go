package writer

import (
	"fmt"
	"sort"

	"github.com/epsxy/flower/pkg/model"
)

func Write() {

}

// TODO: remove (deprecated)
func WriteTable(t *model.Table) string {
	var result string
	result += fmt.Sprintf("entity %s {\n", t.Name)
	for _, field := range t.Fields {
		result += fmt.Sprintf("* %s,%s\n", field.Name, string(field.Type))
	}
	result += "}\n"
	return result
}

type UMLTree struct {
	Tables []*model.Table
	Fks    []*model.ForeignKey
	Links  map[string]*model.EntityLink
}

func NewUMLBuilder() *UMLTree {
	return &UMLTree{
		Tables: []*model.Table{},
		Fks:    []*model.ForeignKey{},
	}
}

func (t *UMLTree) SetTables(tables []*model.Table) *UMLTree {
	t.Tables = tables
	return t
}

func (t *UMLTree) SetFks(fks []*model.ForeignKey) *UMLTree {
	t.Fks = fks
	return t
}

func (t *UMLTree) Build() string {
	// build
	var result string
	for _, table := range t.Tables {
		result += fmt.Sprintf("entity %s {\n", table.Name)
		currentPks := ""
		currentContent := ""
		keys := make([]string, 0, len(table.FieldsByName))
		for k := range table.FieldsByName {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			field := table.FieldsByName[key]
			if field.IsPrimaryKey {
				// a PK is always mandatory
				mandatoryDeclarator := "*"
				currentPks += fmt.Sprintf("\t%s %s, PK, %s\n", mandatoryDeclarator, field.Name, string(field.Type))
			} else {
				mandatoryDeclarator := "*"
				if field.IsNullable {
					mandatoryDeclarator = " "
				}
				currentContent += fmt.Sprintf("\t%s %s, %s\n", mandatoryDeclarator, field.Name, string(field.Type))
			}
		}
		result += currentPks
		result += "--\n"
		result += currentContent
		result += "}\n"
	}
	for _, link := range t.Links {
		result += WriteLink(link)
	}
	return result
}

/*
PLANTUML LINKS SYNTAX
----------------------
Type 	        Symbol
Zero or One 	|o--
Exactly One 	||--
Zero or Many 	}o--
One or Many 	}|--
----------------------
*/
func WriteLink(l *model.EntityLink) string {
	var leftFragment, rightFragment string
	// 0,1<-->0,1 links
	if l.Left != nil && l.Right != nil {
		if l.Left.IsNullable {
			rightFragment = "-o|"
		} else {
			rightFragment = "-||"
		}
		if l.Right.IsNullable {
			leftFragment = "|o-"
		} else {
			leftFragment = "||-"
		}
		return fmt.Sprintf("%s %s%s %s\n", l.Left.SourceName, leftFragment, rightFragment, l.Right.SourceName)
	}
	// 0,1 <-- N links
	if l.Left == nil && l.Right != nil {
		rightFragment = "-|{"
		if l.Right.IsNullable {
			leftFragment = "|o-"
		} else {
			leftFragment = "||-"
		}
		return fmt.Sprintf("%s %s%s %s\n", l.Right.DestinationName, leftFragment, rightFragment, l.Right.SourceName)
	}
	// N --> 0,1 links
	if l.Left != nil && l.Right == nil {
		leftFragment = "}|-"
		if l.Left.IsNullable {
			rightFragment = "-o|"
		} else {
			rightFragment = "-||"
		}
		return fmt.Sprintf("%s %s%s %s\n", l.Left.SourceName, leftFragment, rightFragment, l.Left.DestinationName)
	}
	return ""
}
