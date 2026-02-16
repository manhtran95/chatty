package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"chatty.mtran.io/internal/response"
	"github.com/stretchr/testify/assert"
)

const (
	testUserName     = "John Doe"
	testUserEmail    = "john@example.com"
	testUserPassword = "password123"
)

func TestUserSignup(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, _ := setupTestAppWithServer(t)
	defer cleanup()
	defer server.Close()

	// Sign up a user
	resp := signupUser(t, server, testUserName, testUserEmail, testUserPassword)
	defer resp.Body.Close()

	log.Printf("resp: %v", resp)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response response.APIResponse[SignupResponse]
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, response.Message, "Signup successful")
	assert.Equal(t, response.Data.Form.Name, testUserName)
	assert.Equal(t, response.Data.Form.Email, testUserEmail)
	assert.Equal(t, response.Data.Form.Password, testUserPassword)
	assert.True(t, response.Data.Redirect)
}

func TestUserLogin(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, _ := setupTestAppWithServer(t)
	defer cleanup()
	defer server.Close()

	// First, sign up a user
	signupResp := signupUser(t, server, testUserName, testUserEmail, testUserPassword)
	signupResp.Body.Close()
	assert.Equal(t, http.StatusOK, signupResp.StatusCode)

	// Now login with the same credentials
	loginResp := loginUser(t, server, testUserEmail, testUserPassword)
	defer loginResp.Body.Close()

	// Check the response status
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// Parse the response body
	var apiResponse response.APIResponse[LoginResponse]
	err := json.NewDecoder(loginResp.Body).Decode(&apiResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the response
	assert.Equal(t, "Login successful", apiResponse.Message)
	assert.Equal(t, response.StatusSuccess, apiResponse.Status)
	assert.NotNil(t, apiResponse.Data.UserInfo)
	assert.Equal(t, testUserEmail, apiResponse.Data.UserInfo.Email)
	assert.Equal(t, testUserName, apiResponse.Data.UserInfo.Name)
	assert.NotEmpty(t, apiResponse.Data.AccessToken, "Access token should not be empty")

	// Verify refresh token cookie was set
	cookies := loginResp.Cookies()
	var refreshTokenCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshTokenCookie = cookie
			break
		}
	}
	assert.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")
	assert.NotEmpty(t, refreshTokenCookie.Value, "Refresh token should not be empty")
	assert.True(t, refreshTokenCookie.HttpOnly, "Refresh token should be HttpOnly")
}

func TestUserLoginInvalidCredentials(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, _ := setupTestAppWithServer(t)
	defer cleanup()
	defer server.Close()

	// First, sign up a user with a unique email
	testEmail := "invalid-cred-test@example.com"
	signupResp := signupUser(t, server, "Invalid Cred Test User", testEmail, testUserPassword)
	signupResp.Body.Close()
	assert.Equal(t, http.StatusOK, signupResp.StatusCode)

	// Try to login with wrong password
	loginResp := loginUser(t, server, testEmail, "wrongpassword")
	defer loginResp.Body.Close()

	// Should get 422 Unprocessable Entity
	assert.Equal(t, http.StatusUnprocessableEntity, loginResp.StatusCode)

	// Parse the response body
	var apiResponse response.APIResponse[LoginResponse]
	err := json.NewDecoder(loginResp.Body).Decode(&apiResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify error message
	assert.Equal(t, "Login failed", apiResponse.Message)
	assert.Equal(t, response.StatusError, apiResponse.Status)
	assert.Empty(t, apiResponse.Data.AccessToken, "Access token should be empty on failed login")
}

func signupUser(t *testing.T, server *httptest.Server, name string, email string, password string) *http.Response {
	t.Helper()

	// Create form data
	formData := url.Values{}
	formData.Set("name", name)
	formData.Set("email", email)
	formData.Set("password", password)

	// Create the request
	req, err := http.NewRequest(
		http.MethodPost,
		server.URL+"/user/signup", // Use the test server URL
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		t.Fatalf("Failed to create signup request: %v", err)
	}

	// Set the Content-Type header for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send signup request: %v", err)
	}

	return resp
}

func loginUser(t *testing.T, server *httptest.Server, email string, password string) *http.Response {
	t.Helper()

	// Create form data
	formData := url.Values{}
	formData.Set("email", email)
	formData.Set("password", password)

	// Create the request
	req, err := http.NewRequest(
		http.MethodPost,
		server.URL+"/user/login", // Use the login endpoint
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		t.Fatalf("Failed to create login request: %v", err)
	}

	// Set the Content-Type header for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send login request: %v", err)
	}

	return resp
}

func loginUserAndGetAccessToken(t *testing.T, server *httptest.Server, email string, password string) string {
	t.Helper()

	loginResp := loginUser(t, server, email, password)
	defer loginResp.Body.Close()

	var loginResponse response.APIResponse[LoginResponse]
	err := json.NewDecoder(loginResp.Body).Decode(&loginResponse)
	if err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}
	return loginResponse.Data.AccessToken
}