package main

import (
	"context"
	"fmt"
	"time"

	"github.com/JinY3/gopkg/filex"
	"github.com/JinY3/gopkg/logx"
	"github.com/JinY3/mypower-monitor/checkdaily"
	"github.com/JinY3/mypower-monitor/server"
	"github.com/gin-gonic/gin"
)

var checkdailyYaml struct {
	Token string            `json:"token"`
	Users []checkdaily.User `json:"users"`
}

func init() {
	filex.ReadConfig("config", "userlist", &checkdailyYaml)
	logx.MyAll.Debugf("读取用户列表成功: %v", checkdailyYaml)
}

func main() {
	ctlCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(yaml struct {
		Token string            `json:"token"`
		Users []checkdaily.User `json:"users"`
	}) {
		// for _, user := range yaml.Users {
		// 	go user.Check(yaml.Token)
		// }
		for {
			select {
			case <-ctlCtx.Done():
				return
			case <-time.After(24 * time.Hour):
				for _, user := range yaml.Users {
					go user.Check(yaml.Token)
				}
			}
		}
	}(checkdailyYaml)

	port := 7001 // master的端口
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	server.Init(r)

	logx.MyAll.Infof("server start at :%d", port)
	r.Run(fmt.Sprintf(":%d", port))
}
