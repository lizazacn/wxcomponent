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
			ToUserName:   fmt.Sprintf("<![CDATA[%s]]>", msgInfo.FromUserName),
			FromUserName: fmt.Sprintf("<![CDATA[%s]]>", msgInfo.ToUserName),
			Content:      fmt.Sprintf("<![CDATA[%s]]>", "解析请求数据异常！"),
			MsgType:      "<![CDATA[text]]>",
			CreateTime:   time.Now().Unix(),
		})
		return
	}
	log.Println(*msgInfo)

	var waitStat = true
	go func() {
		if waitStat {
			ctx.String(http.StatusOK, "success")
			time.Sleep(3 * time.Second)
		}
	}()

	answer, err := Comman.GetAnswer(msgInfo.Content)
	log.Printf("响应数据：%s\n", answer)
	waitStat = false
	if err != nil {
		ctx.XML(http.StatusOK, Struct.XML{
			ToUserName:   fmt.Sprintf("<![CDATA[%s]]>", msgInfo.FromUserName),
			FromUserName: fmt.Sprintf("<![CDATA[%s]]>", msgInfo.ToUserName),
			Content:      fmt.Sprintf("<![CDATA[%s]]>", "解析请求数据异常！"),
			MsgType:      "<![CDATA[text]]>",
			CreateTime:   time.Now().Unix(),
		})
		return
	}
	ctx.XML(http.StatusOK, Struct.XML{
		ToUserName:   fmt.Sprintf("<![CDATA[%s]]>", msgInfo.FromUserName),
		FromUserName: fmt.Sprintf("<![CDATA[%s]]>", msgInfo.ToUserName),
		Content:      fmt.Sprintf("<![CDATA[%s]]>", answer),
		MsgType:      "<![CDATA[text]]>",
		CreateTime:   time.Now().Unix(),
	})
}
