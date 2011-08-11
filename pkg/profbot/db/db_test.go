package db


import (
	"testing"
)


func TestTableExists(t *testing.T) {
	conn := Open(":memory:")
	defer conn.Close()

	if TableExists(conn, "test_table") {
		t.Fatalf("Table 'test_table' should not exist yet")
	}


	stmt, err := conn.Prepare("create table test_table (id)")
	if err != nil {
		t.Fatalf("Error in SQL create statement")
	}
	stmt.Exec()
	stmt.Next()

	if !TableExists(conn, "test_table") {
		t.Fatalf("Table 'test_table' should exist")
	}

}



func TestTableCreate(t *testing.T) {
	conn := Open(":memory:")
	defer conn.Close()

	if TableExists(conn, "test_table") {
		t.Fatalf("Table 'test_table' should not exist yet")
	}

	CreateTable(conn, "test_table")
	if TableExists(conn, "test_table") {
		t.Fatalf("Table 'test_table' should not exist after empty column list")
	}

	CreateTable(conn, "test_table", "id integer", "name string")
	if !TableExists(conn, "test_table") {
		t.Fatalf("Table 'test_table' should have been created")
	}
}
