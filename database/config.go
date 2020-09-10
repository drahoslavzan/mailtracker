package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	DbName   string
	EmailCol string
}

func init() {
	godotenv.Load("./database/.env")
}

func GetConfig() (cfg Config) {
	cfg = Config{
		MongoURI: getEnv("MONGO"),
		DbName:   getEnv("DB_NAME"),
		EmailCol: getEnv("DB_COL_EMAIL"),
	}

	return
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		panic(fmt.Errorf("missing env value for '%s'", key))
	}
	return val
}
