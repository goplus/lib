package main

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/os"
	"github.com/goplus/lib/c/sqlite"
)

func main() {
	os.Remove(c.Str("test.db"))

	db, err := sqlite.Open(c.Str("test.db"))
	check(err, db, "sqlite: Open")

	err = db.Exec(c.Str("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)"), nil, nil, nil)
	check(err, db, "sqlite: Exec CREATE TABLE")

	stmt, err := db.PrepareV3("INSERT INTO users (id, name) VALUES (?, ?)", 0, nil)
	check(err, db, "sqlite: PrepareV3 INSERT")

	stmt.BindInt(1, 100)
	stmt.BindText(2, c.Str("Hello World"), -1, nil)

	err = stmt.Step()
	checkDone(err, db, "sqlite: Step INSERT 1")

	stmt.Reset()
	stmt.BindInt(1, 200)
	stmt.BindText(2, c.Str("This is llgo"), -1, nil)

	err = stmt.Step()
	checkDone(err, db, "sqlite: Step INSERT 2")

	stmt.Close()

	stmt, err = db.PrepareV3("SELECT * FROM users", 0, nil)
	check(err, db, "sqlite: PrepareV3 SELECT")

	for {
		if err = stmt.Step(); err != sqlite.HasRow {
			break
		}
		c.Printf(c.Str("==> id=%d, name=%s\n"), stmt.ColumnInt(0), stmt.ColumnText(1))
	}
	checkDone(err, db, "sqlite: Step done")

	stmt.Close()
	db.Close()
}

func check(err sqlite.Errno, db *sqlite.Sqlite3, at string) {
	if err != sqlite.OK {
		c.Printf(c.Str("==> %s Error: (%d) %s\n"), c.AllocaCStr(at), err, db.Errmsg())
		c.Exit(1)
	}
}

func checkDone(err sqlite.Errno, db *sqlite.Sqlite3, at string) {
	if err != sqlite.Done {
		check(err, db, at)
	}
}
