package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"chatty.mtran.io/internal/auth"
	"chatty.mtran.io/internal/models"
	"chatty.mtran.io/internal/response"
	"chatty.mtran.io/internal/validator"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type LoginResponse struct {
	validator.Validator
	UserInfo    *models.UserInfo `json:"userInfo,omitempty"`
	AccessToken string           `json:"accessToken,omitempty"`
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	userInfo, err := app.users.Authenticate(form.Email, form.Password)
	loginResponse := LoginResponse{
		UserInfo: userInfo,
	}
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			loginResponse.AddNonFieldError("Email or password is incorrect")
			res := response.APIResponse[LoginResponse]{
				Data:    loginResponse,
				Message: "Login failed",
				Status:  "error",
			}
			app.writeJSON(w, http.StatusUnprocessableEntity, res)
		} else {
			app.writeJSONServerError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// JWT token generation
	accessToken := auth.GenerateAccessToken(userInfo.ID.String())
	loginResponse.AccessToken = accessToken
	exp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRE_DAYS"))
	refreshToken := auth.GenerateRefreshToken(userInfo.ID.String())

	// Set refresh token as secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((time.Duration(exp) * 24 * time.Hour).Seconds()),
	})

	res := response.APIResponse[LoginResponse]{
		Data:    loginResponse,
		Message: "Login successful",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(res)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm
	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		res := response.APIResponse[userSignupForm]{
			Data:    form,
			Message: "Validation failed",
			Status:  "error",
		}
		app.writeJSON(w, http.StatusBadRequest, res)
		return
	}
	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			res := response.APIResponse[userSignupForm]{
				Data:    form,
				Message: "Duplicate email",
				Status:  "error",
			}
			app.errorLog.Printf("Error inserting user: %v", err)
			app.writeJSON(w, http.StatusUnprocessableEntity, res)
			return
		} else {
			app.errorLog.Printf("Error inserting user: %v", err)
			app.writeJSONServerError(w, http.StatusInternalServerError, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	// app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// Success response — no redirect, just JSON
	res := response.APIResponse[any]{
		Data:    nil,
		Message: "Signup successful",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(res)
}
