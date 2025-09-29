package redis

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"fmt"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()
var redisAddr string
var redisPwd string
var redisDb int

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	redisAddr = os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set")
	}

	redisPwd = os.Getenv("REDIS_PWD")
	redisDbStr := os.Getenv("REDIS_DB")
	if redisDbStr == "" {
		log.Fatal("REDIS_DB is not set")
	}

	redisDb, err = strconv.Atoi(redisDbStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}
}

func Connect() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPwd,
		DB:       redisDb,
	})

	// Ping Redis to check connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis Ping Response:", pong)

	return nil
}

func LPush(key string, data interface{}) error {
	if rdb == nil {
		return fmt.Errorf("no redis instance has been created")
	}

	var jsonData []byte
	var err error

	// Check if data is already a JSON string (byte slice)
	switch v := data.(type) {
	case []byte:
		// Data is already JSON (byte slice), use as is
		jsonData = v
	default:
		// Data is not a JSON string, so we need to marshal it
		jsonData, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("cannot marshal data into JSON: %v", err)
		}
	}

	// Push the (potentially marshaled) JSON data to the Redis list
	_, err = rdb.LPush(ctx, key, jsonData).Result()
	if err != nil {
		log.Printf("Could not LPUSH JSON to key %s: %v", key, err)
		return err
	}

	return nil
}

func LRange(key string, rangeLen int64) ([]interface{}, error) {
	cmd := rdb.LRange(ctx, key, 0, rangeLen)
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	jsonList := make([]interface{}, len(result))

	for i, item := range result {
		var jsonData interface{}
		err := json.Unmarshal([]byte(item), &jsonData)
		if err != nil {
			return nil, err
		}

		jsonList[i] = jsonData
	}

	return jsonList, nil
}

func BLPop(key string) (interface{}, error) {
	result, err := rdb.BLPop(ctx, 500*time.Millisecond, key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func LLen(key string) (int64, error) {
	queLen, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return queLen, nil
}

func RPush(key string, item interface{}) error {
	return rdb.RPush(ctx, key, item).Err()
}
