package Struct

// Conf 基础配置文件
type Conf struct {
	Listen      string `json:"listen"`       // 监听地址
	Port        int    `json:"port"`         // 端口号
	APIKey      string `json:"api_key"`      // APIKey
	SecretKey   string `json:"secret_key"`   // SecretKey
	ServiceName string `json:"service_name"` // 服务名
}

type MsgInfo struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	MsgType      string `json:"MsgType"`
	CreateTime   int64  `json:"CreateTime"`
	Data         string `json:"Data"`
	Content      string `json:"Content" xml:"Content"`
}

type XML struct {
	AppId        string `json:"appid" xml:"appid"`
	ToUserName   string `json:"ToUserName" xml:"ToUserName"`
	FromUserName string `json:"FromUserName" xml:"FromUserName"`
	MsgType      string `json:"MsgType" xml:"MsgType"`
	CreateTime   int64  `json:"CreateTime" xml:"CreateTime"`
	Content      string `json:"Content" xml:"Content"`
}
