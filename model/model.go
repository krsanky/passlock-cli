package model

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/mattn/go-sqlite3"
)

//var DB *sql.DB

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
	CreateTable()
}

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
)`

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

func Create(account_id int, title, password string, release time.Time) *Passlock {
	pl := &Passlock{}
	pl.Title = title
	pl.Password = password
	pl.Release = release
	return pl
}

// transaction this
func (p *Passlock) Save() error {
	db, err := open()
	defer db.Close()

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
	db, err := open()
	defer db.Close()

	_, err = db.Exec(`
DELETE FROM passlock
WHERE id = ?`,
		p.Id)
	return err
}

/*
func Get(id int) (*Passlock, error) {
	pl := &Passlock{}
	row := db.DB.QueryRow(`
SELECT id, account_id, title, password, ts
FROM passlock
WHERE id = $1`, id)

	err := row.Scan(
		&pl.Id,
		&pl.AccountId,
		&pl.Title,
		&pl.Password,
		&pl.Release)

	return pl, err
}

func GetIds(u *account.User) ([]int, error) {
	var (
		id  int
		ids []int
	)
	rows, err := db.DB.Query(`
SELECT id from passlock
WHERE account_id = $1`, u.Id)

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

func GetAll(u *account.User) ([]Passlock, error) {
	ids, err := GetIds(u)
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
*/
