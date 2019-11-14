package middleware

import (
	"context"
	"go_boilerplate/lib"
	"go_boilerplate/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AccessTokenVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidTkn)
			return
		}

		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidTkn)
			return
		}

		tokenPart := parts[1]
		tk := &user.Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			var err error
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, err
			}
			return []byte("123213123123213"), nil
		})

		if err != nil {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidTkn)
			return
		}

		if !token.Valid {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidTkn)
			return
		}
		userContext := user.ContextData{
			UserName: tk.UserName,
		}
		ctx := context.WithValue(r.Context(), "user", &userContext)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func RefreshTokenVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("refresh_token")

		if tokenHeader == "" {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidRefreshTkn)
			return
		}

		tk := &user.Token{}
		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			var err error
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, err
			}
			return []byte("123213123123213RefreshToken"), nil
		})

		if err != nil {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidRefreshTkn)
			return
		}

		if !token.Valid {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidRefreshTkn)
			return
		}
		userContext := user.ContextData{
			UserName: tk.UserName,
		}
		ctx := context.WithValue(r.Context(), "user", &userContext)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
