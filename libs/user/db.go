package user

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"-"`
	CreatedAt time.Time `json:"-"`
}

type UserAccess struct {
	Db      *sql.DB
	SQ      *sq.StatementBuilderType
	DbCache sq.DBProxy
	Queries map[string]*sql.Stmt
}

type TokenData struct {
	AccessToken      string
	RefreshToken     string
	AccessUuid       string
	RefreshUuid      string
	AccessExpiresAt  int64
	RefreshExpiresAt int64
}

func (ua *UserAccess) FindUserByCredentials(email string, password string) (*User, error) {
	u := &User{}

	userQuery := ua.SQ.Select("id, username, email, created_at").From("users").Where(
		"lower(email) = lower(?) AND password = crypt(?, password)", email, password).Limit(1)
	err := userQuery.RunWith(ua.Db).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return u, nil
}

func (ua *UserAccess) FindUserById(id string) (*User, error) {
	u := &User{}

	userQuery := ua.SQ.Select("id, username, email, created_at").From("users").Where(sq.Eq{
		"id": id,
	}).Limit(1)
	err := userQuery.RunWith(ua.Db).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (ua *UserAccess) InsertUser(username string, email string, password string) (*User, error) {
	u := &User{}

	id := uuid.New()

	query:= ua.SQ.
		Insert("users").Columns("id", "username", "email", "password").
		Values(id, username, email, sq.Expr("crypt(?, gen_salt('bf', 8))", password)).Suffix("RETURNING id, username, email, created_at")

	err := query.RunWith(ua.Db).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	// TODO: Log postgres error

	return u, err
}

func (ua *UserAccess) UpdateUserPassword(id string, newPassword string) error {
	query:= ua.SQ.
		Update("users").Set("password", sq.Expr("crypt(?, gen_salt('bf', 8))", newPassword)).Where("id = ?", id)

	_, err := query.RunWith(ua.Db).Exec()

	// TODO: Log postgres error

	return err
}

func InitUserAccess(db *sql.DB, DbCache sq.DBProxy, psql *sq.StatementBuilderType) (*UserAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
	}

	for name, sql := range queriesToPrepare {
		stmt, err := db.Prepare(sql)

		if err != nil {
			return nil, err
		}

		queries[name] = stmt
	}

	ua := &UserAccess{Db: db, SQ: psql, DbCache: DbCache, Queries: queries}

	return ua, nil
}
