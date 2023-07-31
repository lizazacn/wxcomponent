package API

import (
	"bytes"
	"encoding/json"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Comman"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx"
	"github.com/gin-gonic/gin"
	Requests "github.com/lizazacn/requests"
	"log"
	"net/http"
	"time"
)

func WxCallback(ctx *gin.Context) {
	//var data, _ = io.ReadAll(ctx.Request.Body)
	//fmt.Printf("请求数据：%s", string(data))
	//return
	var appid = ctx.GetHeader("X-Wx-From-Appid")
	var msgInfo = new(Struct.XML)
	//err := json.Unmarshal(data, msgInfo)
	err := ctx.Bind(msgInfo)
	if err != nil {
		//ctx.XML(http.StatusOK, Struct.XML{
		//	ToUserName:   fmt.Sprintf("<![CDATA[%s]]>", msgInfo.FromUserName),
		//	FromUserName: fmt.Sprintf("<![CDATA[%s]]>", msgInfo.ToUserName),
		//	Content:      fmt.Sprintf("<![CDATA[%s]]>", "解析请求数据异常！"),
		//	MsgType:      "<![CDATA[text]]>",
		//	CreateTime:   time.Now().Unix(),
		//})
		ctx.XML(http.StatusOK, Struct.XML{
			ToUserName:   msgInfo.FromUserName,
			FromUserName: msgInfo.ToUserName,
			MsgType:      "text",
			CreateTime:   time.Now().Unix(),
			Content:      "解析请求数据异常!",
		})
		return
	}
	log.Println(*msgInfo)

	ctx.String(http.StatusOK, "success")
	answer, err := Comman.GetAnswer(msgInfo.Content)
	log.Printf("响应数据：%s\n", answer)
	//waitStat = false
	if err != nil {
		//ctx.XML(http.StatusOK, Struct.XML{
		//	ToUserName:   fmt.Sprintf("<![CDATA[%s]]>", msgInfo.FromUserName),
		//	FromUserName: fmt.Sprintf("<![CDATA[%s]]>", msgInfo.ToUserName),
		//	Content:      fmt.Sprintf("<![CDATA[%s]]>", "解析请求数据异常！"),
		//	MsgType:      "<![CDATA[text]]>",
		//	CreateTime:   time.Now().Unix(),
		//})
		_ = SendCustomerServiceMsg(appid, "解析请求数据异常!", msgInfo.FromUserName)
		//ctx.XML(http.StatusOK, Struct.XML{
		//	ToUserName:   msgInfo.FromUserName,
		//	FromUserName: msgInfo.ToUserName,
		//	MsgType:      "text",
		//	CreateTime:   time.Now().Unix(),
		//	Content:      "解析请求数据异常!",
		//})
		return
	}
	//ctx.XML(http.StatusOK, Struct.XML{
	//	ToUserName:   fmt.Sprintf("<![CDATA[%s]]>", msgInfo.FromUserName),
	//	FromUserName: fmt.Sprintf("<![CDATA[%s]]>", msgInfo.ToUserName),
	//	Content:      fmt.Sprintf("<![CDATA[%s]]>", answer),
	//	MsgType:      "<![CDATA[text]]>",
	//	CreateTime:   time.Now().Unix(),
	//})
	_ = SendCustomerServiceMsg(appid, answer, msgInfo.FromUserName)
	//ctx.XML(http.StatusOK, Struct.XML{
	//	ToUserName:   msgInfo.FromUserName,
	//	FromUserName: msgInfo.ToUserName,
	//	MsgType:      "text",
	//	CreateTime:   time.Now().Unix(),
	//	Content:      answer,
	//})
}

func SendCustomerServiceMsg(appid, msg, toUser string) error {
	log.Printf("APPID: %s", appid)
	accessToken, err := wx.GetAuthorizerAccessToken(appid)
	if err != nil {
		log.Println(err)
		return err
	}
	//url := "http://api.weixin.qq.com/cgi-bin/message/custom/send"
	//url := "http://api.weixin.qq.com/cgi-bin/message/custom/send?from_appid=" + appid
	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + accessToken
	var data = make(map[string]interface{})
	data["touser"] = toUser
	data["msgtype"] = "text"
	data["text"] = map[string]string{"content": msg}
	buffer, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}
	var buf = bytes.NewBuffer(buffer)
	var header = http.Header{}
	header.Add("Content-Type", "application/json")
	response, err := Requests.Requests(http.MethodPost, url, buf, header, true, false, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(response.Text))
	return nil
}
