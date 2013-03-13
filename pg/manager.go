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
	findByIdStmt       *sql.Stmt // statement to get a user by ID
}

// Setup prepares the manager's database connections
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
	}

	return err
}

// Close cleans up the database connection of the current manager
func (this *t_manager) Close() error {
	if this.session != nil {
		return this.session.Close()
	}
	return nil
}

// CreateUser creates a user record in the database and returns a User structure
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

// Authenticate checks a user against the provided password. Returns an error
// if the user does not exist or not authenticated.
// If the user is authenticated, its login time will be updated if the updateLogin flag is true
func (this *t_manager) Authenticate(u um.User, plainPw string, updateLogin bool) error {
	return nil
}

func (this *t_manager) FindById(id uint64) (um.User, error) {
	row := this.findByIdStmt.QueryRow(id)
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

func (this *t_manager) Find(q string) (um.User, error) {
	return nil, nil
}

// UserNameExists returns true iff there's a user with this user name (case insensitive).
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
