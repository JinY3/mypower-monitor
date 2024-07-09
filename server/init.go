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
	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})
	r.GET("/echarts.js", func(c *gin.Context) {
		c.File("echarts.js")
	})
	r.GET("/my.js", func(c *gin.Context) {
		c.File("my.js")
	})
	r.GET("/data", func(c *gin.Context) {
		// 读取value.txt文件并生成数组
		value, err := readFile("value.txt")
		if err != nil {
			logx.MyAll.Errorf("读取value.txt文件失败: %s", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"message": "读取文件失败",
			})
			return
		}
		// 读取time.txt文件并生成数组
		time, err := readFile("time.txt")
		if err != nil {
			logx.MyAll.Errorf("读取time.txt文件失败: %s", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"message": "读取文件失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"time":  time,
			"value": value,
		})
	})
}

func readFile(s string) ([]string, error) {
	// 读取文件到字节数组
	bytes, err := os.ReadFile(s)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\n"), nil
}
