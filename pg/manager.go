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

	createUserStmt      *sql.Stmt // prepared statement for creating user
	emailAddrExistsStmt *sql.Stmt // statement for checking if a email address exists
	findByIdStmt        *sql.Stmt // statement to look for a user by ID
	findStmt            *sql.Stmt // statement to look for a user by user name and email address
	updateLastLoginStmt *sql.Stmt // statement to update last_login
}

func (this *t_manager) Setup(dns string) error {
	var err error
	this.session, err = sql.Open("postgres", dns)

	if err == nil {
		query := "insert into um_users(email_addr, display_name, status, created_on, last_login) values($1, $2, $3, $4, $5) returning id;"
		this.createUserStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "select exists(select id from um_users where email_addr=$1);"
		this.emailAddrExistsStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "select id, email_addr, display_name, status, hash, salt, created_on, last_login from um_users where id=$1 limit 1;"
		this.findByIdStmt, err = this.session.Prepare(query)
		if err != nil {
			return err
		}
		query = "select id, email_addr, display_name, status, hash, salt, created_on, last_login from um_users where email_addr=$1 limit 1;"
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

func (this *t_manager) CreateUser(emailAddr, displayName string, status int32) (um.User, error) {
	emailAddr = strings.ToLower(strings.Trim(emailAddr, " "))
	displayName = strings.Trim(displayName, " ")
	if emailAddr == "" {
		return nil, errors.New("Email address cannot be empty")
	}
	user := &t_user{
		emailAddr:   emailAddr,
		displayName: displayName,
		status:      status,
		createdOn:   time.Now(),
		lastLogin:   time.Unix(0, 0),
	}

	row := this.createUserStmt.QueryRow(user.emailAddr, user.displayName, user.status, user.createdOn, user.lastLogin)
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
	err := row.Scan(&user.id, &user.emailAddr, &user.displayName, &user.status, &hash, &salt, &user.createdOn, &user.lastLogin)
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
	return rowToUser(this.findStmt.QueryRow(q))
}

func (this *t_manager) EmailAddrExists(emailAddr string) (bool, error) {
	emailAddr = strings.ToLower(strings.Trim(emailAddr, " "))
	if emailAddr == "" {
		return false, errors.New("emailAddr must not be empty")
	}
	var exists bool
	row := this.emailAddrExistsStmt.QueryRow(emailAddr)
	err := row.Scan(&exists)
	return exists, err
}
