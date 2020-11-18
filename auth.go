package cls_vod

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rachihoaoi/vod_toolkit/utils"
)

const (
	TOKEN_EXPIRE_TIME = time.Hour * 23
)

type (
	GetAuthTokensPayload struct {
		Auth auth `json:"auth"`
	}

	auth struct {
		Identity *identity `json:"identity"`
		Scope    *scope    `json:"scope,omitempty"`
	}

	identity struct {
		Methods  []string  `json:"methods,omitempty" desc:"methods"`
		Password *password `json:"password,omitempty" desc:"password"`
	}

	password struct {
		User *user `json:"user,omitempty"`
	}

	user struct {
		Name     string  `json:"name,omitempty"`
		Password string  `json:"password,omitempty"`
		Domain   *domain `json:"domain,omitempty"`
	}

	domain struct {
		Name string `json:"name" desc:"name,omitempty"`
	}

	scope struct {
		Project *project `json:"project"`
	}

	project struct {
		Name string `json:"name" desc:"name"`
	}

	BuildConfig struct {
		Method       []string
		UserName     string
		UserPassword string
		DomainName   string
		ProjectName  string
	}

	GetSecurityTokensRequest struct {
		Auth auth `json:"auth"`
	}

	GetSecurityTokensResp struct {
		Credential *credential `json:"credential"`
	}

	credential struct {
		ExpiresAt     string `json:"expires_at"`
		Access        string `json:"access"`
		Secret        string `json:"secret"`
		SecurityToken string `json:"securitytoken"`
	}

	JSClientConfig struct {
		AccessKeyId     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		SecurityToken   string `json:"security_token"`
		ProjectId       string `json:"project_id"`
		VodServer       string `json:"vod_server"`
		VodPort         string `json:"vod_port"`
	}
)

func BuildGetAuthTokensPayload(config *BuildConfig) *GetAuthTokensPayload {
	return &GetAuthTokensPayload{
		Auth: auth{
			Identity: &identity{
				Methods: config.Method,
				Password: &password{
					User: &user{
						Name:     config.UserName,
						Password: config.UserPassword,
						Domain: &domain{
							Name: config.DomainName,
						},
					}},
			},
			Scope: &scope{
				Project: &project{
					Name: config.ProjectName,
				},
			},
		},
	}
}

func (c *vodClient) RefreshToken() error {
	redisClient := GetRedisClient()
	key := fmt.Sprintf("%s_%s_%s", c.config.vod.projectName, c.config.vod.domain, c.config.auth.userName)
	requestBody := BuildGetAuthTokensPayload(&BuildConfig{
		Method:       []string{"password"},
		UserName:     c.config.auth.userName,
		UserPassword: c.config.auth.password,
		DomainName:   c.config.vod.domain,
		ProjectName:  c.config.vod.projectName,
	})
	_, header, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client:  c.httpClient,
		Url:     AUTH_TOKEN_URL,
		Method:  http.MethodPost,
		Header:  nil,
		Payload: requestBody,
	})
	if err != nil {
		return err
	}
	c.config.auth.token = header.Get(HUAWEI_TOKEN_NAME)
	if c.config.auth.token != "" {
		redisClient.Set(key, c.config.auth.token, TOKEN_EXPIRE_TIME)
		return nil
	}
	return errors.New("empty Token")
}

func (c *vodClient) GetAuthToken() (err error) {
	redisClient := GetRedisClient()

	key := fmt.Sprintf("%s_%s_%s", c.config.vod.projectName, c.config.vod.domain, c.config.auth.userName)
	if token, err := redisClient.Get(key).Result(); err == nil && token != "" {
		fmt.Println("[Get Token From Redis]: " + token)
		c.config.auth.token = token
		return nil
	}
	return c.RefreshToken()
}

func (c *vodClient) GetSecurityTokens() (resp *GetSecurityTokensResp, err error) {
	payload := &GetSecurityTokensRequest{
		Auth: auth{
			Identity: &identity{
				Methods: []string{"token"},
			},
		},
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: c.httpClient,
		Url:    GET_TEMP_SECU_TOKEN_URL,
		Method: http.MethodPost,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: c.config.auth.token,
		},
		Payload: payload,
	})
	if err != nil {
		return resp, err
	}
	if err = json.Unmarshal(b, &resp); err != nil {
		return resp, err
	}
	return
}

func (c *vodClient) GetJSClientConfig() (resp *JSClientConfig, err error) {
	var st = new(GetSecurityTokensResp)
	if st, err = c.GetSecurityTokens(); err != nil {
		return resp, err
	}
	resp = &JSClientConfig{
		AccessKeyId:     st.Credential.Access,
		SecretAccessKey: st.Credential.Secret,
		SecurityToken:   st.Credential.SecurityToken,
		ProjectId:       c.config.vod.projectId,
		VodServer:       VOD_HOST,
		VodPort:         "",
	}
	return resp, err
}
