package cls_vod

import (
	"sync"
)

type MediaType string

const (
	SCHEMA   = "https://"
	IAM_HOST = "iam.cn-east-2.myhuaweicloud.com"
	VOD_HOST = "vod.cn-east-2.myhuaweicloud.com"
)

const (
	TAG_JSON = "json"
)

const (
	HUAWEI_TOKEN_NAME      = "x-subject-token"
	HUAWEI_AUTH_TOKEN_NAME = "X-Auth-Token"
)

const (
	VOD_CATAGORY_TEST   = 119396
	VOD_CATAGORY_OTHERS = -1
)

const (
	AUTH_TOKEN_URL          = SCHEMA + IAM_HOST + "/v3/auth/tokens"
	GET_TEMP_SECU_TOKEN_URL = SCHEMA + IAM_HOST + "/v3.0/OS-CREDENTIAL/securitytokens"
	CREATE_ASSET_URL        = SCHEMA + VOD_HOST + "/v1.0/%s/asset"
	GET_ASSET_DETAIL_URL    = SCHEMA + VOD_HOST + "/v1.0/%s/asset/details"
	GET_ASSET_AUTH_URL      = SCHEMA + VOD_HOST + "/v1.1/%s/asset/authority"
	ASSET_UPLOAD_CHECK_URL  = SCHEMA + VOD_HOST + "/v1.0/%s/asset/status/uploaded"
	PUBLISH_VIDEO_URL       = SCHEMA + VOD_HOST + "/v1.0/%s/asset/status/publish"
)

const (
	VEDIO_TYPE_MP4  MediaType = "MP4"
	VEDIO_TYPE_TS   MediaType = "TS"
	VEDIO_TYPE_MOV  MediaType = "MOV"
	VEDIO_TYPE_MXF  MediaType = "MXF"
	VEDIO_TYPE_MPG  MediaType = "MPG"
	VEDIO_TYPE_FLV  MediaType = "FLV"
	VEDIO_TYPE_WMV  MediaType = "WMV"
	VEDIO_TYPE_AVI  MediaType = "AVI"
	VEDIO_TYPE_M4V  MediaType = "M4V"
	VEDIO_TYPE_F4V  MediaType = "F4V"
	VEDIO_TYPE_MPEG MediaType = "MPEG"
	VEDIO_TYPE_3GP  MediaType = "3GP"
	VEDIO_TYPE_ASF  MediaType = "ASF"
	VEDIO_TYPE_MKV  MediaType = "MKV"
	VEDIO_TYPE_HLS  MediaType = "HLS"
)

var (
	lock        = new(sync.Mutex)
	redisLocker = new(sync.Mutex)
)
