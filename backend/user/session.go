package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"
)

type session struct {
	sessionToken string
	csrftoken    string
}

var (
	loginSessions   = map[string]session{}
	ErrUnauthorized = errors.New("Unauthorized")
)

func generateToken() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		// Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(2 * time.Minute),
	}
	http.SetCookie(w, cookie)
}

func setCSRFCookie(w http.ResponseWriter, CSRFToken string) {
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    CSRFToken,
		Path:     "/",
		HttpOnly: false,
		// Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(2 * time.Minute),
	}
	http.SetCookie(w, cookie)
}

func Authorize(r *http.Request) error {
	userEmail, err := r.Cookie("authInfo")
	session, ok := loginSessions[userEmail.Value]
	if !ok {
		return ErrUnauthorized
	}
	s, err := r.Cookie("session_id")
	if err != nil || s.Value == "" || s.Value != session.sessionToken {
		return ErrUnauthorized
	}
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != session.csrftoken || csrf == "" {
		return ErrUnauthorized
	}
	return nil
}
