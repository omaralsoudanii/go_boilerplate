package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
)

func GetInstance(log *logrus.Logger) *sqlx.DB {
	log.Infoln("Connecting to MySQL...")
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s&columnsWithAlias=true",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARSE_TIME"))

	db, err := sqlx.Connect("mysql", addr)

	if err != nil {
		log.Fatalf("Failed connecting to the database: %v", err)
	}

	log.Infoln("MySQL started at: " + addr)
	return db
}
