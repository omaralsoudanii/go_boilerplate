package database

import (
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetInstance(log *logrus.Logger) (*sqlx.DB, squirrel.StatementBuilderType) {
	log.Infoln("Connecting to MySQL...")
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARSE_TIME"))

	db, err := sqlx.Connect("mysql", addr)

	if err != nil {
		log.Fatalf("Failed connecting to the database: %v", err)
	}

	sb := squirrel.StatementBuilder.RunWith(db)
	log.Infoln("MySQL started at: " + addr)
	return db, sb
}
