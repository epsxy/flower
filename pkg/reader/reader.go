package reader

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/model"
	"github.com/epsxy/flower/pkg/writer"
)

func Read(data string) *writer.UMLTree {
	logger := global.GetLogger()

	// init objects
	var tables []*model.Table
	var fks []*model.ForeignKey
	var pks []*model.PrimaryKey
	tablesByName := map[string]*model.Table{}

	// prepare data
	lines := strings.Split(data, "\n")
	logger.Debug("line split", "lines", lines)

	// remove comments
	lines = global.RemoveComments(lines)
	logger.Debug("removed comments", "lines", lines)

	// remove constraints
	lines = global.RemoveConstraints(lines)
	logger.Debug("removed constraints/indexes", "lines", lines)

	// joining lines again now
	joinedLines := strings.Join(lines, " ")
	logger.Debug("join lines", "lines", joinedLines)

	// split by instruction
	instructions := strings.Split(joinedLines, ";")
	for _, instruction := range instructions {
		logger.Info("parsing instruction", "instruction", instruction)
		line := global.CleanUpString(instruction)
		logger.Info("cleaned up instruction", "instruction", line)
		// now we identify the current instruction
		if strings.Contains(line, "CREATE TABLE") {
			// case -- TABLE
			t := tableMatch(line)
			tables = append(tables, t)
			tablesByName[t.Name] = t
		} else if strings.Contains(line, "FOREIGN KEY") {
			// case -- FOREIGN KEY
			logger.Debug("ForeignKey detected")
			// parse fk declarations
			fk := fkMatchStandaloneInstruction(line)
			if fk != nil {
				fks = append(fks, fk)
			} else {
				logger.Warn("unable to parse detected FK")
			}
		} else if strings.Contains(line, "PRIMARY KEY") {
			// case -- PRIMARY KEY
			logger.Debug("ForeignKey detected")
			// parse pk declarations
			pk := pkMatchStandaloneInstruction(line)
			if pk != nil {
				pks = append(pks, pk)
			} else {
				logger.Warn("unable to parse detected PK")
			}
		}
	}
	// check all the external PK detected and declare their tables as PK/NOT NULL
	for _, pk := range pks {
		for _, field := range pk.FieldNames {
			tablesByName[pk.TableName].FieldsByName[field].IsPrimaryKey = true
			// If a field is a primary key, it has to be NOT NULL, always
			tablesByName[pk.TableName].FieldsByName[field].IsNullable = false
		}
	}
	// links are grouped in a map, with the id being `table1_table2` with table names ordered lexicographically
	links := map[string]*model.EntityLink{}
	// Create links
	for _, fk := range fks {
		// we compte the id; unique whatever the link direction ( -> or <- )
		idArr := []string{fk.SourceTable, fk.DestinationTable}
		sort.Strings(idArr)
		id := strings.Join(idArr, "_")
		if links[id] == nil {
			links[id] = &model.EntityLink{}
		}
		isNull := false
		if tablesByName[fk.SourceTable].FieldsByName[fk.SourceFields[0]] != nil {
			isNull = tablesByName[fk.SourceTable].FieldsByName[fk.SourceFields[0]].IsNullable
		} else {
			logger.Warn("field not found while building links", "table", fk.SourceTable, "field", fk.SourceFields[0])
		}
		// link is `source -> destination`
		if fk.SourceTable < fk.DestinationTable {
			links[id].Left = &model.Link{
				SourceName:      fk.SourceTable,
				DestinationName: fk.DestinationTable,
				IsNullable:      isNull,
			}
		} else {
			// link is `source <- destination`
			links[id].Right = &model.Link{
				SourceName:      fk.SourceTable,
				DestinationName: fk.DestinationTable,
				IsNullable:      isNull,
			}
		}
	}
	logger.Info("data extracted")
	return &writer.UMLTree{
		Tables: tables,
		Fks:    fks,
		Links:  links,
	}
}

