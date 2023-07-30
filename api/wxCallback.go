package API

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Comman"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func WxCallback(ctx *gin.Context) {
	var data, _ = io.ReadAll(ctx.Request.Body)
	fmt.Println(string(data))
	//return
	var msgInfo = new(Struct.MsgInfo)
	err := json.Unmarshal(data, msgInfo)
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
	fmt.Println(*wxMsg)
	answer, err := Comman.GetAnswer(msgInfo.Data)
	fmt.Printf("响应数据：%s\n", answer)
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
