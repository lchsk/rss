package main

import "github.com/lchsk/rss/db"

func installDemo(dba *db.DbAccess) {
	ua := dba.User

	ua.InsertUser("bugs", "bugs@bunny.com", "bunny")
}
