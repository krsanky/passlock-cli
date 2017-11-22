package model

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Passlock struct {
	Id        int64
	AccountId int
	Title     string
	Password  string
	Release   time.Time
}

func (p *Passlock) String() string {
	return fmt.Sprintf("<id:%d name:%s>", p.Id, p.Title)
}

func init() {
	open()
	CreateTable()
}

//func open() (*sql.DB, error) {
func open() {
	//return sql.Open("sqlite3", ".passlock.db")
	var err error
	db, err = sql.Open("sqlite3", ".passlock.db")
	if err != nil {
		panic(err)
	}
}

func Close() {
	db.Close()
}

// safe to run if file/table does/doesnt exist
func CreateTable() {
	sql_ := `
CREATE TABLE passlock (                                                                     
    id INTEGER PRIMARY KEY,
    title text NOT NULL,
    password text NOT NULL,
    ts timeSTAMP NOT NULL
)`

	_, err := db.Exec(sql_)

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

func Create(account_id int, title, password string, release time.Time) *Passlock {
	pl := &Passlock{}
	pl.Title = title
	pl.Password = password
	pl.Release = release
	return pl
}

// transaction this
func (p *Passlock) Save() error {
	res, err := db.Exec(`
INSERT INTO passlock
(title, password, ts)
VALUES (?, ?, ?)`,
		p.Title, p.Password, p.Release)

	id, err2 := res.LastInsertId()
	if err2 != nil {
		panic(fmt.Sprintf("Failed to get LastInsertId: %s", err2.Error()))
	}
	p.Id = id

	return err
}

func (p *Passlock) Delete() error {
	_, err := db.Exec(`DELETE FROM passlock WHERE id = ?`, p.Id)
	return err
}

func GetIds() ([]int, error) {
	var (
		id  int
		ids []int
	)

	rows, err := db.Query(`SELECT id from passlock`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func Get(id int) (*Passlock, error) {
	pl := &Passlock{}
	row := db.QueryRow(`
SELECT id, title, password, ts
FROM passlock
WHERE id = $1`, id)

	err := row.Scan(
		&pl.Id,
		&pl.Title,
		&pl.Password,
		&pl.Release)

	return pl, err
}

func GetAll() ([]Passlock, error) {
	ids, err := GetIds()
	if err != nil {
		return nil, err
	}
	pls := make([]Passlock, len(ids))
	for i, id := range ids {
		pl, err := Get(id)
		if err != nil {
			return nil, err
		}
		pls[i] = *pl
	}
	return pls, nil
}
