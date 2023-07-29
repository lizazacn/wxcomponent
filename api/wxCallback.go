package API

import (
	"encoding/xml"
	"fmt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Comman"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func WxCallback(ctx *gin.Context) {
	//var data, _ = io.ReadAll(ctx.Request.Body)
	//fmt.Println(string(data))
	//return
	var msgInfo = new(Struct.MsgInfo)
	err := ctx.Bind(msgInfo)
	if err != nil {
		ctx.XML(http.StatusOK, Struct.XML{
			ToUserName:   msgInfo.FromUserName,
			FromUserName: msgInfo.ToUserName,
			Content:      "解析请求数据异常！",
			MsgType:      msgInfo.MsgType,
			CreateTime:   time.Now().Unix(),
		})
		return
	}
	fmt.Println(*msgInfo)
	var wxMsg = new(Struct.XML)
	err = xml.Unmarshal([]byte(msgInfo.Data), wxMsg)
	if err != nil {
		ctx.XML(http.StatusOK, Struct.XML{
			ToUserName:   msgInfo.FromUserName,
			FromUserName: msgInfo.ToUserName,
			Content:      "解析请求数据异常！",
			MsgType:      msgInfo.MsgType,
			CreateTime:   time.Now().Unix(),
		})
		return
	}
	answer, err := Comman.GetAnswer(msgInfo.Data)
	if err != nil {
		ctx.XML(http.StatusOK, Struct.XML{
			ToUserName:   msgInfo.FromUserName,
			FromUserName: msgInfo.ToUserName,
			Content:      "获取答案异常！",
			MsgType:      msgInfo.MsgType,
			CreateTime:   time.Now().Unix(),
		})
		return
	}
	ctx.XML(http.StatusOK, Struct.XML{
		ToUserName:   msgInfo.FromUserName,
		FromUserName: msgInfo.ToUserName,
		Content:      answer,
		MsgType:      msgInfo.MsgType,
		CreateTime:   time.Now().Unix(),
	})
}