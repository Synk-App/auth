package tests

import (
	"database/sql"
	"strconv"
	"synk/gateway/app"
	"synk/gateway/app/model"
	"testing"
)

func setupUsersModelDB(t *testing.T) *sql.DB {
	db, err := app.InitDB(true)
	if err != nil {
		t.Fatalf("db connection failed: %v", err)
	}
	return db
}

func TestNewUsers(t *testing.T) {
	db := setupUsersModelDB(t)
	defer db.Close()

	u := model.NewUsers(db)
	if u == nil {
		t.Error("NewUsers returned nil")
	}
}

func TestUsers_Add(t *testing.T) {
	db := setupUsersModelDB(t)
	defer db.Close()
	usersModel := model.NewUsers(db)

	input := model.UserRegisterData{
		UserName:  "Test Add User",
		UserEmail: "add_test@synk.com",
		UserPass:  "123456",
	}

	id, err := usersModel.Add(input)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if id == 0 {
		t.Fatal("Add returned ID 0")
	}

	db.Exec("DELETE FROM user WHERE user_id = ?", id)
}

func TestUsers_ById(t *testing.T) {
	db := setupUsersModelDB(t)
	defer db.Close()
	usersModel := model.NewUsers(db)

	input := model.UserRegisterData{
		UserName:  "Test ById",
		UserEmail: "byid_test@synk.com",
		UserPass:  "pass123",
	}
	id, _ := usersModel.Add(input)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", id)

	user, err := usersModel.ById(id)
	if err != nil {
		t.Errorf("ById failed: %v", err)
	}
	if user.UserId != id {
		t.Errorf("ById returned wrong ID. Got %d, want %d", user.UserId, id)
	}
	if user.UserName != input.UserName {
		t.Errorf("ById returned wrong Name. Got %s, want %s", user.UserName, input.UserName)
	}

	emptyUser, err := usersModel.ById(-999)
	if err != nil {
		t.Errorf("ById(invalid) should not error, got: %v", err)
	}
	if emptyUser.UserId != 0 {
		t.Error("ById(invalid) should return empty user (UserId 0)")
	}
}

func TestUsers_ByEmail(t *testing.T) {
	db := setupUsersModelDB(t)
	defer db.Close()
	usersModel := model.NewUsers(db)

	email := "byemail_test@synk.com"
	input := model.UserRegisterData{
		UserName:  "Test ByEmail",
		UserEmail: email,
		UserPass:  "pass456",
	}
	id, _ := usersModel.Add(input)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", id)

	user, err := usersModel.ByEmail(email)
	if err != nil {
		t.Errorf("ByEmail failed: %v", err)
	}
	if user.UserId != id {
		t.Errorf("ByEmail returned wrong user ID")
	}
	if user.UserEmail != email {
		t.Errorf("ByEmail returned wrong Email")
	}

	emptyUser, err := usersModel.ByEmail("non_existent@synk.com")
	if err != nil {
		t.Errorf("ByEmail(invalid) should not error, got: %v", err)
	}
	if emptyUser.UserId != 0 {
		t.Error("ByEmail(invalid) should return empty user")
	}
}

func TestUsers_List(t *testing.T) {
	db := setupUsersModelDB(t)
	defer db.Close()
	usersModel := model.NewUsers(db)

	input := model.UserRegisterData{
		UserName:  "Test List",
		UserEmail: "list_test@synk.com",
		UserPass:  "pass789",
	}
	id, _ := usersModel.Add(input)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", id)

	listAll, err := usersModel.List("")
	if err != nil {
		t.Fatalf("List(all) failed: %v", err)
	}

	found := false
	for _, u := range listAll {
		if u.UserId == id {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created user not found in List(all)")
	}

	listSpecific, err := usersModel.List(strconv.Itoa(id))
	if err != nil {
		t.Fatalf("List(specific) failed: %v", err)
	}
	if len(listSpecific) != 1 {
		t.Errorf("List(specific) expected 1 result, got %d", len(listSpecific))
	} else {
		if listSpecific[0].UserEmail != input.UserEmail {
			t.Error("List(specific) data mismatch")
		}
	}

	listEmpty, err := usersModel.List("-999")
	if err != nil {
		t.Errorf("List(invalid) should not error, got: %v", err)
	}
	if len(listEmpty) != 0 {
		t.Errorf("List(invalid) expected 0 results, got %d", len(listEmpty))
	}
}
