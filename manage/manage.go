package main

import (
	"flag"
	"log"

	"github.com/lchsk/rss/db"
	"github.com/lchsk/rss/demo"
)

func init() {
	var install = flag.Bool("install", false, "Install new db")
	var demoFlag = flag.Bool("demo", false, "Install demo")

	flag.Parse()

	conn, err := db.GetDBConn("localhost", "rss", "rss", "rss_db", "5432")

	if err != nil {
		log.Println("cannot get connection: %s\n", err)
		return
	}

	if *install {
		err := db.InstallSchema(conn, "../sql/schema.sql")

		if err != nil {
			log.Println("could not install schema: ", err)
			return
		}

		log.Println("installed schema")
	}

	DBA, err := db.InitDbAccess(conn)

	if err != nil {
		log.Println("cannot init db access: %s\n", err)
		return
	}

	if *demoFlag {
		demo.InstallDemo(DBA)

		log.Println("installed demo")
	}
}

func main() {

}
