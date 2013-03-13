package pg

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/golibs/um"
	"strings"
	"time"
)

type t_manager struct {
	session *sql.DB

	createUserStmt     *sql.Stmt // prepared statement for creating user
	userNameExistsStmt *sql.Stmt // statement for checking if a user name exists
	findByIdStmt       *sql.Stmt // statement to look for a user by ID
	findStmt           *sql.Stmt // statement to look for a user by user name and email address
	updateLastLoginStmt *sql.Stmt // statement to update last_login
}

func (this *t_manager) Setup(dns string) error {
	var err error
	this.session, err = sql.Open("postgres", dns)

	if err == nil {
		query := "insert into um_users(user_name, email_addr, status, created_on, last_login) values($1, $2, $3, $4, $5) returning id;"
		this.createUserStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "select exists(select id from um_users where user_name=$1);"
		this.userNameExistsStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "select id, user_name, email_addr, status, hash, salt, created_on, last_login from um_users where id=$1 limit 1;"
		this.findByIdStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "select id, user_name, email_addr, status, hash, salt, created_on, last_login from um_users where user_name=$1 or email_addr=$2 limit 1;"
		this.findStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "update um_users set last_login=$1 where id=$2;"
		this.updateLastLoginStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
	}

	return err
}

func (this *t_manager) Close() error {
	if this.session != nil {
		return this.session.Close()
	}
	return nil
}

func (this *t_manager) CreateUser(userName, emailAddr string, status int32) (um.User, error) {
	if userName == "" {
		return nil, errors.New("User name must not be blank")
	}
	user := &t_user{
		userName:  strings.ToLower(strings.Trim(userName, " ")),
		emailAddr: strings.ToLower(strings.Trim(emailAddr, " ")),
		status:    status,
		createdOn: time.Now(),
		lastLogin: time.Unix(0, 0),
	}

	row := this.createUserStmt.QueryRow(user.userName, user.emailAddr, user.status, user.createdOn, user.lastLogin)
	err := row.Scan(&user.id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (this *t_manager) Authenticate(u um.User, plainPw string, updateLogin bool) error {
	if u == nil || plainPw == "" {
		return errors.New("User and plain password cannot be nil")
	}

	err := um.ComparePassword(u.Hash(), []byte(plainPw), u.Salt())
	if err == nil && updateLogin {
		lastLogin := time.Now()
		u.(*t_user).lastLogin = lastLogin
		_, err = this.updateLastLoginStmt.Exec(lastLogin, u.Id())
	}

	return err
}

func rowToUser(row *sql.Row) (um.User, error) {
	user := &t_user{}
	var hash, salt string
	err := row.Scan(&user.id, &user.userName, &user.emailAddr, &user.status, &hash, &salt, &user.createdOn, &user.lastLogin)
	if err != nil {
		return nil, err
	}
	user.hash, err = hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	user.salt, err = hex.DecodeString(salt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (this *t_manager) FindById(id uint64) (um.User, error) {
	return rowToUser(this.findByIdStmt.QueryRow(id))
}

func (this *t_manager) Find(q string) (um.User, error) {
	q = strings.ToLower(strings.Trim(q, " "))
	return rowToUser(this.findStmt.QueryRow(q, q))
}

func (this *t_manager) UserNameExists(userName string) (bool, error) {
	userName = strings.ToLower(strings.Trim(userName, " "))
	if userName == "" {
		return false, errors.New("userName must not be empty")
	}
	var exists bool
	row := this.userNameExistsStmt.QueryRow(userName)
	err := row.Scan(&exists)
	return exists, err
}
