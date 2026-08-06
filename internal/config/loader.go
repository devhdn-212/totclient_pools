package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnvInt(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	expInt, _ := strconv.Atoi(os.Getenv("JWT_EXP"))

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host:   os.Getenv("DB_HOST"),
			Port:   os.Getenv("DB_PORT"),
			User:   os.Getenv("DB_USER"),
			Pass:   os.Getenv("DB_PASS"),
			Name:   os.Getenv("DB_NAME"),
			Schema: os.Getenv("DB_SCHEMA"),
			Tz:     os.Getenv("DB_TIMEZONE"),
			// Defaults sized for THIS worker's actual concurrency profile
			// (one Kafka message processed at a time per replica), not
			// copy-pasted from totclient_api's HTTP-server sizing — see the
			// doc comment on Database.MaxConns.
			MaxConns: int32(getEnvInt("DB_MAX_CONNS", 10)),
			MinConns: int32(getEnvInt("DB_MIN_CONNS", 2)),
		},
		Jwt: Jwt{
			Key:      os.Getenv("JWT_KEY"),
			Exp:      expInt,
			Issuer:   os.Getenv("JWT_ISSUER"),
			Audience: os.Getenv("JWT_AUDIENCE"),
		},
		Redis: Redis{
			Host: os.Getenv("DB_REDIS_HOST"),
			Port: os.Getenv("DB_REDIS_PORT"),
			Pass: os.Getenv("DB_REDIS_PASSWORD"),
			Name: os.Getenv("DB_REDIS_NAME"),
		},
		Telegram: Telegram{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		},
		BalanceAPI: BalanceAPI{
			APIKey: os.Getenv("BALANCE_API_KEY"),
		},
		Kafka: Kafka{
			Brokers: os.Getenv("KAFKA_BROKERS"),
			Topic:   os.Getenv("KAFKA_TOPIC"),
			GroupID: os.Getenv("KAFKA_GROUP_ID"),
		},
	}
}
