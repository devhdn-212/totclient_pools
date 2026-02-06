package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"gofibergocu/internal/config"
	"strconv"
	"time"
)

var (
	RDB *redis.Client
	ctx = context.Background()
)

func InitRedis(conf config.Redis) error {
	host := conf.Host
	port := conf.Port
	pwd := conf.Pass
	dbStr := conf.Name
	if host == "" || port == "" || dbStr == "" {
		return fmt.Errorf("redis env variables missing")
	}
	dbNum, err := strconv.Atoi(dbStr)
	if err != nil {
		return fmt.Errorf("invalid DB_REDIS_NAME: %v", err)
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pwd,
		DB:       dbNum,
	})

	if _, err := RDB.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("cannot connect to Redis: %v", err)
	}
	log.Info("Connected to Redis")
	return nil
}
func RedisHealth() bool {
	if RDB == nil {
		log.Error("Redis client not initialized. Call InitRedis() first.")
		return false
	}

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		fmt.Errorf("Redis health check failed: %v", err)
		return false
	}

	log.Info("Redis is healthy")
	return true
}
func getClient(db int) *redis.Client {
	var conf config.Redis
	if db == 0 {
		if RDB == nil {
			panic("Redis client not initialized. Call InitRedis() first.")
		}
		return RDB
	}
	// temporary client untuk DB lain
	dbHost := conf.Host + ":" + conf.Port
	dbPass := conf.Pass
	return redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: dbPass,
		DB:       db,
	})
}
func SetRedis(key string, data interface{}, expire time.Duration, db ...int) error {
	targetDB := 0
	if len(db) > 0 {
		targetDB = db[0]
	}

	client := getClient(targetDB)
	defer func() {
		if targetDB != 0 { // jangan close global client
			client.Close()
		}
	}()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = client.Set(ctx, key, jsonData, expire).Err()
	if err != nil {
		fmt.Sprintf("Redis Set failed : %v", err)
		return err
	}
	return nil
}

func GetRedis(key string, db ...int) (string, bool, error) {
	targetDB := 0
	if len(db) > 0 {
		targetDB = db[0]
	}

	client := getClient(targetDB)
	defer func() {
		if targetDB != 0 {
			client.Close()
		}
	}()

	result, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		fmt.Sprintf("Redis Get failed : %v", err)
		return "", false, err
	}
	return result, true, nil
}

func DeleteRedis(key string, db ...int) (int64, error) {
	targetDB := 0
	if len(db) > 0 {
		targetDB = db[0]
	}

	client := getClient(targetDB)
	defer func() {
		if targetDB != 0 {
			client.Close()
		}
	}()

	deleted, err := client.Del(ctx, key).Result()
	if err != nil {
		fmt.Sprintf("Redis Delete failed : %v", err)
		return 0, err
	}
	return deleted, nil
}
