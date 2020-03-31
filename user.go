package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

const sqlInsertUser = `
    INSERT INTO users (id, username, email, password) VALUES
    ($1, $2, $3, crypt($4, gen_salt('bf', 8)))
    RETURNING id, username, email, created_at
`

const sqlFindUserByCredentials = `
    SELECT id, username, email, created_at
    FROM users
    WHERE lower(email) = lower($1) AND password = crypt($2, password)
`

const sqlFindUserById = `
    SELECT id, username, email, created_at
    FROM users
    WHERE id = $1
`

type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}

type UserAccess struct {
	Queries map[string]*sql.Stmt
}

func (ua *UserAccess) FindUserByCredentials(email string, password string) (*User, error) {
	u := &User{}

	stmt := ua.Queries["findUserByCredentials"]
	err := stmt.QueryRow(email, password).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ua *UserAccess) FindUserById(id string) (*User, error) {
	u := &User{}

	stmt := ua.Queries["findUserById"]
	err := stmt.QueryRow(id).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ua *UserAccess) insertUser(username string, email string, password string) (*User, error) {
	u := &User{}

	stmt := ua.Queries["insertUser"]

	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	stmt.QueryRow(id, username, email, password).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	return u, err
}

func initUserAccess(db *sql.DB) (*UserAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
		"insertUser":            sqlInsertUser,
		"findUserByCredentials": sqlFindUserByCredentials,
		"findUserById":          sqlFindUserById,
	}

	for name, sql := range queriesToPrepare {
		stmt, err := db.Prepare(sql)

		if err != nil {
			return nil, err
		}

		queries[name] = stmt
	}

	ua := &UserAccess{Queries: queries}

	return ua, nil
}
