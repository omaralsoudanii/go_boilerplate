package usecase

import (
	"context"
	"errors"
	"fmt"
	_lib "go_boilerplate/lib"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var log = _lib.GetLogger()

type userUseCase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

// NewUserUseCase will create new an userUseCase object representation of user.UseCase interface
func NewUserUseCase(u user.Repository, timeout time.Duration) user.UseCase {
	return &userUseCase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (u *userUseCase) Register(c context.Context, user *models.User) error {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)
	if err != nil {
		log.Error(err)
		return models.ErrBadParamInput
	}
	hashErr := u.userRepo.Insert(ctx, user)
	if hashErr != nil {
		return hashErr
	}
	return nil
}

func (u *userUseCase) SignIn(c context.Context, user *models.User) (string, string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// get password from user repo to validate against sent one
	userModel, err := u.userRepo.FetchByName(ctx, user.UserName)
	if err != nil {
		return "", "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(user.Password))
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	accessTokenString := generateToken(userModel, false)
	refreshTokenString := generateToken(userModel, true)

	// TODO:: save both access_token and refresh_token in redis to detect token leaks on refreshing the access_token
	err = u.userRepo.StoreSession(ctx, userModel, refreshTokenString)
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

func (u *userUseCase) SignOut(c context.Context) error {
	userContext, ok := c.Value("user").(*user.ContextData)
	if !ok {
		log.Error(errors.New("context_retrieve_user_err"))
		return errors.New("context_retrieve_user_err")
	}
	redisKey := "user:" + userContext.UserName
	err := u.userRepo.DeleteSession(redisKey)
	if err != nil {
		log.Error(err)
		return err
	}
	return err
}

func (u *userUseCase) Refresh(c context.Context, refreshToken string) (string, error) {
	userContext, ok := c.Value("user").(*user.ContextData)
	if !ok {
		log.Error(errors.New("context_retrieve_user_err"))
		return "", errors.New("context_retrieve_user_err")
	}

	userData, err := u.userRepo.GetUser(userContext.UserName)
	if err != nil {
		return "", err
	}
	if userData["refreshToken"] != refreshToken {
		log.Error(errors.New("invalid_refresh_token"))
		return "", errors.New("invalid_refresh_token")
	}
	// access_token token creation
	expireToken := time.Now().Add(time.Minute * 20).Unix()
	tk := &user.Token{
		ID:       userData["id"],
		UserName: userData["user_name"],
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "go_boilerplate",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("123213123123213"))
	return tokenString, nil
}
func generateToken(userModel *models.User, refresh bool) string {
	if refresh {
		// refresh token creation
		rTk := &user.Token{
			ID:       fmt.Sprint(userModel.ID),
			UserName: userModel.UserName,
		}
		refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), rTk)
		refreshTokenString, _ := refreshToken.SignedString([]byte("123213123123213RefreshToken"))
		return refreshTokenString

	} else {
		// access_token token creation
		expireToken := time.Now().Add(time.Minute * 10000).Unix()
		tk := &user.Token{
			ID:       fmt.Sprint(userModel.ID),
			UserName: userModel.UserName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expireToken,
				Issuer:    "go_boilerplate",
				IssuedAt:  time.Now().Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte("123213123123213"))
		return tokenString
	}
}
