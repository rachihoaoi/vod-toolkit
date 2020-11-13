package cls_vod

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	REDIS_OPENCONN_ERR              = "Open redis connection set up failed, %s \n"
	REDIS_CLOSECONN_ERR             = "Close redis connection failed, %s \n"
	REDIS_DEFAULT_POOL_SIZE         = 1000
	REDIS_DEFAULT_EXPIRE_TIME       = time.Second * 60 * 5
	REDIS_DEFAULT_TOKEN_EXPIRE_TIME = time.Hour * 24 * 30
)

var (
	redisCache *redis.Client
)

func OpenRedis() *redis.Client {
	redisLocker.Lock()
	defer redisLocker.Unlock()
	if redisCache == nil {
		addr, pwd, database := cfg.redis.redisAddr, cfg.redis.redisPassword, 0
		redisCache = redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: database, PoolSize: REDIS_DEFAULT_POOL_SIZE})
		statusCmd := redisCache.Ping()
		if err := statusCmd.Err(); err != nil {
			panic(err)
		}
	}
	return redisCache
}

func CloseRedis() {
	if err := redisCache.Close(); err != nil {
		panic(err)
	}
}

func GetRedisClient() *redis.Client {
	if redisCache == nil {
		OpenRedis()
	}
	return redisCache
}

func GetWithCache(key string, result interface{}, callback func() (interface{}, time.Duration, error)) (err error) {
	var (
		res        string
		retryCount = 3
		jsonStr    []byte
	)
	for retry := 0; retry < retryCount; retry++ {
		if res, err = redisCache.Get(key).Result(); err == nil || err == redis.Nil {
			break
		}
	}
	if err != nil && err != redis.Nil {
		fmt.Println(err)
		return
	}

	if err == nil {
		if err = json.Unmarshal([]byte(res), result); err != nil {
			fmt.Println(err)
		} else {
			return
		}
	}
	midRes, expireTime, err := callback()
	if err != nil {
		fmt.Println(err)
		return
	}

	if jsonStr, err = json.Marshal(midRes); err != nil {
		fmt.Println(err)
	}

	if err = json.Unmarshal(jsonStr, result); err != nil {
		fmt.Printf("【redis.GetWithCache】catch unmarshal to result err: %s", err)
	}
	if expireTime == 0 {
		expireTime = REDIS_DEFAULT_EXPIRE_TIME
	}
	for retry := 0; retry < retryCount; retry++ {
		if _, err = redisCache.Set(key, string(jsonStr), expireTime).Result(); err == nil {
			break
		}
	}
	if err != nil {
		fmt.Printf("【redis.GetWithCache】set catch err: %s", err)
	}
	return
}
