package cls_vod

import (
	"net/http"
	"sync"
	"time"
)

type VodClient struct {
	config     *config
	httpClient *http.Client
}

var vodClient sync.Map

func NewVodClient() (client *VodClient, err error) {
	client = new(VodClient)
	client.httpClient = &http.Client{Timeout: time.Second * 5}
	client.config = GetConfig()
	if err := client.GetAuthToken(); err != nil {
		return nil, err
	}
	vodClient.Store(client.config.auth.userName, client)
	return
}

func GetVodClient() (client *VodClient, err error) {
	userName := GetConfig().auth.userName
	if c, ok := vodClient.Load(userName); ok {
		if client, ok := c.(*VodClient); ok && client != nil {
			return client, nil
		}
	}
	lock.Lock()
	client, err = NewVodClient()
	lock.Unlock()
	return
}
