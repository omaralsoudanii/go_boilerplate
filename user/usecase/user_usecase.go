package usecase

import (
	"context"
	"fmt"
	"go_boilerplate/lib"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *userUseCase) Register(c context.Context, user *models.User) (*models.User, error) {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)
	if err != nil {
		return nil, err
	}
	id, err := u.userRepo.Insert(ctx, user)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (u *userUseCase) SignIn(c context.Context, data *models.User) (string, string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// get password from user repo to validate against sent one
	userModel, err := u.userRepo.GetByName(ctx, data.UserName)
	if err != nil {
		return "", "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(data.Password))
	if err != nil {
		return "", "", err
	}

	accessTokenString, err := generateAccessToken(userModel)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := generateRefreshToken(userModel)
	if err != nil {
		return "", "", err
	}
	// TODO:: save both access_token and refresh_token in redis to detect token leaks on refreshing the access_token
	sk := os.Getenv("REDIS_SESSION_KEY") + ":" + string(userModel.ID) + ":" + userModel.UserName + ":" + userModel.Email

	err = u.userRepo.StoreSession(ctx, userModel, sk, refreshTokenString)
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

func (u *userUseCase) SignOut(c context.Context) error {
	sk, err := getCtxSessionKey(c)
	if err != nil {
		return err
	}

	err = u.userRepo.DeleteSession(sk)
	if err != nil {
		return err
	}
	return err
}

func (u *userUseCase) Refresh(c context.Context, refreshToken string) (string, error) {
	sk, err := getCtxSessionKey(c)
	if err != nil {
		return "", err
	}
	userData, err := u.userRepo.GetUser(sk)
	if err != nil {
		return "", err
	}
	if userData["refreshToken"] != refreshToken {
		return "", lib.ErrInvalidRefreshTkn
	}
	// access_token token creation
	id, _ := strconv.Atoi(userData["id"])
	userModel := &models.User{
		ID:       uint(id),
		UserName: userData["username"],
	}
	tokenString, err := generateAccessToken(userModel)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func generateRefreshToken(userModel *models.User) (string, error) {
	// refresh token creation
	tknTimeout, _ := strconv.Atoi(os.Getenv("SIGNED_REFRESH_TKN_TIMEOUT"))
	expireToken := time.Now().Add(time.Second * time.Duration(tknTimeout)).Unix()
	rTk := &user.Token{
		ID:       fmt.Sprint(userModel.ID),
		UserName: userModel.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    os.Getenv("TKNS_ISSUER"),
			IssuedAt:  time.Now().Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), rTk)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SIGNED_REFRESH_TKN_SECRET")))
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}

func generateAccessToken(userModel *models.User) (string, error) {
	// access_token token creation
	tknTimeout, _ := strconv.Atoi(os.Getenv("SIGNED_ACCESS_TKN_TIMEOUT"))
	expireToken := time.Now().Add(time.Second * time.Duration(tknTimeout)).Unix()
	tk := &user.Token{
		ID:       fmt.Sprint(userModel.ID),
		UserName: userModel.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    os.Getenv("TKNS_ISSUER"),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNED_ACCESS_TKN_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getCtxSessionKey(c context.Context) (string, error) {
	ctxUnqKey, _ := strconv.Atoi(os.Getenv("CTX_USER_SESSION_KEY"))
	key := &user.ContextKey{
		Key: ctxUnqKey,
	}
	userContext, ok := c.Value(key).(*user.ContextData)
	if !ok {
		return "", lib.ErrInternalServerError
	}
	return userContext.SessionKey, nil
}
