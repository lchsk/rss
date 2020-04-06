package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lchsk/rss/cache"
	"github.com/lchsk/rss/db"
	"github.com/lchsk/rss/libs/api"
)

const (
	// URLs
	registerUserUrl   = "/api/users"
	authenticationUrl = "/api/authentication"
	refreshTokenUrl   = "/api/authentication/refresh"
	logoutUrl         = "/api/logout"
	fetchUserUrl      = "/api/users/{user_id}"

	// API Errors
	errInvalidUsernameLen = "invalid_username_len"
	errInvalidPasswordLen = "invalid_password_len"
	errInvalidEmailLen    = "invalid_email_len"
	errInvalidEmail       = "invalid_email"
	errInvalidInputFormat = "invalid_input_format"
	// Generic database error, such as violated index
	errDbError = "db_error"
)

var DBA *db.DbAccess
var Cache *cache.CacheAccess

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(registerUserUrl, handlerRegisterUser).Methods(http.MethodPost)
	router.HandleFunc(authenticationUrl, handlerAuthentication).Methods(http.MethodPost)
	router.HandleFunc(logoutUrl, checkValidToken(handlerLogout)).Methods(http.MethodPost)
	router.HandleFunc(refreshTokenUrl, Refresh).Methods(http.MethodPost)
	router.HandleFunc(fetchUserUrl, checkValidToken(handlerFetchUser)).Methods(http.MethodGet)

	return router
}

func runAPI() {
	router := getRouter()

	log.Fatal(http.ListenAndServe(":8000", api.CommonMiddleware(router)))
}

func main() {
	conn, _ := db.GetDBConn("rss", "rss", "rss_db", "5432")
	dba, _ := db.InitDbAccess(conn)
	DBA = dba

	redis, _ := cache.GetRedisConn()
	cacheAccess, _ := cache.InitCacheAccess(redis)
	Cache = cacheAccess

	runAPI()
}
