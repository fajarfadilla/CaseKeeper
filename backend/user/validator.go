package user

import (
	"errors"
	"net"
	"net/mail"
	"strings"
	"unicode"
)

func validatePassword(password string) error {
	isMoreThan8 := len(password) < 8
	if isMoreThan8 {
		return errors.New("Pasword less than 8")
	}
	var isLower, isUpper, isSym, isDigit bool
	for _, r := range password {
		switch {
		case !isLower && unicode.IsLower(r):
			isLower = true
		case !isUpper && unicode.IsUpper(r):
			isUpper = true
		case !isDigit && unicode.IsDigit(r):
			isDigit = true
		case !isSym && unicode.IsSymbol(r) || unicode.IsPunct(r):
			isSym = true
		}
	}
	if !isLower {
		return errors.New("password must contains at least one Lowercase character")
	}
	if !isUpper {
		return errors.New("password must contains at least one Uppercase character")
	}
	if !isDigit {
		return errors.New("password must contains at least one Digit character")
	}
	if !isSym {
		return errors.New("password must contains at least one Symbol character")
	}
	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("Email Invalid")
	}
	emailParts := strings.Split(email, "@")
	_, err = net.LookupMX(emailParts[1])
	if err != nil {
		return errors.New("Domain not found")
	}
	return nil
}
