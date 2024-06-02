package util

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName         string `envconfig:"DB_NAME"`
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	LevelLoger     string
	ServerPort     string
	NatsClussterID string
	NatsClientID   string
	NatsURL        string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	c.DBName = os.Getenv("DB_NAME")
	c.DBPassword = os.Getenv("DB_PASSWORD")
	c.DBUser = os.Getenv("DB_USER")
	c.DBPort = os.Getenv("DB_PORT")
	c.DBHost = os.Getenv("DB_HOST")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.LevelLoger = os.Getenv("LELEV_LOGGER")
	c.NatsClussterID = os.Getenv("NATS_CLUSSTER_ID")
	c.NatsClientID = os.Getenv("NATS_CLIENT_ID")
	c.NatsURL = os.Getenv("NATS_URL")
	return nil
}
