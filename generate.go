package boilerplate

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Driver string

const (
	DriverSqlite   = "sqlite"
	DriverPostgres = "postgres"
)

type GenerateQueriesOptions struct {
	TableName          string
	Model              any
	AutoGeneratingCols []string
	PrimaryKeys        []string
	Driver             Driver
}

type GeneratedQueries struct {
	CreateTable string
	DropTable   string
	Select      string
	Insert      string
	Update      string
	Upsert      string
	Delete      string
}

// Uses the struct fields 'db' and 'dbtype' to auto generate most queries.
//
// `AutoGeneratingCols` is a list of col names that are auto generated on INSERT or UPSERT, like 'BIGSERIAL', etc. These will be excluded from INSERT statements to let the database auto generate.
//
// `PrimaryKeys` is a list of col names that are used to control what is done on an UPDATE query.
func GenerateQueries(opts GenerateQueriesOptions) (queries GeneratedQueries) {
	type col struct {
		name string
		sql  string
	}

	if len(opts.PrimaryKeys) == 0 {
		panic("No primary key specified for table: " + opts.TableName)
	}

	modelType := reflect.TypeOf(opts.Model)

	cols := []col{}

	for i := range modelType.NumField() {
		field := modelType.Field(i)
		name := field.Tag.Get("db")
		dbtype := strings.ToUpper(field.Tag.Get("dbtype"))
		switch opts.Driver {
		case DriverSqlite:
			// sqlite does not support SERIAL, so we switch it to AUTOINCREMENT
			if strings.Contains(dbtype, "SERIAL") {
				dbtype = strings.ReplaceAll(dbtype, "BIGSERIAL", "INTEGER")
				dbtype = strings.ReplaceAll(dbtype, "SERIAL", "INTEGER")
				dbtype = strings.ReplaceAll(dbtype, "NOT NULL", "")
				dbtype = dbtype + " AUTOINCREMENT"
				dbtype = strings.ReplaceAll(dbtype, "  ", " ")
			}
			// sqlite does not support UUID, so we switch it to TEXT
			if strings.Contains(dbtype, "UUID") {
				dbtype = strings.ReplaceAll(dbtype, "UUID", "TEXT")
			}
		case DriverPostgres:
			// postgres does not support AUTOINCREMENT, so we switch it to SERIAL
			if strings.Contains(dbtype, "AUTOINCREMENT") {
				dbtype = strings.ReplaceAll(dbtype, "AUTOINCREMENT", "")
				dbtype = strings.ReplaceAll(dbtype, "INT", "SERIAL")
			}
		}

		if name == "" || dbtype == "" || name == "-" || dbtype == "-" {
			continue
		}

		cols = append(cols, col{
			name: name,
			sql:  dbtype,
		})
	}

	// CREATE TABLE
	colStrings := []string{}
	for _, c := range cols {
		colStrings = append(colStrings, fmt.Sprintf("%s %s", c.name, c.sql))
	}
	queries.CreateTable = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", opts.TableName, strings.Join(colStrings, ", "))

	// DROP TABLE
	queries.DropTable = fmt.Sprintf("DROP TABLE IF EXISTS %s", opts.TableName)

	// SELECT
	queries.Select = fmt.Sprintf("SELECT * FROM %s", opts.TableName)

	// INSERT
	inserts := []string{}
	vals := []string{}
	for _, c := range cols {
		if slices.Contains(opts.AutoGeneratingCols, c.name) {
			continue
		}
		inserts = append(inserts, c.name)
		vals = append(vals, ":"+c.name)
	}
	queries.Insert = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *", opts.TableName, strings.Join(inserts, ", "), strings.Join(vals, ", "))

	// UPDATE
	sets := []string{}
	wheres := []string{}
	for _, c := range cols {
		if slices.Contains(opts.PrimaryKeys, c.name) {
			wheres = append(wheres, fmt.Sprintf("%s = :%s", c.name, c.name))
			continue
		}
		if slices.Contains(opts.AutoGeneratingCols, c.name) {
			continue
		}
		sets = append(sets, fmt.Sprintf("%s = :%s", c.name, c.name))
	}
	queries.Update = fmt.Sprintf("UPDATE %s SET %s WHERE %s RETURNING *", opts.TableName, strings.Join(sets, ", "), strings.Join(wheres, " AND "))

	// UPSERT
	upserts := []string{}
	for _, c := range cols {
		if slices.Contains(opts.AutoGeneratingCols, c.name) || slices.Contains(opts.PrimaryKeys, c.name) {
			continue
		}

		upserts = append(upserts, fmt.Sprintf("%s = EXCLUDED.%s", c.name, c.name))
	}
	queries.Upsert = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s RETURNING *", opts.TableName, strings.Join(inserts, ", "), strings.Join(vals, ", "), opts.PrimaryKeys[0], strings.Join(upserts, ", "))

	// DELETE
	queries.Delete = fmt.Sprintf("DELETE FROM %s WHERE %s", opts.TableName, strings.Join(wheres, " AND "))

	return
}
