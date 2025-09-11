package auth

import (
	"context"
	"encoding/json"
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

			tokenStr := ""
			if len(header) > 7 && header[:7] == "Bearer " {
				tokenStr = header[7:]
			} else {
				tokenStr = header
			}
			cl, err := jwt.ValidateToken(tokenStr)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				response := map[string]interface{}{
					"errors": []map[string]interface{}{
						{
							"message": "Invalid token",
							"extensions": map[string]interface{}{
								"code": "FORBIDDEN",
							},
						},
					},
					"data": nil,
				}

				json.NewEncoder(w).Encode(response)
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
