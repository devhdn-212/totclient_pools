package connection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"strconv"
	"time"

	"github.com/devhdn-212/totagen_api/internal/config"
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
		Log.Info("Redis env variables missing")
		return fmt.Errorf("redis env variables missing")
	}
	dbNum, err := strconv.Atoi(dbStr)
	if err != nil {
		return fmt.Errorf("invalid DB_REDIS_NAME: %v", zap.Error(err))
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pwd,
		DB:       dbNum,
	})

	if _, err := RDB.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("cannot connect to Redis: %v", err)
	}
	Log.Info("Connected to Redis")
	return nil
}
func RedisHealth() bool {
	if RDB == nil {
		Log.Fatal("Redis client not initialized. Call InitRedis() first.")
		return false
	}

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		Log.Fatal("Redis health check failed: ", zap.Error(err))
		return false
	}

	Log.Info("Redis is healthy")
	return true
}
func getClient(db int) *redis.Client {
	var conf config.Redis
	if db == 0 {
		if RDB == nil {
			Log.Panic("Redis client not initialized. Call InitRedis() first.")
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
		Log.Fatal("Redis Set failed : ", zap.Error(err))
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
		Log.Fatal("Redis Get failed : ", zap.Error(err))
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
		Log.Fatal("Redis Delete failed : ", zap.Error(err))
		return 0, err
	}
	return deleted, nil
}

func BlacklistJWT(jti string, ttl time.Duration) error {
	if jti == "" {
		return fmt.Errorf("empty jti")
	}
	return SetRedis("master:jwt:blacklist:"+jti, "1", ttl)
}

func IsJWTBlacklisted(jti string) (bool, error) {
	if jti == "" {
		return false, fmt.Errorf("empty jti")
	}
	_, found, err := GetRedis("master:jwt:blacklist:" + jti)
	return found, err
}
