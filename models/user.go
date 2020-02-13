package models

type User struct {
	ID            uint       `json:"id"`
	FirstName     string     `json:"first_name,omitempty" db:"first_name" valid:"required~First name is required"`
	LastName      string     `json:"last_name,omitempty" db:"last_name" valid:"required~Last name is required"`
	UserName      string     `json:"username" db:"username"  valid:"required~User name is required"`
	Email         string     `json:"email" valid:"email~Enter a correct email,required~Email is required"`
	Password      string     `json:"password,omitempty" valid:"required~Password is required"`
	Avatar        NullString `json:"avatar,omitempty"`
	CountryID     NullInt64  `json:"country_id,omitempty" db:"country_id"`
	CityID        NullInt64  `json:"city_id,omitempty" db:"city_id"`
	NationalityID NullInt64  `json:"nationality_id,omitempty"  db:"nationality_id"`
	Gender        NullString `json:"gender,omitempty"`
	BirthDate     NullTime   `json:"birth_date,omitempty" db:"birth_date"`
	CreatedAt     NullTime   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     NullTime   `json:"updated_at,omitempty" db:"updated_at"`
	// EmailConfirmed
	Flags []uint8 `json:"flags,omitempty" db:"flags"`
}
