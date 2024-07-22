package helpers

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var sqltracing = false

// History table
const tableHistory = "history"

// History table columns
const colDateStamp = "date_stamp"
const colTimeStamp = "time_stamp"
const colState = "state"
const colPctDem = "pct_dem"
const colPctGop = "pct_gop"
const colPctZero = "pct_zero"
const colStartDate = "start_date"
const colEndDate = "end_date"
const colPollster = "pollster"

// Record insertion interface struct
const ixEndDate = "ix_end_date"

// Database parameters
type dbparams struct {
	state     string
	startDate string
	endDate   string
	pctDem    float64
	pctGop    float64
	pollster  string
}

// Assigned and used at run-time
var pathDatabase string
var sqliteDatabase *sql.DB

/*
Internal function to run an SQL statement and handle any errors.
*/
func sqlFunc(text string) {

	if sqltracing {
		log.Printf("sqlFunc: %s\n", text)
	}

	statement, err := sqliteDatabase.Prepare(text) // Prepare SQL Statement
	if err != nil {
		log.Fatalf("sqlFunc: sqliteDatabase.Prepare failed\n%s\nreason: %s\n",
			text,
			err.Error())
	}

	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatalf("sqlFunc: statement.Exec failed\n%s\nreason:%s\n",
			text,
			err.Error())
	}

}

/*
Internal function to run an SQL select query and handle any errors. The output is returned to caller.
*/
func sqlQuery(text string) *sql.Rows {

	if sqltracing {
		log.Printf("sqlQuery: %s\n", text)
	}

	rows, err := sqliteDatabase.Query(text)
	if err != nil {
		log.Fatalf("sqlQuery: sqliteDatabase.Query failed\n%s\nreason: %s\n",
			text,
			err.Error())
	}

	return rows

}

/*
Internal function to initialise a jacotest database.

* Create database (includes file creation/re-creation).
* Create history table and all of its columns, a combination of which is the primary index.
* Create secondary indexes.
*/
func initDB() {

	if sqltracing {
		log.Println("initDB: Begin")
	}

	sqlText := "CREATE TABLE " + tableHistory + " ("
	sqlText += colDateStamp + " VARCHAR NOT NULL, "
	sqlText += colTimeStamp + " VARCHAR NOT NULL, "
	sqlText += colState + " VARCHAR NOT NULL, "
	sqlText += colStartDate + " VARCHAR NOT NULL, "
	sqlText += colEndDate + " VARCHAR NOT NULL, "
	sqlText += colPctDem + " FLOAT NOT NULL, "
	sqlText += colPctGop + " FLOAT NOT NULL, "
	sqlText += colPollster + " VARCHAR NOT NULL, "
	sqlText += "PRIMARY KEY (" + colState + ", " + colEndDate + ") )"
	sqlFunc(sqlText)

	sqlText = "CREATE INDEX " + ixEndDate + " ON " + tableHistory + " (" + colEndDate + ")"
	sqlFunc(sqlText)

	if sqltracing {
		log.Println("initDB: End")
	}

}

/*
DBOpen - Database Open

* If the database directory is not present, create it.
* If the history.db file in the database directory is not present, call initDB.
* Connect to DB.
* Validate DB.
*/
func DBOpen(driverDatabase, dirDatabase, fileDatabase string) {

	if sqltracing {
		log.Printf("DBOpen: Begin")
	}

	// Database file
	pathDatabase = dirDatabase + "/" + fileDatabase
	_, err := os.Stat(pathDatabase)
	if err != nil {
		if sqltracing {
			log.Printf("DBOpen: database file(%s) inaccessible, will create it.",
				pathDatabase)
		}
		sqliteDatabase, err = sql.Open(driverDatabase, pathDatabase)
		if err != nil {
			log.Fatalf("DBOpen: sql.Open/create(%s) failed, reason: %s",
				pathDatabase,
				err.Error())
		}
		initDB()

		if sqltracing {
			log.Printf("DBOpen: End, database created")
		}
		return
	}

	// Connect to pre-existing database
	if sqltracing {
		log.Printf("DBOpen database file exists")
	}
	sqliteDatabase, err = sql.Open(driverDatabase, pathDatabase)
	if err != nil {
		log.Fatalf("DBOpen: sql.Open/pre-existing(%s) failed, reason: %s",
			pathDatabase,
			err.Error())
	}

	// sqliteDatabase stays open until process exit

	if sqltracing {
		log.Printf("DBOpen: End, existing database opened")
	}

}

/*
DBClose - Close the database.
*/
func DBClose() {

	if sqltracing {
		log.Printf("DBClose: Begin")
	}

	err := sqliteDatabase.Close()
	if err != nil {
		log.Fatalf("DBOpen: sql.Close(%s) failed, reason: %s",
			pathDatabase,
			err.Error())
	}

	if sqltracing {
		log.Printf("DBClose: End")
	}

}

/*
DBStore - Store a poll record.
*/
func DBStore(fields dbparams) {

	dateUTC := "'" + GetUtcDate() + "'"
	timeUTC := "'" + GetUtcTime() + "'"
	sqlText := "INSERT OR REPLACE INTO " + tableHistory + " VALUES("
	sqlText += dateUTC + ", " + timeUTC + ",\"" + fields.state + "\", \"" + fields.startDate + "\", \"" + fields.endDate
	caboose := fmt.Sprintf("\", %f, %f, \"%s\" )", fields.pctDem, fields.pctGop, fields.pollster)
	sqlText += caboose

	sqlFunc(sqlText)
}
