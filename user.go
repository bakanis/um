package um

import "time"

type User interface {
	Id() uint64
	UserName() string
	EmailAddr() string
	Status() int32
	CreatedOn() time.Time
	LastLogin() time.Time

	SetEmailAddr(email string) error
	SetPassword(pw string) error
	SetStatus(status int32) error

	Hash() []byte
	Salt() []byte
}
