package cls_vod

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/rachihoaoi/vod_toolkit/utils"
)

type (
	VideoInfo struct {
		Client    *vodClient
		Title     string
		Desc      string
		Type      string
		VideoName string
		VideoType MediaType

		assetId       string
		authorizedUrl string
		target        target
		uploadResult  InitiateMultipartUploadResult
	}

	CreateAssetRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryId  int64  `json:"category_id"`
		VideoName   string `json:"video_name"`
		VideoType   string `json:"video_type"`
		VideoMd5    string `json:"video_md5"`
		AutoPublish int64  `json:"auto_publish"`
	}

	CreateAssetResponse struct {
		AssetId            string   `json:"asset_id"`
		VideoUploadUrl     string   `json:"video_upload_url"`
		CoverUploadUrl     string   `json:"cover_upload_url"`
		SubtitleUploadUrls []string `json:"subtitle_upload_urls"`
		Target             target   `json:"target"`
	}

	target struct {
		Bucket   string `json:"bucket"`
		Location string `json:"location"`
		Object   string `json:"object"`
	}

	ErrorInfo struct {
		ErrorCode string `json:"error_code" desc:"错误码"`
		ErrorMsg  string `json:"error_msg" desc:"错误描述"`
	}

	UploadedRequest struct {
		AssetId string `json:"asset_id"`
		Status  string `json:"status"`
	}

	UploadedResp struct {
		ErrorInfo
		AssetId string `json:"asset_id"`
	}

	PublishRequest struct {
		AssetId []string `json:"asset_id"`
	}

	PublishResp struct {
		AssetInfoArray []AssetInfo `json:"asset_info_array"`
	}

	AssetInfo struct {
		AssetId       string    `json:"asset_id"`
		Status        string    `json:"status"`
		Description   string    `json:"description"`
		BaseInfo      *BaseInfo `json:"base_info"`
		PlayInfoArray *PlayInfo `json:"play_info_array"`
	}

	BaseInfo struct {
		Title        string `json:"title"`
		VideoName    string `json:"video_name"`
		Description  string `json:"description"`
		CategoryId   int64  `json:"category_id"`
		CategoryName string `json:"category_name"`
		CreateTime   string `json:"create_time"`
		LastModified string `json:"last_modified"`
		VideoType    string `json:"video_type"`
		VideoUrl     string `json:"video_url"`
	}

	TranscodeInfo struct {
	}

	ThumbnailInfo struct {
	}

	ReviewInfo struct {
	}

	PlayInfo struct {
		PlayType  string     `json:"play_type"`
		Url       string     `json:"url"`
		Encrypted int        `json:"encrypted"`
		MeteData  []MetaData `json:"meta_data"`
	}

	MetaData struct {
		Duration  float64 `json:"duration"`
		VideoSize float64 `json:"video_size"`
		Width     float64 `json:"width"`
		Hight     float64 `json:"hight"`
	}

	GetAssetDetailRequest struct {
		AssetId    string `json:"asset_id"`
		Categories string `json:"categories"`
	}

	GetAssetDetailResponse struct {
		AssetId  string    `json:"asset_id"`
		BaseInfo *BaseInfo `json:"base_info"`
		// TranscodeInfo *TranscodeInfo `json:"transcode_info"`
		// ThumbnailInfo *ThumbnailInfo `json:"thumbnail_info"`
		// ReviewInfo    *ReviewInfo    `json:"review_info"`
	}

	InitAssetAuthorityRequest struct {
		HttpVerb    string `json:"http_verb"`
		ContentType string `json:"content_type"`
		Bucket      string `json:"bucket"`
		ObjectKey   string `json:"object_key"`
	}

	InitAssetAuthorityResponse struct {
		SignStr string `json:"sign_str"`
	}

	GetAssetAuthorityRequest struct {
		HttpVerb    string `json:"http_verb"`
		ContentType string `json:"content_type"`
		Bucket      string `json:"bucket"`
		ObjectKey   string `json:"object_key"`
		ContentMd5  string `json:"content_md5"`
		UploadId    string `json:"upload_id"`
		PartNumber  int64  `json:"part_number"`
	}

	GetAssetAuthorityResponse struct {
		SignStr string `json:"sign_str"`
	}

	InitiateMultipartUploadResult struct {
		XMLName     xml.Name `xml:"InitiateMultipartUploadResult"`
		Version     string   `xml:"version,attr"`
		Bucket      string   `xml:"Bucket"`
		Key         string   `xml:"Key"`
		UploadId    string   `xml:"UploadId"`
		Description string   `xml:",innerxml"`
	}
)

var Catagory = map[string]int64{
	"test": VOD_CATAGORY_TEST,
}

