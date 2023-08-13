package Comman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Global"
	Requests "github.com/lizazacn/requests"
	"net/http"
	"strings"
	"time"
)

// GetAccessToken 获取AccessToken
func GetAccessToken(apiKey, secretKey string) (string, error) {
	if apiKey == "" || secretKey == "" {
		apiKey = Global.Conf.APIKey
		secretKey = Global.Conf.SecretKey
	}
	Global.EndTime = time.Now().AddDate(0, 0, 20)
	fmt.Printf("当前Token：%s,\n 预计过期时间: %v\n", Global.AccessToken, Global.EndTime)
	var url = fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", apiKey, secretKey)
	response, err := Requests.Requests(http.MethodGet, url, nil, nil, true, true, false, nil)
	if err != nil {
		return "", err
	}
	var accessToken, ok = response.Json["access_token"]
	if !ok || accessToken.(string) == "" {
		return "", errors.New("获取access_token异常！")
	}
	Global.AccessToken = accessToken.(string)
	return accessToken.(string), nil
}

func GetAnswer(question string) (string, error) {
	Global.MaxCallback++
	defer func() {
		Global.MaxCallback--
	}()
	if Global.MaxCallback >= 2 {
		return "系统异常请稍后重试！", nil
	}
	if question == "" {
		return "请求数据不能为空！", errors.New("请求数据不能为空！")
	}
	if len(question) >= 2000 {
		return "请求数据超出最大长度限制，最大长度限制为2000字符！", errors.New("请求数据超出最大长度限制，最大长度限制为2000字符！")
	}
	var url = fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/plugin/%s/?access_token=%s", Global.Conf.ServiceName, Global.AccessToken)
	//var url = fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions?access_token=%s", Global.AccessToken)
	//var data = make(map[string][]map[string]string)
	//data["messages"] = []map[string]string{
	//	{
	//		"role":    "user",
	//		"content": question,
	//	},
	//}
	var data = make(map[string]string)
	data["query"] = question
	buffer, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	var buf = bytes.NewBuffer(buffer)
	var header = http.Header{}
	header.Add("Content-Type", "application/json")
	if !NotOverdue() {
		newToken, err := GetAccessToken("", "")
		if err != nil {
			return "服务端异常，请稍后重试！", err
		}
		Global.AccessToken = newToken
		url = fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/plugin/%s/?access_token=%s", Global.Conf.ServiceName, newToken)
	}
	var tryCount = 0
reTry:
	response, err := Requests.Requests(http.MethodPost, url, buf, header, true, true, false, nil)
	if err != nil {
		tryCount++
		newToken, err := GetAccessToken("", "")
		if err != nil {
			return "服务端异常，请稍后重试！", err
		}
		Global.AccessToken = newToken
		url = fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/plugin/%s/?access_token=%s", Global.Conf.ServiceName, newToken)
		if tryCount > 10 {
			return "服务端异常次数过多，请联系管理员！", err
		}
		goto reTry
	}
	var res, ok = response.Map["result"].(string)
	if !ok {
		return "", errors.New("响应数据格式异常！")
	}
	return DataCollate(res), nil
}

func NotOverdue() bool {
	return time.Now().After(Global.EndTime)
}

func DataCollate(in string) string {
	var out string
	dataList := strings.Split(in, "[")
	if len(dataList) < 2 {
		return in
	}
	dataList = dataList[:len(dataList)-1]
	out = strings.Join(dataList, "[")
	return out
}
