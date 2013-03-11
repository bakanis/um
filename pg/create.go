package pg

import (
	"errors"
	"github.com/golibs/um"
	"strings"
	"time"
)

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
