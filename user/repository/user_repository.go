package repository

import (
	"context"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"time"

	_lib "go_boilerplate/lib"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

var log = _lib.GetLogger()

type userRepository struct {
	DbConn    *sqlx.DB
	RedisConn *redis.Client
}

func NewUserRepository(db *sqlx.DB, r *redis.Client) user.Repository {

	return &userRepository{db, r}
}

func (repo *userRepository) Insert(ctx context.Context, user *models.User) error {

	tx := repo.DbConn.MustBegin()
	tx.MustExecContext(ctx, "INSERT INTO user ( first_name, last_name , username , email , password , created_at) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.UserName, user.Email, user.Password, time.Now())
	saveErr := tx.Commit()
	if saveErr != nil {
		log.Error(saveErr)
		return _lib.ErrInternalServerError
	}
	return nil
}

func (repo *userRepository) FetchByName(ctx context.Context, userName string) (*models.User, error) {
	query := `SELECT * from user where username = $1`
	userModel := &models.User{}
	err := repo.DbConn.GetContext(ctx, userModel, query, userName)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return userModel, nil
}

func (repo *userRepository) FetchByID(ctx context.Context, ID uint) (*models.User, error) {
	query := `SELECT * from user where id = $1`
	userModel := &models.User{}
	err := repo.DbConn.GetContext(ctx, userModel, query, ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return userModel, nil
}

func (repo *userRepository) StoreSession(ctx context.Context, user *models.User, token string) error {
	userMap := map[string]interface{}{
		"id":           user.ID,
		"username":     user.UserName,
		"email":        user.Email,
		"avatar":       user.Avatar.String,
		"birth_date":   user.BirthDate.Time,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"gender":       user.Gender.String,
		"refreshToken": token,
	}
	redisKey := "user:" + user.UserName

	// create/set user in redis
	err := repo.RedisConn.HMSet(redisKey, userMap).Err()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
func (repo *userRepository) DeleteSession(key string) error {
	err := repo.RedisConn.Del(key).Err()
	if err != nil {
		log.Error(err)
	}
	return err
}
func (repo *userRepository) GetUser(userName string) (map[string]string, error) {
	key := "user:" + userName
	data := make(map[string]string)
	data, err := repo.RedisConn.HGetAll(key).Result()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return data, nil
}
