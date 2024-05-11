package main

import (
	"GoChatServer/helper"
	"GoChatServer/router"
	"GoChatServer/ws"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	helper.InitConfig("./config") //初始化配置
	helper.InitLogger()           //初始化日志
	helper.InitSqlLogger()        //初始化日志
	helper.InitChatDatabase()     //初始化DB
	helper.InitWeiXin()           //初始化微信实例
	//helper.InitRedis()

	//_ = gin.Default()                             //创建gin实例
	r := gin.New()                           //创建gin实例
	router.InitRoute(r)                      //初始化API路由
	ws.InitWebsocket(r)                      //初始化Ws
	r.Static("/static/uploads", "./uploads") //开启静态文件服务

	//优雅启动
	run(r)
}

func run(r *gin.Engine) {
	srv := &http.Server{
		Addr:    helper.Configs.Server.Address,
		Handler: r,
	}

	go func() {
		// 服务连接
		helper.Logger.Info("Server Start ...")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	helper.Logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		helper.Logger.Infof("Server Shutdown: %s", err.Error())
	}
	helper.Logger.Infof("Server exiting")

}
