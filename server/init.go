package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/JinY3/gopkg/logx"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logx.MyAll.Out,
	}))
	r.Use(cors)
	r.LoadHTMLFiles("index.html")
	r.GET("/:homeid", func(c *gin.Context) {
		homeid := c.Param("homeid")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"homeid": homeid,
		})
	})
	r.GET("/static/echarts.js", func(c *gin.Context) {
		c.File("static/echarts.js")
	})
	r.GET("/static/my.js/:homeid", func(c *gin.Context) {
		homeid := c.Param("homeid")
		// 从my.js文件中读取内容并替换其中的ID
		myjsContent, err := os.ReadFile("static/my.js")
		if err != nil {
			logx.MyAll.Errorf("读取my.js文件失败: %s", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"message": "读取文件失败",
			})
			return
		}
		c.Data(http.StatusOK, "text/javascript", []byte(strings.ReplaceAll(string(myjsContent), "{{.homeid}}", homeid)))
	})
	r.GET("/data/:homeid", data)
}
