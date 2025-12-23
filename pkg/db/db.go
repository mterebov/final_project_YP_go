package db

import (
	"fmt"
    "os"
	"database/sql"

	_ "modernc.org/sqlite"
	
)


var DB *sql.DB
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

    DB, err = sql.Open("sqlite", dbFile)
    if err != nil {
        return fmt.Errorf("database open error: %s", err)
    }

    if install{
        DB.Exec(schema)
    }
    return nil
}


func Close() error {
	if DB != nil {
		err := DB.Close()
		DB = nil
		return err
	}
	return nil
}
