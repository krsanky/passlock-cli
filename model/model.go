package model

import (
	"database/sql"
	"regexp"

	"github.com/mattn/go-sqlite3"
)

func open() (*sql.DB, error) {
	return sql.Open("sqlite3", ".passlock.db")
}

// safe to run if file/table does/doesnt exist
func CreateTable() {
	db, err := open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql_ := `
CREATE TABLE passlock (                                                                     
    id INTEGER PRIMARY KEY,
    title text NOT NULL,
    password text NOT NULL,
    ts timeSTAMP NOT NULL
);`

	_, err = db.Exec(sql_)

	if err != nil {
		//fmt.Printf("type: %s\n", reflect.TypeOf(err))
		if _, ok := err.(sqlite3.Error); ok {
			//fmt.Printf("Code: %s\n", e.Code)
			//fmt.Printf("ECode: %s\n", e.ExtendedCode)
			matched, err2 := regexp.MatchString("table .* already exists", err.Error())
			//fmt.Println(matched, err2)
			if (!matched) || (err2 != nil) {
				panic(err)
			}
		} else {
			panic(err.Error())
		}
	}
	//type: sqlite3.Error
	//panic: table passlock already exists
}
