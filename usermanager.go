package um

type UserManager interface {
	// Setup prepares a user manager and database connection
	Setup(dns string) error

	// Close cleans up the user manager and database connection
	Close() error

	// CreateUser creates new user in the database and returns a new user.
	// If any error occures (such as database error or duplicate username),
	// the returning User will be nil along with the error.
	// The user name and email address will be converted to lower case.
	CreateUser(userName, emailAddr string, status int32) (User, error)

	// FindById looks for a User by its ID number.
	// If the user with that ID does not exist, the returning User will be nil
	// along with an error.
	FindById(id uint64) (User, error)

	// Find looks for a User with a matching user name or email address.
	// This operation is case insensitive. If no match is found, the returning User
	// will be nil along with an error.
	Find(q string) (User, error)

	// UserNameExists returns true iff the user name already exists in the database.
	// This operation is case insensitve.
	// If an error occurs during the operation, false will be return along with the error.
	UserNameExists(userName string) (bool, error)

	// Authenticate checks a user against the provided password.
	// This function returns an error if the user does not exist or not authenticated.
	// If the user is authenticated, the function will return nil.
	// The user login time will be updated if the updateLogin flag is true.
	Authenticate(u User, plainPw string, updateLogin bool) error
}
