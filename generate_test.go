package boilerplate

import (
	"testing"

	"github.com/google/uuid"
)

type T1 struct {
	ID     int64     `json:"id"     db:"id"     dbtype:"BIGSERIAL NOT NULL PRIMARY KEY"`
	Field1 string    `json:"field1" db:"field1" dbtype:"TEXT NOT NULL"`
	Field2 int64     `json:"field2" db:"field2" dbtype:"BIGSERIAL NOT NULL"`
	Field3 bool      `json:"field3" db:"field3" dbtype:"BOOLEAN NOT NULL"`
	Field4 uuid.UUID `json:"field4" db:"field4" dbtype:"UUID"`
}

func TestGenerateQueries(t *testing.T) {
	t.Run("sqlite", func(t *testing.T) {
		queries := GenerateQueries(GenerateQueriesOptions{
			TableName:          "t1",
			Model:              T1{},
			AutoGeneratingCols: []string{"id", "field2"},
			PrimaryKeys:        []string{"id"},
			Driver:             DriverSqlite,
		})
		expected := struct {
			CreateTable string
			DropTable   string
			Select      string
			Insert      string
			Update      string
			Upsert      string
			Delete      string
		}{
			CreateTable: `CREATE TABLE IF NOT EXISTS t1 (id INTEGER PRIMARY KEY AUTOINCREMENT, field1 TEXT NOT NULL, field2 INTEGER AUTOINCREMENT, field3 BOOLEAN NOT NULL, field4 TEXT)`,
			DropTable:   `DROP TABLE IF EXISTS t1`,
			Select:      `SELECT * FROM t1`,
			Insert:      `INSERT INTO t1 (field1, field3, field4) VALUES (:field1, :field3, :field4) RETURNING *`,
			Update:      `UPDATE t1 SET field1 = :field1, field3 = :field3, field4 = :field4 WHERE id = :id RETURNING *`,
			Upsert:      `INSERT INTO t1 (field1, field3, field4) VALUES (:field1, :field3, :field4) ON CONFLICT (id) DO UPDATE SET field1 = EXCLUDED.field1, field3 = EXCLUDED.field3, field4 = EXCLUDED.field4 RETURNING *`,
			Delete:      `DELETE FROM t1 WHERE id = :id`,
		}

		if queries.CreateTable != expected.CreateTable {
			t.Errorf("CreateTable is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.CreateTable, queries.CreateTable)
		}
		if queries.DropTable != expected.DropTable {
			t.Errorf("DropTable is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.DropTable, queries.DropTable)
		}
		if queries.Select != expected.Select {
			t.Errorf("Select is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Select, queries.Select)
		}
		if queries.Insert != expected.Insert {
			t.Errorf("Insert is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Insert, queries.Insert)
		}
		if queries.Update != expected.Update {
			t.Errorf("Update is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Update, queries.Update)
		}
		if queries.Upsert != expected.Upsert {
			t.Errorf("Upsert is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Upsert, queries.Upsert)
		}
		if queries.Delete != expected.Delete {
			t.Errorf("Delete is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Delete, queries.Delete)
		}
	})
	t.Run("postgres", func(t *testing.T) {
		queries := GenerateQueries(GenerateQueriesOptions{
			TableName:          "t1",
			Model:              T1{},
			AutoGeneratingCols: []string{"id", "field2"},
			PrimaryKeys:        []string{"id"},
			Driver:             DriverPostgres,
		})
		expected := struct {
			CreateTable string
			DropTable   string
			Select      string
			Insert      string
			Update      string
			Upsert      string
			Delete      string
		}{
			CreateTable: `CREATE TABLE IF NOT EXISTS t1 (id BIGSERIAL NOT NULL PRIMARY KEY, field1 TEXT NOT NULL, field2 BIGSERIAL NOT NULL, field3 BOOLEAN NOT NULL, field4 UUID)`,
			DropTable:   `DROP TABLE IF EXISTS t1`,
			Select:      `SELECT * FROM t1`,
			Insert:      `INSERT INTO t1 (field1, field3, field4) VALUES (:field1, :field3, :field4) RETURNING *`,
			Update:      `UPDATE t1 SET field1 = :field1, field3 = :field3, field4 = :field4 WHERE id = :id RETURNING *`,
			Upsert:      `INSERT INTO t1 (field1, field3, field4) VALUES (:field1, :field3, :field4) ON CONFLICT (id) DO UPDATE SET field1 = EXCLUDED.field1, field3 = EXCLUDED.field3, field4 = EXCLUDED.field4 RETURNING *`,
			Delete:      `DELETE FROM t1 WHERE id = :id`,
		}

		if queries.CreateTable != expected.CreateTable {
			t.Errorf("CreateTable is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.CreateTable, queries.CreateTable)
		}
		if queries.DropTable != expected.DropTable {
			t.Errorf("DropTable is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.DropTable, queries.DropTable)
		}
		if queries.Select != expected.Select {
			t.Errorf("Select is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Select, queries.Select)
		}
		if queries.Insert != expected.Insert {
			t.Errorf("Insert is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Insert, queries.Insert)
		}
		if queries.Update != expected.Update {
			t.Errorf("Update is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Update, queries.Update)
		}
		if queries.Upsert != expected.Upsert {
			t.Errorf("Upsert is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Upsert, queries.Upsert)
		}
		if queries.Delete != expected.Delete {
			t.Errorf("Delete is not what is expected:\n\nExpected:\n%s\n\nActual:\n%s\n", expected.Delete, queries.Delete)
		}
	})
}
