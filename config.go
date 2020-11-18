package cls_vod

import (
	"os"
	"sync"
)

var (
	locker = new(sync.Mutex)
	cfg    *config
)

type authSettings struct {
	ak       string
	sk       string
	userName string
	password string
	token    string
}

type vodSettings struct {
	endpoint    string
	domain      string
	projectName string
	projectId   string
}

type redisSettings struct {
	redisAddr     string
	redisPassword string
}

type config struct {
	auth  authSettings
	vod   vodSettings
	redis redisSettings
}

func GetConfig() *config {
	if cfg == nil {
		locker.Lock()
		InitConfig()
		locker.Unlock()
	}
	return cfg
}

func InitConfig() {
	if cfg == nil {
		_authSettings := authSettings{
			ak:       os.Getenv("AK"),
			sk:       os.Getenv("SK"),
			userName: os.Getenv("HW_USERNAME"),
			password: os.Getenv("HW_PASSWORD"),
		}
		_vodSettings := vodSettings{
			endpoint:    os.Getenv("VOD_ENDPOINT"),
			domain:      os.Getenv("VOD_DOMAIN"),
			projectName: os.Getenv("VOD_PROJECT_NAME"),
			projectId:   os.Getenv("VOD_PROJECT_ID"),
		}
		_redisSettings := redisSettings{
			redisAddr:     os.Getenv("REDIS_ADDR"),
			redisPassword: os.Getenv("REDIS_PWD"),
		}
		cfg = &config{
			auth:  _authSettings,
			vod:   _vodSettings,
			redis: _redisSettings,
		}
	}
}
