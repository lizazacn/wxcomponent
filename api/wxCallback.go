package API

import (
	"encoding/json"
	"fmt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Comman"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

func WxCallback(ctx *gin.Context) {
	var data, _ = io.ReadAll(ctx.Request.Body)
	fmt.Printf("请求数据：%s", string(data))
	//return
	var msgInfo = new(Struct.XML)
	err := json.Unmarshal(data, msgInfo)
	//err := ctx.Bind(msgInfo)
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
	log.Println(*msgInfo)

	answer, err := Comman.GetAnswer(msgInfo.Content)
	log.Printf("响应数据：%s\n", answer)
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
