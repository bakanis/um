package pg

import (
	"bytes"
	"fmt"
	"github.com/golibs/rand"
	"github.com/golibs/um"
	"testing"
)

func TestSetDisplayName(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1)
	if err != nil {
		panic(err)
	}

	err = user.SetDisplayName("A New Name")
	if err != nil {
		panic(err)
	}
	if name := user.DisplayName(); name != "A New Name" {
		t.Errorf("Display name returns '%s', but 'A New Name' is expected", name)
		t.Fail()
	}
	if !assertRecord(session, "um_users", map[string]interface{}{"id": 1, "display_name": "A New Name"}) {
		t.Error("Display name is not updated in the database")
		t.Fail()
	}
}

func TestSetEmailAddr(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1)
	if err != nil {
		panic(err)
	}

	err = user.SetEmailAddr("newemail@example.com")
	if err != nil {
		panic(err)
	}
	if emailAddr := user.EmailAddr(); emailAddr != "newemail@example.com" {
		t.Errorf("Email address returns '%s', but 'newemail@example.com' is expected", emailAddr)
		t.Fail()
	}
	if !assertRecord(session, "um_users", map[string]interface{}{"id": 1, "email_addr": "newemail@example.com"}) {
		t.Error("Email address is not updated in the database")
		t.Fail()
	}
}

func TestSetStatus(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1)
	if err != nil {
		panic(err)
	}

	status := int32(rand.Int63())

	err = user.SetStatus(status)
	if err != nil {
		panic(err)
	}
	if newStatus := user.Status(); newStatus != status {
		t.Errorf("Status returns '%d', but '%d' is expected", newStatus, status)
		t.Fail()
	}
	if !assertRecord(session, "um_users", map[string]interface{}{"id": 1, "status": status}) {
		t.Error("Email address is not updated in the database")
		t.Fail()
	}
}

func TestSetPassword(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1)
	if err != nil {
		panic(err)
	}

	oldHash := user.Hash()
	oldSalt := user.Salt()
	err = user.SetPassword([]byte("NewPassword1234"))
	if err != nil {
		panic(err)
	}
	err = manager.Authenticate(user, []byte("NewPassword1234"), false)
	if err != nil {
		t.Error("Fail to authenticate with the new password")
		t.Fail()
	}
	if hash := user.Hash(); bytes.Equal(hash, oldHash) {
		t.Error("SetPassword did not modify hash")
		t.Fail()
	}
	if salt := user.Salt(); bytes.Equal(salt, oldSalt) {
		t.Error("SetPassword did not generate a new salt")
		t.Fail()
	}
	if !assertRecord(session, "um_users", map[string]interface{}{"id": 1, "hash": fmt.Sprintf("%02x", user.Hash()), "salt": fmt.Sprintf("%02x", user.Salt())}) {
		t.Error("SetPassword did not update the database")
		t.Fail()
	}
}
