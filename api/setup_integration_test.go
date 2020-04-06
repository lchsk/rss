// +build integration

package main

import (
	"database/sql"
	"fmt"

	"github.com/lchsk/rss/cache"
	"github.com/lchsk/rss/db"
)

func setupSchema(conn *sql.DB) error {
	err := db.InstallSchema(conn, "../schema.sql")

	if err != nil {
		return fmt.Errorf("installschema err %s", err)
	}

	return nil
}

func setupIntegrationTests() error {
	conn, err := db.GetDBConn("rsstest", "rsstest", "rss_testdb", "5433")

	if err != nil {
		return fmt.Errorf("getdbconn err: %v", err)
	}

	err = setupSchema(conn)

	if err != nil {
		return fmt.Errorf("unable to set schema up: %s", err)
	}

	dba, err := db.InitDbAccess(conn)

	if err != nil {
		return fmt.Errorf("initdbaccess err: %v", err)
	}

	DBA = dba

	redis, err := cache.GetRedisConn()

	if err != nil {
		return fmt.Errorf("getredisconn err: %v", err)
	}

	cacheAccess, err := cache.InitCacheAccess(redis)

	if err != nil {
		return fmt.Errorf("initcacheaccess err: %v", err)
	}

	Cache = cacheAccess

	return nil
}
