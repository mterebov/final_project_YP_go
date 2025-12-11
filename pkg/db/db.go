package db

import (
	"fmt"
    "os"
	"database/sql"

	_ "modernc.org/sqlite"
	
)


var db *sql.DB
const schema string = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR,
	comment TEXT,
	repeat VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`


func Init(dbFile string) error {
    var install bool
    _, err := os.Stat(dbFile)
    if err != nil {
        install = true
    }

    db, err = sql.Open("sqlite", dbFile)
    if err != nil {
        return fmt.Errorf("database open error: %s", err)
    }

    if install{
        db.Exec(schema)
    }
    return nil
}


func Close() error {
	if db != nil {
		err := db.Close()
		db = nil
		return err
	}
	return nil
}
