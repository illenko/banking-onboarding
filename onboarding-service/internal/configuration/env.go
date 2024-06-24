package configuration

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	dir, err := os.Getwd()

	if err != nil {
		slog.Warn("Unable to get current working directory")
	}

	path := dir + "/.env"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir = filepath.Dir(dir)
		path = dir + "/.env"
	}

	if err := godotenv.Load(path); err != nil {
		slog.Warn("No .env file found in path", slog.String("path", path))
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
