package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/lchsk/rss/user"

	_ "github.com/lib/pq"
)

func main() {
}

type DbAccess struct {
	DB   *sql.DB
	User *user.UserAccess
}

func initDbAccess(db *sql.DB) (*DbAccess, error) {
	ua, err := user.InitUserAccess(db)

	if err != nil {
		return nil, err
	}

	return &DbAccess{DB: db, User: ua}, nil
}

func init() {
	var install = flag.Bool("install", false, "Install new db")
	var demo = flag.Bool("demo", false, "Install demo")

	flag.Parse()

	db, err := getDBConn()

	if err != nil {
		fmt.Printf("cannot get connection: %s\n", err)
		return
	}

	if *install {
		installSchema(db)
	}

	dba, _ := initDbAccess(db)

	if *demo {
		installDemo(dba)
	}

}

func installDemo(dba *DbAccess) {
	ua := dba.User

	ua.InsertUser("bugs", "bugs@bunny.com", "bunny")

	fmt.Println("created demo")
}

func installSchema(db *sql.DB) {
	f, err := ioutil.ReadFile("schema.sql")

	if err != nil {
		fmt.Printf("cannot open schema.sql: %s\n", err)
		return
	}

	schema := string(f)

	if err != nil {
		fmt.Printf("cannot get db conn: %s\n", err)
		return
	}

	_, err = db.Exec(schema)

	if err == nil {
		fmt.Println("created schema successfully")
	} else {
		fmt.Printf("error creating schema: %s\n", err)
	}

}

func getDBConn() (*sql.DB, error) {
	psqlInfo := "host=localhost port=4432 user=rss password=rss dbname=rss sslmode=disable"

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
