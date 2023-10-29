package writer

import (
	"fmt"
	"sort"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/graph"
	"github.com/epsxy/flower/pkg/model"
)

func Build(t *model.UMLTree) []string {
	//logger := global.GetLogger()
	// list vertexes
	var vertexes []string
	for _, table := range t.Tables {
		vertexes = append(vertexes, table.Name)
	}
	if !t.Options.SplitUnconnected && !t.Options.SplitDistance {
		// in that case, we're not going to split anything, so we can build the ERD
		return []string{_buildImpl(vertexes, t)}
	}
	// build graph structure
	g := graph.Gen(vertexes, t.Links)
	// split graph in connected partitions if enabled
	var partitions [][]string
	if t.Options.SplitUnconnected {
		partitions = graph.Dfs(vertexes, g)
	} else {
		partitions = [][]string{vertexes}
	}
	// if distance split enabled, refine each connected partitions
	var newPartitions [][]string
	if t.Options.SplitDistance {
		// TODO
		for _, p := range partitions {
			// TODO process each partition
			res := graph.Split(p, g, *t.Options)
			newPartitions = append(newPartitions, res...)
		}
	} else {
		// nothing else to do
		newPartitions = partitions
	}
	var response []string
	for _, p := range newPartitions {
		response = append(response, _buildImpl(p, t))
	}
	return response
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

func _generateDocumentName(vertexes []string) string {
	if len(vertexes) == 1 {
		return vertexes[0]
	}
	return "tables"
}

func _buildImpl(vertexes []string, t *model.UMLTree) string {
	logger := global.GetLogger()
	start := fmt.Sprintf("@startuml %s\n", _generateDocumentName(vertexes))
	end := "@enduml"
	tables := ""
	links := ""
	visitedLink := map[string]bool{}
	for _, v := range vertexes {
		tables += _buildCurrentTable(t.TablesByName[v])
		for _, link := range t.LinksByTableName[v] {
			if !visitedLink[link.Id()] {
				visitedLink[link.Id()] = true
				links += WriteLink(link)
			}
		}
	}
	logger.Info("basic build finished")
	return start + tables + links + end
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
