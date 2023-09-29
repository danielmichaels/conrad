package server

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

// Ptr takes in non-pointer and returns a pointer
func Ptr[T any](v T) *T {
	return &v
}

func convStrBoolToBool(v string) bool {
	return v == "true"
}

func formatURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

// isOnCheckbox checks the value of a form.Insecure and returns a bool as a string
func isOnCheckbox(v string) string {
	if v == "on" || v == "true" {
		return "true"
	}
	return "false"
}

func urlIDParam(r *http.Request, optName *string) (int64, error) {
	key := "id"
	if optName != nil {
		key = *optName
	}
	param := chi.URLParam(r, key)
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, err
}

func Hash(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Matches(plaintextPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
