package cls_vod

import (
	"net/http"
	"sync"
	"time"
)

type vodClient struct {
	config     *config
	httpClient *http.Client
}

var clientMap sync.Map

func NewVodClient() (client *vodClient, err error) {
	client = new(vodClient)
	client.httpClient = &http.Client{Timeout: time.Second * 5}
	client.config = GetConfig()
	if err := client.GetAuthToken(); err != nil {
		return nil, err
	}
	clientMap.Store(client.config.auth.userName, client)
	return
}

func GetVodClient() (client *vodClient, err error) {
	userName := GetConfig().auth.userName
	if c, ok := clientMap.Load(userName); ok {
		if client, ok := c.(*vodClient); ok && client != nil {
			return client, nil
		}
	}
	lock.Lock()
	client, err = NewVodClient()
	lock.Unlock()
	return
}