func tableMatch(text string) *model.Table {
	logger := global.GetLogger()

	var table model.Table
	tableRe := regexp.MustCompile(global.TableRegexp)
	match := tableRe.FindStringSubmatch(text)
	if match != nil {
		groupNames := tableRe.SubexpNames()
		captures := make(map[string]string)

		for i, name := range groupNames {
			if i != 0 && name != "" {
				captures[name] = match[i]
			}
		}
		tableName := captures["tableName"]
		tableFields := captures["tableFields"]
		fields, pk := fieldsMatch(tableFields)
		if pk == nil {
			logger.Warn("a primary key was not detected in the table declaraction", "table", tableName)
		}
		fByName := map[string]*model.Field{}
		for _, field := range fields {
			fByName[field.Name] = field
			if pk != nil && global.Contains(pk.FieldNames, field.Name) {
				fByName[field.Name].IsPrimaryKey = true
			}
		}
		table = model.Table{
			Name:         tableName,
			Fields:       fields,
			FieldsByName: fByName,
			Pk:           pk,
		}
	}
	return &table
}

func fieldsMatch(text string) ([]*model.Field, *model.PrimaryKey) {
	logger := global.GetLogger()
	logger.Warn("reading fields", "text", text)
	var fields []*model.Field
	// split fields
	fieldsSplit := strings.Split(text, ",")
	logger.Warn("split fields", "fields", fieldsSplit)

	// looping over all fields
	for _, f := range fieldsSplit {
		logger.Debug("processing field", "field", f)
		currentField := &model.Field{
			Name:         "",
			Type:         "",
			IsPrimaryKey: strings.Contains(f, "PRIMARY KEY"),
			IsNullable:   !strings.Contains(f, "NOT NULL"),
		}
		purgedString := strings.ReplaceAll(f, "NOT NULL", "")
		purgedString = strings.ReplaceAll(purgedString, "PRIMARY KEY", "")
		pkRe := regexp.MustCompile(global.FieldRegexp)
		match := pkRe.FindStringSubmatch(purgedString)
		if len(match) >= 3 {
			fieldName := match[1]
			dataType := match[2]
			length := match[3]

			fieldType := dataType
			if length != "" {
				fieldType = fmt.Sprintf("%s[%s]", dataType, length)
			}
			currentField.Name = fieldName
			currentField.Type = fieldType
		} else {
			logger.Warn("field didn't matched regex", "text", text)
		}
		if currentField.Name != "" {
			logger.Warn("parsed field, but with an empty name", "field", f)
			fields = append(fields, currentField)
		}
	}
	return fields, pkMatchInTableDeclaration(text)
}

func pkMatchInTableDeclaration(text string) *model.PrimaryKey {
	logger := global.GetLogger()

	pkRe := regexp.MustCompile(global.InTablePkRegexp)
	match := pkRe.FindStringSubmatch(text)
	if len(match) > 0 {
		fields := global.SplitAndGetFields(match[1])
		return &model.PrimaryKey{FieldNames: fields}
	} else {
		logger.Info("table is not containing a primary key declaration", "text", text)
	}
	return nil
}

func pkMatchStandaloneInstruction(text string) *model.PrimaryKey {
	logger := global.GetLogger()
	logger.Debug("trying to get pk from string", "text", text)

	pkRe := regexp.MustCompile(global.StandalonePkRegex)
	matches := pkRe.FindStringSubmatch(text)
	if len(matches) > 2 {
		return &model.PrimaryKey{TableName: matches[1], FieldNames: global.SplitAndGetFields(matches[2])}
	} else {
		logger.Warn("pk didn't matched regex", "text", text)
	}
	return nil
}

func fkMatchStandaloneInstruction(text string) *model.ForeignKey {
	logger := global.GetLogger()

	re := regexp.MustCompile(global.StandaloneFkRegexp)
	matches := re.FindStringSubmatch(text)

	if len(matches) > 0 {
		tableSource := matches[1]
		fieldSource := matches[2]
		tableTarget := matches[3]
		fieldTarget := matches[4]

		return &model.ForeignKey{
			SourceTable:       tableSource,
			DestinationTable:  tableTarget,
			SourceFields:      global.SplitAndGetFields(fieldSource),
			DestinationFields: global.SplitAndGetFields(fieldTarget),
		}
	} else {
		logger.Warn("fk didn't matched", "text", text)
		return nil
	}
}
