package server

import (
	"context"
	"github.com/danielmichaels/conrad/internal/repository"
	"net/http"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

type ctxUserDTO struct {
	ID    int64
	Email string
}

func contextSetAuthenticatedUser(r *http.Request, user *repository.Users) *http.Request {
	u := ctxUserDTO{
		Email: user.Email,
		ID:    user.ID,
	}
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, &u)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *ctxUserDTO {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*ctxUserDTO)
	if !ok {
		return nil
	}
	return user
}
