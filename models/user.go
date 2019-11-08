package models

import (
	"database/sql"
)

type User struct {
	ID            uint           `json:"id"`
	FirstName     string         `json:"first_name" db:"first_name" valid:"required~First name is required"`
	LastName      string         `json:"last_name" db:"last_name" valid:"required~Last name is required"`
	UserName      string         `json:"username" db:"username"  valid:"required~User name is required"`
	Email         string         `json:"email" valid:"email~Enter a correct email,required~Email is required"`
	Password      string         `json:"password" valid:"required~Password is required"`
	Avatar        sql.NullString `json:"avatar"`
	CountryID     sql.NullInt64  `json:"country_id" db:"country_id"`
	CityID        sql.NullInt64  `json:"city_id" db:"city_id"`
	NationalityID sql.NullInt64  `json:"nationality_id"  db:"nationality_id"`
	Gender        sql.NullString `json:"gender"`
	BirthDate     sql.NullTime   `json:"birth_date" db:"birth_date"`
	CreatedAt     sql.NullTime   `json:"created_at" db:"created_at"`
	UpdatedAt     sql.NullTime   `json:"updated_at" db:"updated_at"`
	// EmailConfirmed
	Flags uint8 `json:"flags" db:"flags"`
}
