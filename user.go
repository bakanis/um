package um

import "time"

type User interface {
	Id() uint64
	DisplayName() string
	EmailAddr() string
	Status() int32
	CreatedOn() time.Time
	LastLogin() time.Time

	SetDisplayName(name string) error
	SetEmailAddr(email string) error
	SetPassword(pw []byte) error
	SetStatus(status int32) error

	Hash() []byte
	Salt() []byte
}
