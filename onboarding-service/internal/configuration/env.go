package configuration

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func LoadEnv() {

	dir, err := os.Getwd()
	path := dir + "/.env"

	if err != nil {
		log.Println("Not able to get current working director")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir = filepath.Dir(dir)
		path = dir + "/.env"
	}

	if err := godotenv.Load(path); err != nil {
		log.Println("No .env file found in path " + path)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
