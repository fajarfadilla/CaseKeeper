package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

var (
	ErrRequestEmpty        = errors.New("email and password cannot be empty")
	ErrEmailCannotEmpty    = errors.New("email cannot be empty")
	ErrPasswordCannotEmpty = errors.New("password cannot be empty")
	ErrEmailAlreadyExist   = errors.New("Email Already Exist")
)

func Router() *chi.Mux {
	r := chi.NewMux()
	r.Post("/register", createUserHandler)
	r.Post("/login", loginUserHandler)
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		if err := Authorize(r); err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		writeMessage(w, http.StatusOK, "success")
	})

	return r
}

func writeMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(j)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeMessage(w, status, err.Error())
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	email := r.FormValue("email")
	password := r.FormValue("password")
	switch {
	case email == "" && password == "":
		writeError(w, http.StatusBadRequest, ErrRequestEmpty)
		return
	case email == "":
		writeError(w, http.StatusBadRequest, ErrEmailCannotEmpty)
		return
	case password == "":
		writeError(w, http.StatusBadRequest, ErrPasswordCannotEmpty)
		return
	}
	id, err := createUser(ctx, email, password)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			writeError(w, http.StatusConflict, err)
			return
		}
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var resp struct {
		Id string `json:"id"`
	}
	resp.Id = id.String()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	email := r.FormValue("email")
	password := r.FormValue("password")
	switch {
	case email == "" && password == "":
		writeError(w, http.StatusBadRequest, ErrRequestEmpty)
		return
	case email == "":
		writeError(w, http.StatusBadRequest, ErrEmailCannotEmpty)
		return
	case password == "":
		writeError(w, http.StatusBadRequest, ErrPasswordCannotEmpty)
		return
	}
	user, err := findUser(ctx, email)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	validPassword, err := verifyPassword(password, user.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if !validPassword {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	authInfo := &http.Cookie{
		Name:     "authInfo",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: false,
		// Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(2 * time.Minute),
	}
	http.SetCookie(w, authInfo)

	loginSession := loginSessions[user.Email]
	sessionID, _ := generateToken()
	csrfToken, _ := generateToken()
	setSessionCookie(w, sessionID)
	setCSRFCookie(w, csrfToken)
	loginSession.sessionToken = sessionID
	loginSession.csrftoken = csrfToken
	loginSessions[user.Email] = loginSession

	var resp struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	resp.Id = user.Id.String()
	resp.Email = user.Email
	resp.Role = user.Role.String
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
