package pg

import (
	"errors"
	"fmt"
	"github.com/golibs/um"
	"strings"
	"time"
)

type t_user struct {
	id          uint64
	emailAddr   string
	displayName string
	status      int32
	hash        []byte
	salt        []byte
	createdOn   time.Time
	lastLogin   time.Time

	manager *t_manager
}

func (this *t_user) Id() uint64 {
	return this.id
}

func (this *t_user) DisplayName() string {
	return this.displayName
}

func (this *t_user) EmailAddr() string {
	return this.emailAddr
}

func (this *t_user) Status() int32 {
	return this.status
}

func (this *t_user) CreatedOn() time.Time {
	return this.createdOn
}

func (this *t_user) LastLogin() time.Time {
	return this.lastLogin
}

func (this *t_user) Hash() []byte {
	return this.hash
}

func (this *t_user) Salt() []byte {
	return this.salt
}

func (this *t_user) SetDisplayName(name string) error {
	name = strings.Trim(name, " ")
	if this.displayName != name {
		this.displayName = name
		_, err := this.manager.setDisplayNameStmt.Exec(name, this.id)
		return err
	}
	return nil
}

func (this *t_user) SetEmailAddr(email string) error {
	email = strings.ToLower(strings.Trim(email, " "))
	if this.emailAddr != email {
		this.emailAddr = email
		_, err := this.manager.setEmailAddrStmt.Exec(email, this.id)
		return err
	}
	return nil
}

func (this *t_user) SetPassword(pw []byte) error {
	if len(pw) == 0 {
		return errors.New("Password cannot be nil")
	}
	this.salt = um.Salt128()
	this.hash = um.EncryptPassword(pw, this.salt)
	_, err := this.manager.setPasswordStmt.Exec(fmt.Sprintf("%02x", this.hash), fmt.Sprintf("%02x", this.salt), this.id)
	return err
}

func (this *t_user) SetStatus(status int32) error {
	if this.status != status {
		this.status = status
		_, err := this.manager.setStatusStmt.Exec(status, this.id)
		return err
	}
	return nil
}
