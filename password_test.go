package um

import (
	"github.com/golibs/rand"
	"testing"
)

// TestSalt128 makes sure the result is always 32 hex characters
func TestSalt128(t *testing.T) {
	var salt []byte
	for i := 0; i < 1024; i++ {
		salt = Salt128()
		if len(salt) != 16 {
			t.Errorf("Salt128 fail to generate 128 bit salt: %02x", salt)
			t.FailNow()
		}
	}
}

// TestPasswordSanity makes sure EncryptPassword and ComparePassword work well together
func TestPasswordSanity(t *testing.T) {
	pw := make([]byte, 10)
	pw2 := make([]byte, 10)
	var salt, hash, hash2 []byte

	for i := 0; i < 32; i++ {
		rand.Read(pw)
		rand.Read(pw2)
		salt = Salt128()
		hash = EncryptPassword(pw, salt)
		hash2 = EncryptPassword(pw2, salt)

		if err := ComparePassword(hash, pw, salt); err != nil {
			t.Error(err)
			t.FailNow()
		}
		if err := ComparePassword(hash2, pw, salt); err == nil {
			t.Error("ComparePassword is giving false positive")
			t.FailNow()
		}
	}
}
