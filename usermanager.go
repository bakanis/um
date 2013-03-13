package um

type UserManager interface {
	Setup(dns string) error
	Close() error

	CreateUser(userName, emailAddr string, status int32) (User, error)

	FindById(id uint64) (User, error)
	Find(q string) (User, error)
	UserNameExists(userName string) (bool, error)

	Authenticate(u User, plainPw string, updateLogin bool) error
}
