package um

type UserManager interface {
	Setup(dns string) error

	CreateUser(userName, emailAddr string, status int32) (*User, error)

	FindById(id uint64) (*User, error)
	Find(q string) (*User, error)

	Authenticate(u *User, plainPw string, updateLogin bool) error
	SetPassword(u *User, plainPw string) error
	SetStatus(u *User, status int32) error
}
