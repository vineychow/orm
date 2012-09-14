package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"engine/beedb"
	_ "engine/pq"
)

var (
	mux sync.Mutex
	db  *sql.DB
	orm beedb.Model
)

var userTableSql string = `
    create table if not exists tb_user 
    (
        uid serial,
        name varchar(20) not null,
        email varchar(50) not null,
        created varchar(20) not null,
        primary key(uid)
    )
`

func init() {
	mux.Lock()
	defer mux.Unlock()

	// check
	if db != nil {
		return
	}

	psql, err := sql.Open("postgres", "host=192.168.1.138 port=5432 user=postgres password=admin dbname=t2m sslmode=disable")
	checkErr(err)

	// new db
	db = psql

	// new orm
	orm = beedb.New(db, "pg")

	// open debug test
	beedb.OnDebug = true

	// create database table
	_, err = db.Exec(userTableSql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("psql: " + err.Error())
	}
}

// table to struct
type TbUser struct {
	Uid     int `PK`
	Name    string
	Email   string
	Created string
}

// insert data to table
func Insert() error {
	user := &TbUser{
		Name:    "viney",
		Email:   "viney.chow@gmail.com",
		Created: time.Now().Format("2006-01-02 15:04:05"),
	}
	return orm.Save(user)
}

// update data of table
func Update() error {
	user := map[string]interface{}{
		"name":    "viney.chow",
		"created": time.Now().Format("2006-01-02 15:04:05"),
	}

	i, err := orm.SetTable("tb_user").SetPK("uid").Where("uid=$1", 2).Update(user)
	if err == nil {
		fmt.Println(i)
		return nil
	}

	return err
}

// query data from table
func Query() error {
	var user TbUser
	err := orm.Where("name=$1", "viney").Limit(1).Find(&user)
	if err == nil {
		fmt.Println(user)
		return nil
	}

	return err
}

// delete data from table
func Delete() error {
	i, err := orm.SetTable("tb_user").SetPK("uid").Where("name=$1 and uid>$2", "viney", 3).DeleteRow()
	if err == nil {
		fmt.Println(i)
		return nil
	}
	return err
}

func main() {
	// insert
	err := Insert()
	checkErr(err)

	// update
	err = Update()
	checkErr(err)

	// query
	err = Query()
	checkErr(err)

	// delete
	err = Delete()
	checkErr(err)
}