func (i *VideoInfo) CreateAsset() (resp *CreateAssetResponse, err error) {
	resp = new(CreateAssetResponse)
	request := &CreateAssetRequest{
		Title:       i.Title,
		Description: i.Desc,
		VideoName:   i.VideoName,
		VideoType:   string(i.VideoType),
		AutoPublish: 0,
	}
	if catalogId, ok := Catagory[i.Type]; ok {
		request.CategoryId = catalogId
	} else {
		request.CategoryId = VOD_CATAGORY_OTHERS
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: i.Client.httpClient,
		Url:    fmt.Sprintf(CREATE_ASSET_URL, i.Client.config.vod.projectId),
		Method: http.MethodPost,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: i.Client.config.auth.token,
			"Content-Type":         "application/json;charset=UTF-8",
		},
		Payload: request,
	})
	if err != nil {
		return
	}
	err = json.Unmarshal(b, resp)
	i.target = resp.Target
	return
}

func (i *VideoInfo) SetAssetId(assetId string) {
	i.assetId = assetId
}

func (i *VideoInfo) ConfirmUpload(assetId string, status string) (resp *UploadedResp, err error) {
	resp = new(UploadedResp)
	payload := &UploadedRequest{
		AssetId: assetId,
		Status:  status,
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: i.Client.httpClient,
		Url:    fmt.Sprintf(ASSET_UPLOAD_CHECK_URL, i.Client.config.vod.projectId),
		Method: http.MethodPost,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: i.Client.config.auth.token,
		},
		Payload: payload,
		TagName: TAG_JSON,
	})
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	if resp.ErrorCode != "" {
		return nil, errors.New(resp.ErrorCode + ":" + resp.ErrorMsg)
	}
	return resp, nil
}

func (i *VideoInfo) Publish() (resp PublishResp, err error) {
	if i.assetId == "" {
		return resp, errors.New("请先创建上传凭证或指定assetId")
	}
	payload := &PublishRequest{
		AssetId: []string{i.assetId},
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: i.Client.httpClient,
		Url:    fmt.Sprintf(PUBLISH_VIDEO_URL, i.Client.config.vod.projectId),
		Method: http.MethodPost,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: i.Client.config.auth.token,
		},
		Payload: payload,
		TagName: TAG_JSON,
	})
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	return
}

func (i *VideoInfo) GetAssetDetail() (resp *GetAssetDetailResponse, err error) {
	resp = new(GetAssetDetailResponse)
	if i.assetId == "" {
		return resp, errors.New("empty asset id")
	}
	payload := &GetAssetDetailRequest{
		AssetId:    i.assetId,
		Categories: "base_info",
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: i.Client.httpClient,
		Url:    fmt.Sprintf(GET_ASSET_DETAIL_URL, i.Client.config.vod.projectId),
		Method: http.MethodGet,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: i.Client.config.auth.token,
		},
		Payload: payload,
		TagName: TAG_JSON,
	})
	if err != nil {
		return resp, err
	}
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	return
}

func (i *VideoInfo) InitAssetAuthority(bucket, objectKey string) (resp *InitAssetAuthorityResponse, err error) {
	resp = new(InitAssetAuthorityResponse)
	queryPayload := &InitAssetAuthorityRequest{
		HttpVerb:    http.MethodPost,
		ContentType: "video/mp4",
		Bucket:      i.target.Bucket,
		ObjectKey:   i.target.Object,
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: i.Client.httpClient,
		Url:    fmt.Sprintf(GET_ASSET_AUTH_URL, i.Client.config.vod.projectId),
		Method: http.MethodGet,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: i.Client.config.auth.token,
		},
		Payload: queryPayload,
		TagName: TAG_JSON,
	})
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &resp); err == nil {
		i.authorizedUrl = resp.SignStr
	}
	return
}

func (i *VideoInfo) InitUploadJob() (resp *InitiateMultipartUploadResult, err error) {
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client: i.Client.httpClient,
		Url:    i.authorizedUrl,
		Method: http.MethodPost,
		Header: map[string]string{
			"Content-Type": "video/mp4",
		},
	})
	if err != nil {
		return
	}
	if err = xml.Unmarshal(b, &resp); err != nil {
		return
	}
	i.uploadResult = *resp
	return
}

func (i *VideoInfo) GetAssetAuthority() (resp *GetAssetAuthorityResponse, err error) {
	resp = new(GetAssetAuthorityResponse)
	payload := &GetAssetAuthorityRequest{
		HttpVerb:    http.MethodPut,
		ContentType: "video/mp4",
		Bucket:      i.target.Bucket,
		ObjectKey:   i.target.Object,
		UploadId:    i.uploadResult.UploadId,
		PartNumber:  1,
		// ContentMd5:  "MmIzMGQ3MjQzZDc2NzA3MjBmMzEzY2JlY2Y4NDRhZjA=",
	}
	b, _, err := utils.DoHttpRequest(&utils.HttpRequestConfig{
		Client:  i.Client.httpClient,
		Url:     fmt.Sprintf(GET_ASSET_AUTH_URL, i.Client.config.vod.projectId),
		Method:  http.MethodGet,
		TagName: TAG_JSON,
		Header: map[string]string{
			HUAWEI_AUTH_TOKEN_NAME: i.Client.config.auth.token,
		},
		Payload: payload,
	})
	if err != nil {
		return
	}
	fmt.Println(string(b))
	err = json.Unmarshal(b, &resp)
	return
}
