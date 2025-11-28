package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"synk/gateway/app"
	"synk/gateway/app/controller"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func setupUserCtrlDB(t *testing.T) *sql.DB {
	db, err := app.InitDB(true)
	if err != nil {
		t.Fatalf("controller: db connection failed [%v]", err.Error())
	}

	os.Setenv("JWT_SECRET", "test_secret_key_123")
	os.Setenv("JWT_REFRESH_SECRET", "test_refresh_secret_key_456")
	return db
}

func createDummyUserForController(db *sql.DB) (int, string, string) {
	email := "controller_auth_test@synk.com"
	rawPass := "password123"

	hashed, _ := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)

	res, _ := db.ExecContext(context.Background(),
		"INSERT INTO user (user_name, user_email, user_pass) VALUES (?, ?, ?)",
		"Auth Controller User", email, string(hashed))

	id, _ := res.LastInsertId()
	return int(id), email, rawPass
}

func generateTestAccessToken(userId int) string {
	claims := controller.AccessTokenClaims{}
	claims.User.UserId = userId
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return s
}

func generateTestRefreshToken(userId int) string {
	claims := controller.RefreshTokenClaims{}
	claims.User.UserId = userId
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	return s
}

func TestUsers_HandleShow(t *testing.T) {
	db := setupUserCtrlDB(t)
	defer db.Close()

	uid, _, _ := createDummyUserForController(db)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", uid)

	uController := controller.NewUsers(db)

	req, _ := http.NewRequest("GET", "/users?user_id="+strconv.Itoa(uid), nil)
	rr := httptest.NewRecorder()

	uController.HandleShow(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleShow status error: %d", rr.Code)
	}

	var resp controller.HandleShowResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if len(resp.Data) == 0 {
		t.Error("HandleShow returned empty list")
	}
}

func TestUsers_HandleRegister(t *testing.T) {
	db := setupUserCtrlDB(t)
	defer db.Close()

	uController := controller.NewUsers(db)

	email := "new_register_controller@synk.com"

	db.Exec("DELETE FROM user WHERE user_email = ?", email)

	reqBody := controller.HandleUserRegisterRequest{
		UserName:  "New Registrant",
		UserEmail: email,
		UserPass:  "pass123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	uController.HandleRegister(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Register failed: %v Body: %s", rr.Code, rr.Body.String())
	}

	var resp controller.HandleUserRegisterResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp.Data.UserId == 0 {
		t.Error("Register returned 0 ID")
	}
	if resp.Data.Token == "" {
		t.Error("Register did not return access token")
	}

	cookies := rr.Result().Cookies()
	foundRefresh := false
	for _, c := range cookies {
		if c.Name == "refresh_token" && c.Value != "" {
			foundRefresh = true
		}
	}
	if !foundRefresh {
		t.Error("Register did not set refresh_token cookie")
	}

	db.Exec("DELETE FROM user WHERE user_id = ?", resp.Data.UserId)
}

func TestUsers_HandleLogin(t *testing.T) {
	db := setupUserCtrlDB(t)
	defer db.Close()

	uid, email, pass := createDummyUserForController(db)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", uid)

	uController := controller.NewUsers(db)

	reqBody := controller.HandleUserLoginRequest{
		UserEmail: email,
		UserPass:  pass,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	uController.HandleLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Login failed: %v. Body: %s", rr.Code, rr.Body.String())
	}

	var resp controller.HandleUserLoginResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp.Data.UserId != uid {
		t.Errorf("Login returned wrong user ID. Got %d, want %d", resp.Data.UserId, uid)
	}
	if resp.Data.Token == "" {
		t.Error("Login token empty")
	}
}

func TestUsers_HandleRefresh(t *testing.T) {
	db := setupUserCtrlDB(t)
	defer db.Close()

	uid, _, _ := createDummyUserForController(db)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", uid)

	uController := controller.NewUsers(db)

	refreshToken := generateTestRefreshToken(uid)

	req, _ := http.NewRequest("POST", "/users/refresh", nil)

	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	rr := httptest.NewRecorder()

	uController.HandleRefresh(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Refresh failed: %v Body: %s", rr.Code, rr.Body.String())
	}

	var resp controller.HandleUserLoginResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp.Data.Token == "" {
		t.Error("Refresh did not issue new access token")
	}
}

func TestUsers_HandleCheck(t *testing.T) {
	db := setupUserCtrlDB(t)
	defer db.Close()

	uid, _, _ := createDummyUserForController(db)
	defer db.Exec("DELETE FROM user WHERE user_id = ?", uid)

	uController := controller.NewUsers(db)

	validToken := generateTestAccessToken(uid)
	req, _ := http.NewRequest("GET", "/users/check", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rr := httptest.NewRecorder()

	uController.HandleCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Check failed for valid token: %v Body: %s", rr.Code, rr.Body.String())
	}

	var resp controller.HandleUserCheckResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp.Data.UserId != uid {
		t.Errorf("Check returned wrong user ID. Got %d, want %d", resp.Data.UserId, uid)
	}

	req2, _ := http.NewRequest("GET", "/users/check", nil)
	req2.Header.Set("Authorization", "Bearer invalid.token.here")
	rr2 := httptest.NewRecorder()

	uController.HandleCheck(rr2, req2)

	if rr2.Code != http.StatusUnauthorized {
		t.Errorf("Check should fail for invalid token but got code: %v", rr2.Code)
	}
}

func TestUsers_HandleLogout(t *testing.T) {
	db := setupUserCtrlDB(t)
	defer db.Close()

	uController := controller.NewUsers(db)

	req, _ := http.NewRequest("POST", "/users/logout", nil)
	rr := httptest.NewRecorder()

	uController.HandleLogout(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Logout status error: %v", rr.Code)
	}

	cookies := rr.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			found = true
			if c.Value != "" {
				t.Error("Logout did not clear cookie value")
			}
			if c.MaxAge >= 0 {
				t.Error("Logout cookie should have negative MaxAge (expiry)")
			}
		}
	}
	if !found {
		t.Error("Logout response did not contain set-cookie header")
	}
}
