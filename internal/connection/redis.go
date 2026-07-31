package connection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"strconv"
	"time"

	"github.com/devhdn-212/totclient_pools/internal/config"
)

// LimitCounterDB is the logical Redis DB that limittotal/limitglobal
// counters live in — kept separate from the main DB (settings cache, JWT
// blacklist, etc.) since checkout traffic can churn through a lot more keys
// than everything else combined, and the two shouldn't compete for eviction
// or get accidentally wiped together.
const LimitCounterDB = 2

var (
	RDB *redis.Client
	// RDBLimit is a dedicated, long-lived client pointed at LimitCounterDB —
	// deliberately NOT built through getClient()/SetRedis() below, since
	// those open+close a fresh connection per call and the limit-counter
	// path (two Lua script calls per bet item, every checkout) is far too
	// hot for that.
	RDBLimit *redis.Client
	ctx      = context.Background()

	// redisConf is the config InitRedis connected with — kept around so
	// getClient() can build a client for a different logical DB against the
	// same real host/port/password instead of an empty config.Redis{}.
	redisConf config.Redis
)

func InitRedis(conf config.Redis) error {
	redisConf = conf
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

	RDBLimit = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pwd,
		DB:       LimitCounterDB,
	})
	if _, err := RDBLimit.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("cannot connect to Redis (limit counter DB %d): %v", LimitCounterDB, err)
	}

	Log.Info("Connected to Redis")
	return nil
}
func RedisHealth() bool {
	if RDB == nil || RDBLimit == nil {
		Log.Fatal("Redis client not initialized. Call InitRedis() first.")
		return false
	}

	if _, err := RDB.Ping(ctx).Result(); err != nil {
		Log.Fatal("Redis health check failed: ", zap.Error(err))
		return false
	}
	if _, err := RDBLimit.Ping(ctx).Result(); err != nil {
		Log.Fatal("Redis health check failed (limit counter DB): ", zap.Error(err))
		return false
	}

	Log.Info("Redis is healthy")
	return true
}
func getClient(db int) *redis.Client {
	if db == 0 {
		if RDB == nil {
			Log.Panic("Redis client not initialized. Call InitRedis() first.")
		}
		return RDB
	}
	// temporary client untuk DB lain
	return redis.NewClient(&redis.Options{
		Addr:     redisConf.Host + ":" + redisConf.Port,
		Password: redisConf.Pass,
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
		// A transient Redis blip must not take down the whole process — this
		// runs on every cache write during normal request handling, not just
		// at startup (unlike RedisHealth). Log.Fatal here would call
		// os.Exit(1) and kill totclient_api over what's usually just a failed
		// cache write (several call sites even fire this via `go
		// connection.SetRedis(...)` and never check the error at all) — the
		// data still comes from Postgres on the next read either way, so
		// failing soft is the correct behavior.
		Log.Error("Redis Set failed", zap.String("key", key), zap.Error(err))
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
		// Same reasoning as SetRedis above — a Redis outage here should fall
		// through to a database read, not crash the process.
		Log.Error("Redis Get failed", zap.String("key", key), zap.Error(err))
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
		// Same reasoning as SetRedis above — a stale/missed cache invalidation
		// is recoverable (next read just recomputes), a crashed process isn't.
		Log.Error("Redis Delete failed", zap.String("key", key), zap.Error(err))
		return 0, err
	}
	return deleted, nil
}

// DeleteRedisByPrefix deletes every key under prefix+":*" — for cache
// namespaces where the caller only knows a prefix (e.g. an agent code), not
// every individual key underneath it (period, username, etc. segments the
// caller isn't tracking). Uses SCAN (cursor-based, non-blocking) rather than
// KEYS, since KEYS blocks the whole Redis instance while it walks a large
// keyspace.
func DeleteRedisByPrefix(prefix string, db ...int) (int64, error) {
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

	pattern := prefix + ":*"
	var cursor uint64
	var deleted int64
	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			Log.Error("Redis Scan failed", zap.String("pattern", pattern), zap.Error(err))
			return deleted, err
		}
		if len(keys) > 0 {
			n, err := client.Del(ctx, keys...).Result()
			if err != nil {
				Log.Error("Redis Delete (by prefix) failed", zap.String("pattern", pattern), zap.Error(err))
				return deleted, err
			}
			deleted += n
		}
		cursor = nextCursor
		if cursor == 0 {
			return deleted, nil
		}
	}
}

// AddRedisSet adds member to the Redis SET at key and (re)sets its TTL —
// SADD is atomic on its own (safe to call concurrently from many goroutines
// without a read-modify-write race, unlike "read the whole list, append,
// write it back"), but doesn't touch TTL by itself, so EXPIRE is applied
// right after in the same pipeline. Calling this again on an existing key
// keeps rolling its TTL forward from "now" — intentional: a period that's
// still receiving checkouts shouldn't have its index expire out from under it.
func AddRedisSet(key, member string, expire time.Duration, db ...int) error {
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

	pipe := client.TxPipeline()
	pipe.SAdd(ctx, key, member)
	pipe.Expire(ctx, key, expire)
	if _, err := pipe.Exec(ctx); err != nil {
		Log.Error("Redis SAdd failed", zap.String("key", key), zap.String("member", member), zap.Error(err))
		return err
	}
	return nil
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
