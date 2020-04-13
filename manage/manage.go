package main

import (
	"flag"
	"log"

	"github.com/lchsk/rss/db"
)

func init() {
	var install = flag.Bool("install", false, "Install new db")
	var demo = flag.Bool("demo", false, "Install demo")

	flag.Parse()

	conn, err := db.GetDBConn("rss", "rss", "rss_db", "5432")

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

	if *demo {
		installDemo(DBA)

		log.Println("installed demo")
	}
}

func main() {

}
