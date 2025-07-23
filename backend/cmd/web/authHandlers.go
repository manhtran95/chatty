package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

type SignupResponse struct {
	Form     userSignupForm `json:"form"`
	Redirect bool           `json:"redirect"`
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

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm
	var signupResponse SignupResponse
	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, response.APIResponse[any]{
			Data:    nil,
			Message: "Invalid form data",
			Status:  response.StatusError,
			Error:   response.ErrorBadRequest,
		})
		return
	}
	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	passLen, _ := strconv.Atoi(os.Getenv("PASSWORD_MIN_LENGTH"))
	form.CheckField(validator.MinChars(form.Password, passLen), "password", fmt.Sprintf("This field must be at least %d characters long", passLen))
	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		signupResponse.Form = form
		signupResponse.Redirect = false
		res := response.APIResponse[SignupResponse]{
			Data:    signupResponse,
			Message: "Validation failed",
			Status:  response.StatusError,
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
			signupResponse.Form = form
			signupResponse.Redirect = false
			res := response.APIResponse[SignupResponse]{
				Data:    signupResponse,
				Message: "Duplicate email",
				Status:  response.StatusError,
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
	signupResponse.Redirect = true
	res := response.APIResponse[SignupResponse]{
		Data:    signupResponse,
		Message: "Signup successful",
		Status:  response.StatusSuccess,
	}
	json.NewEncoder(w).Encode(res)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, response.APIResponse[any]{
			Data:    nil,
			Message: "Invalid form data",
			Status:  response.StatusError,
			Error:   response.ErrorBadRequest,
		})
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
				Status:  response.StatusError,
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
		Secure:   os.Getenv("HTTPS") == "true",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((time.Duration(exp) * 24 * time.Hour).Seconds()),
	})

	res := response.APIResponse[LoginResponse]{
		Data:    loginResponse,
		Message: "Login successful",
		Status:  response.StatusSuccess,
	}
	json.NewEncoder(w).Encode(res)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		res := response.APIResponse[any]{
			Data:    nil,
			Message: "Refresh token not found",
			Status:  response.StatusError,
		}
		app.writeJSON(w, http.StatusUnauthorized, res)
		return
	}

	// Validate the refresh token
	userID, err := auth.ValidateRefreshToken(cookie.Value)
	if err != nil {
		// Clear the invalid cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1, // Delete the cookie
		})
		res := response.APIResponse[any]{
			Data:    nil,
			Message: "Invalid refresh token",
			Status:  response.StatusError,
		}
		app.writeJSON(w, http.StatusUnauthorized, res)
		return
	}

	// Generate new access token
	newAccessToken := auth.GenerateAccessToken(userID)

	// Create response
	refreshResponse := RefreshResponse{
		AccessToken: newAccessToken,
	}

	res := response.APIResponse[RefreshResponse]{
		Data:    refreshResponse,
		Message: "Token refreshed successfully",
		Status:  response.StatusSuccess,
	}

	app.writeJSON(w, http.StatusOK, res)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	// Clear the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   os.Getenv("HTTPS") == "true",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete the cookie
	})

	res := response.APIResponse[any]{
		Data:    nil,
		Message: "Logout successful",
		Status:  response.StatusSuccess,
	}

	app.writeJSON(w, http.StatusOK, res)
}
