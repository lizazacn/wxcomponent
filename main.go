package main

import (
	"fmt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Comman"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Global"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Utils"
	API "github.com/WeixinCloud/wxcloudrun-wxcomponent/api"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/inits"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/routers"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/sync/errgroup"
	"os"
	"strconv"
)

func main() {
	log.Infof("system begin")
	if err := inits.Init(); err != nil {
		log.Errorf("inits failed, err:%v", err)
		return
	}
	log.Infof("inits.Init Succ")

	var g errgroup.Group

	// 内部服务
	g.Go(func() error {
		r := routers.InnerServiceInit()
		if err := r.Run("127.0.0.1:8081"); err != nil {
			log.Error("startup inner service failed, err:%v", err)
			return err
		}
		return nil
	})

	// 外部服务
	g.Go(func() error {
		r := routers.Init()
		if err := r.Run(":80"); err != nil {
			log.Error("startup service failed, err:%v", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		var filePath = "etc/config.json"
		conf, err := Utils.LoadConf(filePath)
		//if err != nil {
		//	return err
		//}
		if conf == nil {
			conf = new(Struct.Conf)
		}
		Global.Conf = conf
		if Global.Conf == nil {
			return err
		}
		Global.Conf.APIKey = os.Getenv("API_KEY")
		Global.Conf.SecretKey = os.Getenv("SECRET_KEY")
		Global.Conf.ServiceName = os.Getenv("SERVICE_NAME")
		Global.Conf.Listen = os.Getenv("LISTEN")
		Global.Conf.Port, _ = strconv.Atoi(os.Getenv("PORT"))
		// 初始化Token
		_, _ = Comman.GetAccessToken("", "")
		var r = gin.Default()
		var api = r.Group("/api")

		{
			api.POST("/wx_callback", API.WxCallback)
		}
		var listenAddr = "127.0.0.1:8099"
		if conf.Listen != "" || conf.Port != 0 {
			listenAddr = fmt.Sprintf("%s:%d", conf.Listen, conf.Port)
		}
		err = r.Run(listenAddr)
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error(err)
	}
}
