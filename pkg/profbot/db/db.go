package db

import (
	"log"
	"strings"
	"fmt"

	sqlite "gosqlite.googlecode.com/hg/sqlite"
)


func Open(dbname string) *sqlite.Conn {
	conn, err := sqlite.Open(dbname)
	if err != nil {
		log.Fatalf("Cannot open database: " + dbname)
	}
	return conn
}




func TableExists(conn *sqlite.Conn, tableName string) bool {
	stmt, err := conn.Prepare(
		"select 1 from sqlite_master where type = ? and tbl_name = ?")
	defer stmt.Finalize()
	if err != nil {
		log.Println("sql error: " + err.String())
		return false
	}

	stmt.Exec("table", tableName)
	return stmt.Next()
}



func CreateTable(conn *sqlite.Conn, tableName string, columns ...string) {
	if len(columns) == 0 {
		return
	}

	parameter_list := strings.Join(columns, ", ")
	query := fmt.Sprintf("create table %s (%s)", tableName, parameter_list)
	stmt, err := conn.Prepare(query)
	defer stmt.Finalize()

	if err != nil {
		log.Fatalf("Cannot create table %s; %s (%s)\n", tableName, err, query)
	}

	stmt.Exec()
	stmt.Next()
}

