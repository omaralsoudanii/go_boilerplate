package database

import (
	"fmt"
	"go_boilerplate/lib"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetInstance(log *logrus.Logger) (*sqlx.DB, squirrel.StatementBuilderType) {
	lib.GetENV()
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARSE_TIME")))

	if err != nil {
		log.Fatalf("Failed connecting to the database: %v", err)
	}

	sb := squirrel.StatementBuilder.RunWith(db)
	return db, sb
}
