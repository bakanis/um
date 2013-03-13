package pg

import "time"

type t_user struct {
	id          uint64
	emailAddr   string
	displayName string
	status      int32
	hash        []byte
	salt        []byte
	createdOn   time.Time
	lastLogin   time.Time
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
	return nil
}

func (this *t_user) SetEmailAddr(email string) error {
	return nil
}

func (this *t_user) SetPassword(pw string) error {
	return nil
}

func (this *t_user) SetStatus(status int32) error {
	return nil
}
