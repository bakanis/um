package um

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/golibs/rand"
)

// Salt128 generates 128 bits of random data.
func Salt128() []byte {
	x := make([]byte, 16)
	rand.Read(x)
	return x
}

// makePassword makes the actual password from the plain password and salt
func makePassword(plainPw, salt []byte) []byte {
	password := make([]byte, 0, len(plainPw)+len(salt))
	password = append(password, salt...)
	password = append(password, plainPw...)
	return password
}

// EncryptPassword uses bcrypt to encrypt a password and salt combination.
// It returns the encrypted password in hex form.
func EncryptPassword(plainPw, salt []byte) []byte {
	hash, _ := bcrypt.GenerateFromPassword(makePassword(plainPw, salt), 10)
	return hash
}

// ComparePassword compares a hash with the plain password and the salt.
// The function returns nil on success or an error on failure.
func ComparePassword(hash, plainPw, salt []byte) error {
	return bcrypt.CompareHashAndPassword(hash, makePassword(plainPw, salt))
}
