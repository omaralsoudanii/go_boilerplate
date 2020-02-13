package repository

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_lib "go_boilerplate/lib"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"time"
)

type userRepository struct {
	Conn      *sqlx.DB
	RedisConn *redis.Client
}

func NewUserRepository(Conn *sqlx.DB, r *redis.Client) user.Repository {

	return &userRepository{Conn, r}
}

func (repo *userRepository) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	tx, err := repo.Conn.BeginTxx(ctx, nil)
	if err != nil {
		_ = tx.Rollback()
		return nil, _lib.ErrInternalServerError
	}

	result, err := tx.ExecContext(ctx, "INSERT INTO user ( first_name, last_name , username , email , password , created_at) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.UserName, user.Email, user.Password, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return nil, _lib.ErrInternalServerError
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, _lib.ErrInternalServerError
	}

	u, err := repo.GetByID(ctx, uint(id))
	if err != nil {
		return nil, _lib.ErrInternalServerError
	}

	err = tx.Commit()
	if err != nil {
		return nil, _lib.ErrInternalServerError
	}
	return u, nil
}

func (repo *userRepository) GetByName(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT * from user where username = ?`
	userModel := &models.User{}
	err := repo.Conn.GetContext(ctx, userModel, query, username)
	if err != nil {
		return nil, _lib.ErrNotFound
	}
	return userModel, nil
}

func (repo *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	query := `SELECT * from user where id = ?`
	userModel := &models.User{}
	err := repo.Conn.GetContext(ctx, userModel, query, id)
	if err != nil {
		return nil, _lib.ErrNotFound
	}
	return userModel, nil
}

func (repo *userRepository) StoreSession(user *models.User, key string, token string) error {
	userMap := map[string]interface{}{
		"id":           user.ID,
		"username":     user.UserName,
		"email":        user.Email,
		"avatar":       user.Avatar,
		"birth_date":   user.BirthDate,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"gender":       user.Gender,
		"refreshToken": token,
	}

	// create/set user in redis
	err := repo.RedisConn.HMSet(key, userMap).Err()
	if err != nil {
		return err
	}

	return nil
}
func (repo *userRepository) DeleteSession(key string) error {
	err := repo.RedisConn.Del(key).Err()
	return err
}

func (repo *userRepository) GetUser(key string) (map[string]string, error) {
	data := make(map[string]string)
	data, err := repo.RedisConn.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}
