package repository

import (
	"context"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/go-redis/redis"
)

type userRepository struct {
	sb        squirrel.StatementBuilderType
	db        *sqlx.DB
	RedisConn *redis.Client
}

func NewUserRepository(sb squirrel.StatementBuilderType, db *sqlx.DB, r *redis.Client) user.Repository {

	return &userRepository{sb, db, r}
}

func (repo *userRepository) Insert(ctx context.Context, user *models.User) error {

	q := repo.sb.Insert("user").
		Columns("first_name", "last_name", "username", "email", "password", "created_at").
		Values(user.FirstName, user.LastName, user.UserName, user.Email, user.Password, time.Now())
	_, err := q.ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) GetByName(ctx context.Context, username string) (*models.User, error) {
	q, args, err := repo.sb.Select("*").
		From("user").
		Where("username = ?", username).
		ToSql()

	if err != nil {
		return nil, err
	}
	userModel := &models.User{}
	err = repo.db.GetContext(ctx, userModel, q, args...)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (repo *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	q, args, err := repo.sb.Select("*").
		From("user").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	userModel := &models.User{}
	err = repo.db.GetContext(ctx, userModel, q, args)
	if err != nil {
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
		return err
	}

	return nil
}
func (repo *userRepository) DeleteSession(key string) error {
	err := repo.RedisConn.Del(key).Err()
	return err
}

func (repo *userRepository) GetUser(username string) (map[string]string, error) {
	key := "user:" + username
	data := make(map[string]string)
	data, err := repo.RedisConn.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}
