// +build integration

package main

import (
	"fmt"

	"github.com/lchsk/rss/db"
)

func setupIntegrationTests() {
	conn, err := db.GetDBConn("rsstest", "rsstest", "rss_testdb", "5433")

	if err != nil {
		fmt.Println(fmt.Errorf("getdbconn err: %v", err))
		return
	}

	err = db.InstallSchema(conn, "../schema.sql")

	if err != nil {
		fmt.Println(fmt.Errorf("installschema err: %v", err))
		return
	}

	dba, err := db.InitDbAccess(conn)

	if err != nil {
		fmt.Println(fmt.Errorf("initdbaccess err: %v", err))
		return
	}

	DBA = dba
}
