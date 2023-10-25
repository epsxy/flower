package writer

import (
	"fmt"
	"sort"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/graph"
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
	startUml := "@startuml todo\n"
	endUml := "@enduml"
	return startUml + result + endUml
}

type UMLTree struct {
	LinksByTableName map[string][]*model.EntityLink
	TablesByName     map[string]*model.Table
	Tables           []*model.Table
	Fks              []*model.ForeignKey
	Links            map[string]*model.EntityLink
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

func _buildCurrentTable(table *model.Table) string {
	result := ""
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
	return result
}

func (t *UMLTree) Build() string {
	logger := global.GetLogger()

	// build
	var result string
	for _, table := range t.Tables {
		result += _buildCurrentTable(table)
	}
	for _, link := range t.Links {
		result += WriteLink(link)
	}
	logger.Info("basic build finished")
	return result
}

func (t *UMLTree) BuildWithPartitions() []string {
	logger := global.GetLogger()
	res := []string{}

	var tableNames []string
	var g map[string][]string = map[string][]string{}
	var visited map[string]bool = map[string]bool{}
	for _, table := range t.Tables {
		tableNames = append(tableNames, table.Name)
		g[table.Name] = []string{}
		visited[table.Name] = false
	}
	for _, link := range t.Links {
		if link.Left != nil {
			if !global.Contains(g[link.Left.SourceName], link.Left.DestinationName) {
				g[link.Left.SourceName] = append(g[link.Left.SourceName], link.Left.DestinationName)
			}
			if !global.Contains(g[link.Left.DestinationName], link.Left.SourceName) {
				g[link.Left.DestinationName] = append(g[link.Left.DestinationName], link.Left.SourceName)
			}
		}
		if link.Right != nil {
			if !global.Contains(g[link.Right.SourceName], link.Right.DestinationName) {
				g[link.Right.SourceName] = append(g[link.Right.SourceName], link.Right.DestinationName)
			}
			if !global.Contains(g[link.Right.DestinationName], link.Right.SourceName) {
				g[link.Right.DestinationName] = append(g[link.Right.DestinationName], link.Right.SourceName)
			}
		}
	}
	// split graph into partitions
	partitions := graph.Dfs_root(tableNames, g, visited)
	// export each partition to file
	for _, p := range partitions {
		currentRes := ""
		currentResLinks := ""
		linksBuiltMap := map[string]bool{}
		for _, tableName := range p {
			var res string
			currentRes += _buildCurrentTable(t.TablesByName[tableName])
			res, linksBuiltMap = _buildLinks(t.LinksByTableName[tableName], linksBuiltMap)
			currentResLinks += res
		}
		startUml := "@startuml todo\n"
		endUml := "@enduml"
		res = append(res, startUml+currentRes+currentResLinks+endUml)
	}
	logger.Warn("finished build with partitions")
	return res
}

func _buildLinks(links []*model.EntityLink, linksAlreadyBuilt map[string]bool) (string, map[string]bool) {
	res := ""
	for _, link := range links {
		if link.Left != nil {
			if !linksAlreadyBuilt[link.Id()] {
				res += WriteLink(link)
				linksAlreadyBuilt[link.Id()] = true
			}
		} else if link.Right != nil {
			if !linksAlreadyBuilt[link.Id()] {
				res += WriteLink(link)
				linksAlreadyBuilt[link.Id()] = true
			}
		}
	}
	return res, linksAlreadyBuilt
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
