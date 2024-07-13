package server

import (
	"net/http"
	"os"
	"strconv"
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
		c.File("static/index.html")
	})
	r.GET("/echarts.js", func(c *gin.Context) {
		c.File("static/echarts.js")
	})
	r.GET("/my.js", func(c *gin.Context) {
		c.File("static/my.js")
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
		// 将 value 数据转化为 float64 类型
		valueFloat := make([]float64, len(value))
		for i, v := range value {
			if v == "" {
				continue
			}
			valueFloat[i], err = strconv.ParseFloat(v, 64)
		}
		// 将 value 数据转化为 差分数据
		valueDiff := make([]float64, len(valueFloat)-1)
		for i := 1; i < len(valueFloat); i++ {
			valueDiff[i-1] = valueFloat[i-1] - valueFloat[i]
			// logx.MyAll.Debugf("valueDiff[%d]: %.2f", i-1, valueDiff[i-1])
		}
		// 保留两位小数
		for i := 0; i < len(valueDiff); i++ {
			valueDiff[i] = float64(int(valueDiff[i]*100)) / 100
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
			"current": valueDiff[len(valueDiff)-1],
			"time":    time[1:],
			"value":   valueDiff[:len(valueDiff)-1],
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
