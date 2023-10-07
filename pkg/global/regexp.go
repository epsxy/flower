package global

const (
	TableRegexp             = `CREATE TABLE(?: IF NOT EXISTS)? (?P<tableName>\w+.\w+) \((?P<tableFields>.+)\)`
	FieldRegexp             = `^\s*(\w+)\s+([a-zA-Z0-9_ ]*)(?:\((\d+)\))?`
	InTablePkRegexp         = `PRIMARY\s+KEY\s+\((.*?)\)`
	StandaloneFkRegexp      = `ALTER\s+TABLE\s+(?:ONLY\s+)?(\w+.?\w+)\s+ADD\s+(?:CONSTRAINT\s+)?\w+\s+FOREIGN\s+KEY\s+\(([\w\s,]+)\)\s+REFERENCES\s+(\w+.?\w+)\s*\(([\w\s,]+)\)`
	StandalonePkRegex       = `ALTER\s+TABLE\s+(?:ONLY\s+)?(\w+.?\w+)\s+ADD\s+CONSTRAINT\s+\w+\s+PRIMARY\s+KEY\s+\(([\w\s,]+)\)`
	ConstraintRemovalRegexp = `^\s*CREATE INDEX.*`
	IndexRemovalRegexp      = `^\s*CONSTRAINT.*`
)
