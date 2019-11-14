package database

import (
	"fmt"
	"go_boilerplate/config"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetInstance(log *logrus.Logger) (*sqlx.DB, squirrel.StatementBuilderType) {
	v, err := config.ReadConfig("db")
	if err != nil {
		log.Fatalf("Error when reading config: %v", err)
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s",
		v.Get("Username"), v.Get("Password"), v.Get("Host"), v.Get("Port"), v.Get("Database"), v.Get("ParseTime")))
	if err != nil {
		log.Fatalf("Failed connecting to the database: %v", err)
	}

	sb := squirrel.StatementBuilder.RunWith(db)
	return db, sb
}
