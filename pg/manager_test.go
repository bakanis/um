package pg

import (
	"github.com/golibs/um"
	"math/rand"
	"testing"
	"time"
)

// TestCreateValidUser makes sure a the user manager can handle a valid user creation
func TestCreateValidUser(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)

	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	status := rand.Int31()
	user, err := manager.CreateUser("test@example.com", "Test User", status)
	if err != nil {
		panic(err)
	}
	if user == nil {
		t.Error("CreateUser did not return a valid User structure")
		t.FailNow()
	}
	if user.Id() == 0 {
		t.Error("CreateUser did not give the User structure a valid ID number")
		t.Fail()
	}
	if user.EmailAddr() != "test@example.com" {
		t.Errorf("EmailAddr is '%s' instead of 'test@example.com'", user.EmailAddr())
		t.Fail()
	}
	if user.Status() != status {
		t.Errorf("Status is '%d' instead of '%d' as expected", user.Status(), status)
		t.Fail()
	}

	props := map[string]interface{}{
		"id":           user.Id(),
		"email_addr":   user.EmailAddr(),
		"status":       user.Status(),
		"display_name": user.DisplayName(),
	}
	if !assertRecord(session, "um_users", props) {
		t.Error("Cannot find the coresponding user record in the database")
		t.FailNow()
	}
}

func TestCreateInvalidUser(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.CreateUser("  ", "Some User", 0)
	if user != nil || err == nil {
		t.Error("CreateUser must return nil and an error when the email address is empty")
	}
}

// TestEmailExists makes sure EmailExists returns true iff there's a user record in the database
func TestEmailExists(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	exists, err := manager.EmailAddrExists("fixtureuser1@example.com")
	if err != nil {
		panic(err)
	}
	if exists != true {
		t.Error("fixtureuser1 exists in the database")
	}
	exists, err = manager.EmailAddrExists("FIXTUREUSER2@Example.com")
	if err != nil {
		panic(err)
	}
	if exists != true {
		t.Error("FIXTUREUSER2 (case insensitive) exists in the database")
	}
	exists, err = manager.EmailAddrExists("fixtureuser3doesnotexist@something.com")
	if err != nil {
		panic(err)
	}
	if exists != false {
		t.Error("fixtureuser3doesnotexist does not exist in the database")
	}
}

// TestFindByIdSuccess makes sure FindById returns the correct existing user
func TestFindByIdSuccess(t *testing.T) {
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
	if user == nil {
		t.Error("FindById did not return a valid user for ID #1")
		t.FailNow()
	}
	if id := user.Id(); id != 1 {
		t.Errorf("User ID returns %d, but 1 is expected", id)
		t.Fail()
	}
	if emailAddr := user.EmailAddr(); emailAddr != "fixtureuser1@example.com" {
		t.Errorf("Email address returns '%s', but 'fixtureuser1@example.com' is expected", emailAddr)
		t.Fail()
	}
}

// TestFindByIdFail makes sure FindById returns nil and an error when the user does not exist
func TestFindByIdFail(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1000)
	if !(err != nil && user == nil) {
		t.Error("FindById must return an error and a nil-um.User")
		t.FailNow()
	}
}

func testFindSuccessHelper(t *testing.T, manager um.UserManager, q string, expectedId uint64) {
	user, err := manager.Find(q)
	if err != nil {
		panic(err)
	}
	if user == nil {
		t.Errorf("Find did not return a valid user for '%s'", q)
		t.FailNow()
	}
	if id := user.Id(); id != expectedId {
		t.Errorf("Find for '%s' returned ID #%d, but #%d is expected", q, id, expectedId)
		t.Fail()
	}
}

// TestFindSuccess makes sure Find returns a valid User when user exists
func TestFindSuccess(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	testFindSuccessHelper(t, manager, "fixtureuser1@example.com", 1)
	testFindSuccessHelper(t, manager, "fixtureUSER2@EXAMPLe.COM", 2)
}

// TestFindFailture makes sure Find returns nil & error when user doesn't exist
func TestFindFailure(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.Find("userdoesnotexist")
	if user != nil || err == nil {
		t.Error("Find must return nil and an error when user does not exist")
		t.FailNow()
	}
}

// TestAuthenticateSuccess makes sure authenticate returns nil
// when the password matches
func TestAuthenticateSuccess(t *testing.T) {
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
	if manager.Authenticate(user, "Password123", false) != nil {
		t.Error("Authenticate must return nil when the password matches")
		t.FailNow()
	}
}

func TestAuthenticateFail(t *testing.T) {
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
	if manager.Authenticate(user, "WrongPassword", false) == nil {
		t.Error("Authenticate must return an error when password does not match")
		t.FailNow()
	}
}

func TestAuthenticateUpdateLogin(t *testing.T) {
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
	lastLogin := user.LastLogin()
	time.Sleep(3 * time.Second)
	if manager.Authenticate(user, "Password123", true) != nil {
		t.Error("Authenticate must return nil when the password matches")
		t.FailNow()
	}
	if !user.LastLogin().After(lastLogin) {
		t.Error("Successful authentication must update the User last login when flag is set to true")
		t.Fail()
	}
	if !assertRecord(session, "um_users", map[string]interface{}{"id": user.Id(), "last_login": user.LastLogin()}) {
		t.Error("Successful authentication must update last login within the database when flag is set to true")
		t.Fail()
	}
}
