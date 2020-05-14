package main

import (
	"flag"
	"log"

	"github.com/lchsk/rss/demo"
	"github.com/lchsk/rss/libs/db"
)

func main() {
	var install = flag.Bool("install", false, "Install new db")
	var demoFlag = flag.Bool("demo", false, "Install demo")
	var migrate = flag.Bool("migrate", false, "Migrate SQL")

	flag.Parse()

	conn, err := db.GetDBConn("localhost", "rss", "rss", "rss_db", "5432")

	if err != nil {
		log.Printf("cannot get connection: %s\n", err)
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
		log.Printf("cannot init db access: %s\n", err)
		return
	}

	if *migrate {
		log.Println("running migrations")
		Migrate(DBA)
		log.Println("migrations finished")
	}

	if *demoFlag {
		demo.InstallDemo(DBA)

		log.Println("installed demo")
	}
}
