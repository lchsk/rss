package main

import (
	"flag"
	"fmt"

	"github.com/lchsk/rss/db"
	_ "github.com/lib/pq"
)

func main() {
}

func init() {
	var install = flag.Bool("install", false, "Install new db")
	var demo = flag.Bool("demo", false, "Install demo")

	flag.Parse()

	conn, err := db.GetDBConn("rss", "rss", "rss_db", "5432")

	if err != nil {
		fmt.Printf("cannot get connection: %s\n", err)
		return
	}

	if *install {
		db.InstallSchema(conn, "schema.sql")
	}

	DBA, err := db.InitDbAccess(conn)

	if err != nil {
		fmt.Printf("cannot init db access: %s\n", err)
		return
	}

	if *demo {
		installDemo(DBA)
	}
}

func installDemo(dba *db.DbAccess) {
	ua := dba.User

	ua.InsertUser("bugs", "bugs@bunny.com", "bunny")
}
