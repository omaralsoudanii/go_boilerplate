package http

import (
	"context"
	"github.com/asaskevich/govalidator"
	"go_boilerplate/lib"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"net/http"
)

var log = lib.GetLogger()

type NewHttpUserHandler struct {
	UserUseCase user.UseCase
}

type newUserPayLoad struct {
	Tokens tokenPayload
	User *models.User
}

type tokenPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type accessTokenPayload struct {
	AccessToken string `json:"access_token"`
}

func (user *NewHttpUserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userRow models.User
	if err := lib.GetJSON(r, &userRow); err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusUnprocessableEntity, nil, err)
		return
	}

	if ok, err := govalidator.ValidateStruct(&userRow); !ok {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusBadRequest, nil, err)
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	u, accessToken, refreshToken, err := user.UserUseCase.Register(ctx, &userRow)

	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err)
		return
	}
	res := newUserPayLoad{
		Tokens:tokenPayload{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		},
		User:u,
	}

	lib.RespondJSON(w, http.StatusCreated, res, nil)
}

func (user *NewHttpUserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userRow models.User

	if err := lib.GetJSON(r, &userRow); err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusUnprocessableEntity, nil, err)
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	accessToken, refreshToken, err := user.UserUseCase.SignIn(ctx, &userRow)
	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusUnauthorized, nil, lib.ErrNotFound)
		return
	}
	lib.RespondJSON(w, http.StatusOK, tokenPayload{AccessToken: accessToken, RefreshToken: refreshToken}, nil)
}

func (user *NewHttpUserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("refresh_token")
	grantType := r.Header.Get("grant_type")
	if grantType != "refresh_token" {
		lib.RespondJSON(w, http.StatusBadRequest, nil, lib.ErrBadGrantType)
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	accessToken, err := user.UserUseCase.Refresh(ctx, refreshToken)
	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusUnauthorized, nil, lib.ErrInvalidRefreshTkn)
		return
	}
	lib.RespondJSON(w, http.StatusOK, accessTokenPayload{AccessToken: accessToken}, nil)
}

func (user *NewHttpUserHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err := user.UserUseCase.SignOut(ctx)
	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusUnauthorized, nil, err)
		return
	}
	lib.RespondJSON(w, http.StatusOK, nil, nil)
}
