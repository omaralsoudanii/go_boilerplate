package middleware

import (
	"context"
	"go_boilerplate/lib"
	"go_boilerplate/user"
	"net/http"
	"os"
	"strconv"
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
			return []byte(os.Getenv("SIGNED_ACCESS_TKN_SECRET")), nil
		})

		if err != nil {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidTkn)
			return
		}

		if !token.Valid {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidTkn)
			return
		}
		ctx := prepareCtx(tk, r)
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
			return []byte(os.Getenv("SIGNED_REFRESH_TKN_SECRET")), nil
		})

		if err != nil {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidRefreshTkn)
			return
		}

		if !token.Valid {
			lib.RespondJSON(w, http.StatusForbidden, nil, lib.ErrInvalidRefreshTkn)
			return
		}
		ctx := prepareCtx(tk, r)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}

func prepareCtx(tk *user.Token, r *http.Request) context.Context {
	ctxUnqKey, _ := strconv.Atoi(os.Getenv("CTX_USER_SESSION_KEY"))
	key := &user.ContextKey{
		Key: ctxUnqKey,
	}
	sk := os.Getenv("REDIS_SESSION_KEY") + ":" + tk.ID + ":" + tk.UserName + ":" + tk.Email
	userContext := user.ContextData{
		UserName:   tk.UserName,
		SessionKey: sk,
	}
	ctx := context.WithValue(r.Context(), key, &userContext)
	return ctx
}
