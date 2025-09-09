package auth

import (
	"context"
	"net/http"

	"github.com/daffadon/graphy/internal/domain/dto"
	"github.com/daffadon/graphy/internal/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	id string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := header
			cl, err := jwt.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := dto.User{ID: cl.UserID}
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *dto.User {
	raw, _ := ctx.Value(userCtxKey).(*dto.User)
	return raw
}
