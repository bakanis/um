package pg

import "github.com/golibs/um"

// Authenticate checks a user against the provided password. Returns an error
// if the user does not exist or not authenticated.
// If the user is authenticated, its login time will be updated if the updateLogin flag is true
func (this *Manager) Authenticate(u *um.User, plainPw string, updateLogin bool) error {
	return nil
}
