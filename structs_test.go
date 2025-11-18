package boilerplate

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB

func LoadDB(t *testing.T) {
	if db != nil {
		return
	}
	var err error
	db, err = Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
}

func AssertStructEqual(t *testing.T, expected any, got any, note string) {
	expectedJSON, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		panic(err)
	}
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		panic(err)
	}

	if string(expectedJSON) != string(gotJSON) {
		t.Fatalf("%s:\n\nExpected:\n%s\n\nGot:\n%s\n\n", note, expectedJSON, gotJSON)
	}
}

func TestInit(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	LoadDB(t)
}

func TestStringSlice(t *testing.T) {
	type Table struct {
		ID          int         `db:"id" dbtype:"BIGSERIAL NOT NULL PRIMARY KEY"`
		StringSlice StringSlice `db:"string_slice" dbtype:"TEXT[] NOT NULL"`
	}
	queries := GenerateQueries(GenerateQueriesOptions{
		TableName:          "table_string_slice",
		Model:              Table{},
		AutoGeneratingCols: []string{"id"},
		PrimaryKeys:        []string{"id"},
		Driver:             DriverSqlite,
	})

	t.Cleanup(func() {
		_ = Exec(db, queries.DropTable)
	})

	err := Exec(db, queries.CreateTable)
	if err != nil {
		t.Fatal(err)
	}
	row := Table{
		StringSlice: StringSlice{"test1", "test2", "test3"},
	}
	err = NamedExecReturning(db, &row, queries.Insert, &row)
	if err != nil {
		t.Fatal(err)
	}
	if row.ID == 0 {
		t.Fatal("Expected inserted row to have an id, but it's still 0.")
	}
	AssertStructEqual(t, StringSlice{"test1", "test2", "test3"}, row.StringSlice, "Expected inserted StringSlice to match what was inserted")
	out, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), row.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, row.StringSlice, out.StringSlice, "Expected selected StringSlice to match inserted StringSlice")

	out.StringSlice = StringSlice{"test2", "test3", "test4"}
	err = NamedExecReturning(db, &out, queries.Update, &out)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, StringSlice{"test2", "test3", "test4"}, out.StringSlice, "Expected updated StringSlice to match what was updated")

	out2, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), row.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, out.StringSlice, out2.StringSlice, "Expected selected StringSlice to match updated StringSlice")

	insertNilStringSlice := Table{
		StringSlice: nil,
	}
	err = NamedExecReturning(db, &insertNilStringSlice, queries.Insert, &insertNilStringSlice)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, nil, insertNilStringSlice.StringSlice, "Expected inserted nil StringSlice to be nil")
	err = NamedExecReturning(db, &insertNilStringSlice, queries.Insert, &insertNilStringSlice)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, nil, insertNilStringSlice.StringSlice, "Expected inserted Nil StringSlice to be nil on select")

	selectNilStringSlice, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), insertNilStringSlice.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, selectNilStringSlice.StringSlice, nil, "Expected selected nil StringSlice to be nil")
}

func TestIntSlice(t *testing.T) {
	type Table struct {
		ID       int      `db:"id" dbtype:"BIGSERIAL NOT NULL PRIMARY KEY"`
		IntSlice IntSlice `db:"int_slice" dbtype:"INT[] NOT NULL"`
	}
	queries := GenerateQueries(GenerateQueriesOptions{
		TableName:          "table_int_slice",
		Model:              Table{},
		AutoGeneratingCols: []string{"id"},
		PrimaryKeys:        []string{"id"},
		Driver:             DriverSqlite,
	})

	t.Cleanup(func() {
		_ = Exec(db, queries.DropTable)
	})

	err := Exec(db, queries.CreateTable)
	if err != nil {
		t.Fatal(err)
	}
	row := Table{
		IntSlice: IntSlice{1, 2, 3},
	}
	err = NamedExecReturning(db, &row, queries.Insert, &row)
	if err != nil {
		t.Fatal(err)
	}
	if row.ID == 0 {
		t.Fatal("Expected inserted row to have an id, but it's still 0.")
	}
	AssertStructEqual(t, IntSlice{1, 2, 3}, row.IntSlice, "Expected inserted IntSlice to match what was inserted")
	out, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), row.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, row.IntSlice, out.IntSlice, "Expected selected IntSlice to match inserted IntSlice")

	out.IntSlice = IntSlice{2, 3, 4}
	err = NamedExecReturning(db, &out, queries.Update, &out)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, IntSlice{2, 3, 4}, out.IntSlice, "Expected updated IntSlice to match what was updated")

	out2, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), row.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, out.IntSlice, out2.IntSlice, "Expected selected IntSlice to match updated IntSlice")

	insertNilIntSlice := Table{
		IntSlice: nil,
	}
	err = NamedExecReturning(db, &insertNilIntSlice, queries.Insert, &insertNilIntSlice)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, nil, insertNilIntSlice.IntSlice, "Expected inserted nil IntSlice to be nil")
	err = NamedExecReturning(db, &insertNilIntSlice, queries.Insert, &insertNilIntSlice)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, nil, insertNilIntSlice.IntSlice, "Expected inserted Nil IntSlice to be nil on select")

	selectNilIntSlice, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), insertNilIntSlice.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, selectNilIntSlice.IntSlice, nil, "Expected selected nil IntSlice to be nil")
}

func TestJsonObject(t *testing.T) {
	type Table struct {
		ID         int        `db:"id" dbtype:"BIGSERIAL NOT NULL PRIMARY KEY"`
		JsonObject JsonObject `db:"json_object" dbtype:"JSONB NOT NULL"`
	}
	queries := GenerateQueries(GenerateQueriesOptions{
		TableName:          "table_json_object",
		Model:              Table{},
		AutoGeneratingCols: []string{"id"},
		PrimaryKeys:        []string{"id"},
		Driver:             DriverSqlite,
	})

	t.Cleanup(func() {
		_ = Exec(db, queries.DropTable)
	})

	err := Exec(db, queries.CreateTable)
	if err != nil {
		t.Fatal(err)
	}
	row := Table{
		JsonObject: JsonObject{
			"test1":  123,
			"array1": []any{1, "2", nil},
			"test2":  nil,
			"test3": JsonObject{
				"test1": "abc",
			},
		},
	}
	err = NamedExecReturning(db, &row, queries.Insert, &row)
	if err != nil {
		t.Fatal(err)
	}
	if row.ID == 0 {
		t.Fatal("Expected inserted row to have an id, but it's still 0.")
	}
	AssertStructEqual(t,
		JsonObject{
			"test1":  123,
			"array1": []any{1, "2", nil},
			"test2":  nil,
			"test3": JsonObject{
				"test1": "abc",
			},
		},
		row.JsonObject, "Expected inserted JsonObject to match what was inserted")
	out, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), row.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, row.JsonObject, out.JsonObject, "Expected selected JsonObject to match inserted JsonObject")

	out.JsonObject = JsonObject{}
	err = NamedExecReturning(db, &out, queries.Update, &out)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, JsonObject{}, out.JsonObject, "Expected updated JsonObject to match what was updated")

	out2, err := Get[Table](db, fmt.Sprintf("%s %s", queries.Select, "WHERE id = $1"), row.ID)
	if err != nil {
		t.Fatal(err)
	}
	AssertStructEqual(t, out.JsonObject, out2.JsonObject, "Expected selected JsonObject to match updated JsonObject")
}
