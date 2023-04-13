package config

var Config DownHubConfig

// 配置信息
type DownHubConfig struct {
	BDUSS               string `json:"bduss"`
	Token               string `json:"token"`
	DingtalkAccessToken string `json:"dingtalk_access_token"`
	DingtalkSecret      string `json:"dingtalk_secret"`
}
