package lib

import (
	"os"

	"github.com/joho/godotenv"
)

func GetENV() {
	env := os.Getenv("GO_BOILER_ENV")
	var loc string
	if "" == env {
		env = "dev"
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed loading your .env file for environment (%v): \n%v", env, err)
	}
	if env == "prod" {
		loc = cwd + "/env/.env.prod"
	} else {
		loc = cwd + "/env/.env.dev"
	}

	err = godotenv.Load(loc)
	if err != nil {
		log.Fatalf("Failed loading your .env file for environment (%v): \n%v", env, err)
	}
}
