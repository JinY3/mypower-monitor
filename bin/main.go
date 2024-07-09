package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/JinY3/gopkg/logx"
	"github.com/JinY3/mypower-monitor/checkdaily"
	"github.com/JinY3/mypower-monitor/server"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

var (
	account string
	pwd     string
)

func init() {
	flag.StringVar(&account, "account", "", "学号")
	flag.StringVar(&pwd, "pwd", "", "密码")
	flag.StringVar(&checkdaily.Token, "token", "", "pushplus token")
}

func main() {
	flag.Parse()
	ctlCtx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	go func(_account, _pwd string) {
		checkdaily.Check(_account, _pwd)
		for {
			select {
			case <-ctlCtx.Done():
				return
			case <-time.After(24 * time.Hour):
				checkdaily.Check(_account, _pwd)
			}
		}
	}(account, pwd)

	port := 7001 // master的端口
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	server.Init(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// 启动端口监听, 失败则停止阻塞
	logx.MyAll.Infof("server listen on %d", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logx.MyAll.Fatalf("listen: %s", err)
	}

	// 关闭master
	logx.MyAll.Info("server shutdown")
}
