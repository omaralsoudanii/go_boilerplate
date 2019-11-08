package models

import (
	"database/sql"

	pg "github.com/lib/pq"
)

type User struct {
	ID             uint           `json:"id"`
	FirstName      string         `json:"first_name" db:"first_name" valid:"required~First name is required"`
	LastName       string         `json:"last_name" db:"last_name" valid:"required~Last name is required"`
	UserName       string         `json:"user_name" db:"user_name"  valid:"required~User name is required"`
	Email          string         `json:"email" valid:"email~Enter a correct email,required~Email is required"`
	Password       string         `json:"password" valid:"required~Password is required"`
	EmailConfirmed bool           `json:"email_confirmed" db:"email_confirmed"`
	Avatar         sql.NullString `json:"avatar"`
	CountryId      sql.NullInt64  `json:"country_id" db:"country_id"`
	CityId         sql.NullInt64  `json:"city_id" db:"city_id"`
	NationalityId  sql.NullInt64  `json:"nationality_id"  db:"nationality_id"`
	Gender         sql.NullString `json:"gender"`
	BirthDate      pg.NullTime    `json:"birth_date" db:"birth_date"`
	CreatedAt      pg.NullTime    `json:"created_at" db:"created_at"`
	UpdatedAt      pg.NullTime    `json:"updated_at" db:"updated_at"`
}
