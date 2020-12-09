package cls_vod

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AlgorithmType string

var Algorithm = struct {
	AlgorithmA AlgorithmType
	AlgorithmB AlgorithmType
	AlgorithmC AlgorithmType
	AlgorithmD AlgorithmType
}{AlgorithmA: "Algorithm_A", AlgorithmB: "Algorithm_B", AlgorithmC: "Algorithm_C", AlgorithmD: "Algorithm_D"}

type referEncodingAlgorithm interface {
	setOriginalUrl(u string)
	GetAuthorizedUrl() string
}

type (
	algorithmA struct {
		timestamp  int64
		rand       string
		uid        int
		privateKey string
		referUrl   string
	}
	algorithmB struct {
		date       string
		fileName   string
		privateKey string
		referUrl   string
	}
	algorithmC struct {
		timeHex    string
		fileName   string
		privateKey string
		referUrl   string
	}
	algorithmD struct {
		path       string
		timeStamp  string
		privateKey string
	}
)

// ------------------------------algorithmA------------------------------
// https://support.huaweicloud.com/usermanual-vod/vod010014.html#section4
func (a *algorithmA) setOriginalUrl(u string) {
	a.referUrl = u
}

func (a *algorithmA) GetAuthorizedUrl() string {
	// url_template = URL?auth_key={timestamp}-{rand}-{uid}-{auth_key}
	// auth_key=MD5(/asset/{assetId}/{file_name}-{timestamp}-{rand}-{uid}-{private_key})
	assetPart := strings.SplitAfterN(a.referUrl, "/", 4)
	rawStr := fmt.Sprintf("%s-%d-%s-%d-%s", assetPart[len(assetPart)-1], a.timestamp, a.rand, a.uid, a.privateKey)
	authKey := fmt.Sprintf("%s", md5.Sum([]byte(rawStr)))
	template := "auth_key=%d-%s-%d-%x"
	return a.referUrl + "?" + fmt.Sprintf(template, a.timestamp, a.rand, a.uid, authKey)
}

// ------------------------------algorithmB------------------------------
// https://support.huaweicloud.com/usermanual-vod/vod010014.html#section5
func (b *algorithmB) setOriginalUrl(u string) {
	b.referUrl = u
}

func (b *algorithmB) GetAuthorizedUrl() string {
	// md5sum = md5({private_key}{date_yyyyMMddHHmm}/asset/{asset_id}/{file_name})
	// url_template = https://{cdn_domain}/{date_YYYYmmddHHMM}/{md5sum}/asset/{asset_id}/{file_name}
	assetPart := strings.Split(b.referUrl, "/")
	cdnName, assetName, fileName := assetPart[2], assetPart[4], strings.Join(assetPart[5:], "/")
	b.fileName = fileName
	rawStr := fmt.Sprintf("%s%s/asset/%s/%s", b.privateKey, b.date, assetName, b.fileName)
	md5Sum := fmt.Sprintf("%s", md5.Sum([]byte(rawStr)))
	return fmt.Sprintf("https://%s/%s/%x/asset/%s/%s", cdnName, b.date, md5Sum, assetName, b.fileName)
}

// ------------------------------algorithmC------------------------------
// https://support.huaweicloud.com/usermanual-vod/vod010014.html#section6
func (c *algorithmC) setOriginalUrl(u string) {
	c.referUrl = u
}

func (c *algorithmC) GetAuthorizedUrl() string {
	// md5sum = md5({private_key}{date_yyyyMMddHHmm}/asset/{asset_id}/{file_name})
	// url_template = https://{cdn_domain}/{md5hash}/{time_hex}/asset/{asset_id}/{file_name}
	assetPart := strings.Split(c.referUrl, "/")
	cdnName, assetName, fileName := assetPart[2], assetPart[4], strings.Join(assetPart[5:], "/")
	c.fileName = fileName
	rawStr := fmt.Sprintf("%s/asset/%s/%s%s", c.privateKey, assetName, c.fileName, c.timeHex)
	md5Sum := fmt.Sprintf("%s", md5.Sum([]byte(rawStr)))
	return fmt.Sprintf("https://%s/%x/%s/asset/%s/%s", cdnName, md5Sum, c.timeHex, assetName, c.fileName)
}

// ------------------------------algorithmD------------------------------
// https://support.huaweicloud.com/usermanual-vod/vod010014.html#section7
func (d *algorithmD) setOriginalUrl(u string) {

}

func (d *algorithmD) GetAuthorizedUrl() string {
	return ""
}

func (c *vodClient) GetAuthorizedUrl() string {
	return c.config.refer.algorithm.GetAuthorizedUrl()
}

func (c *vodClient) SetAlgorithm(algorithmType AlgorithmType) *vodClient {
	if c.config == nil {
		return c
	}
	c.config.refer.algorithm = generateAlgorithm(algorithmType, c.config.refer.privateKey)
	return c
}

func (c *vodClient) SetOriginalUrl(url string) {
	if c.config == nil {
		return
	}
	if c.config.refer.algorithm == nil {
		fmt.Println("invalid encoding algorithm")
	}
	c.config.refer.originUrl = url
	c.config.refer.algorithm.setOriginalUrl(url)
}

func generateAlgorithm(algorithmType AlgorithmType, privateKey string) referEncodingAlgorithm {
	switch algorithmType {
	case Algorithm.AlgorithmA:
		return newAlgorithmA(privateKey)
	case Algorithm.AlgorithmB:
		return newAlgorithmB(privateKey)
	case Algorithm.AlgorithmC:
		return newAlgorithmC(privateKey)
	case Algorithm.AlgorithmD:
		return newAlgorithmD(privateKey)
	default:
		return newAlgorithmA(privateKey)
	}
}

func newAlgorithmA(privateKey string) *algorithmA {
	return &algorithmA{
		timestamp:  time.Now().Unix(),
		rand:       strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		uid:        0,
		privateKey: privateKey,
	}
}

func newAlgorithmB(privateKey string) *algorithmB {
	return &algorithmB{
		date:       time.Now().Format("200601021504"),
		privateKey: privateKey,
	}
}

func newAlgorithmC(privateKey string) *algorithmC {
	return &algorithmC{
		timeHex:    fmt.Sprintf("%X", time.Now().Unix()),
		privateKey: privateKey,
	}
}

func newAlgorithmD(privateKey string) *algorithmD {
	return &algorithmD{
		path:       "",
		timeStamp:  "",
		privateKey: privateKey,
	}
}
