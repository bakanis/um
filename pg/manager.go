package pg

import (
	"database/sql"
	"errors"
	"github.com/golibs/um"
	"strings"
	"time"
)

type Manager struct {
	session *sql.DB

	createUserStmt *sql.Stmt // prepared statement for creating user
}

// Setup prepares the manager's database connections
func (this *Manager) Setup(dns string) error {
	var err error
	this.session, err = sql.Open("postgres", dns)

	if err == nil {
		query := "insert into um_users(user_name, email_addr, status, created_on, last_login) values($1, $2, $3, $4, $5) returning id;"
		this.createUserStmt, err = this.session.Prepare(query)
	}

	return err
}

// Close cleans up the database connection of the current manager
func (this *Manager) Close() error {
	if this.session != nil {
		return this.session.Close()
	}
	return nil
}

// CreateUser creates a user record in the database and returns a User structure
func (this *Manager) CreateUser(userName, emailAddr string, status int32) (*um.User, error) {
	if userName == "" {
		return nil, errors.New("User name must not be blank")
	}
	user := &um.User{
		UserName:  strings.ToLower(strings.Trim(userName, " ")),
		EmailAddr: strings.ToLower(strings.Trim(emailAddr, " ")),
		Status:    status,
		CreatedOn: time.Now(),
		LastLogin: time.Unix(0, 0),
	}

	row := this.createUserStmt.QueryRow(user.UserName, user.EmailAddr, user.Status, user.CreatedOn, user.LastLogin)
	err := row.Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Authenticate checks a user against the provided password. Returns an error
// if the user does not exist or not authenticated.
// If the user is authenticated, its login time will be updated if the updateLogin flag is true
func (this *Manager) Authenticate(u *um.User, plainPw string, updateLogin bool) error {
	return nil
}

func (this *Manager) FindById(id uint64) (*um.User, error) {
	return nil, nil
}

func (this *Manager) Find(q string) (*um.User, error) {
	return nil, nil
}

func (this *Manager) UserNameExists(userName string) (bool, error) {
	return false, nil
}

func (this *Manager) SetPassword(u *um.User, plainPw string) error {
	return nil
}

func (this *Manager) SetStatus(u *um.User, status int32) error {
	return nil
}
