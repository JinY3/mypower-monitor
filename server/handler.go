package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JinY3/gopkg/logx"
	"github.com/gin-gonic/gin"
)

// CORS中间件
func cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

// 数据解析
func data(c *gin.Context) {
	homeid := c.Param("homeid")
	// 读取value.txt文件并生成数组
	value, err := readData(fmt.Sprintf("data/%s/value.txt", homeid))
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
		valueFloat[i], _ = strconv.ParseFloat(v, 64)
	}
	// 将 value 数据转化为 差分数据
	valueDiff := make([]float64, len(valueFloat)-1)
	for i := 1; i < len(valueFloat); i++ {
		valueDiff[i-1] = valueFloat[i-1] - valueFloat[i]
	}

	// 读取timeTxt.txt文件并生成数组
	timeTxt, err := readData(fmt.Sprintf("data/%s/time.txt", homeid))
	if err != nil {
		logx.MyAll.Errorf("读取time.txt文件失败: %s", err)
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"message": "读取文件失败",
		})
		return
	}
	// 计算天数差
	for i := 1; i < len(timeTxt)-1; i++ {
		// 解析时间
		// logx.MyAll.Debugf("%s", timeTxt[i-1])
		t1, err := time.Parse("2006-01-02", timeTxt[i-1][:10])
		if err != nil {
			logx.MyAll.Errorf("解析时间失败: %s", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"message": "解析时间失败",
			})
			return
		}
		t2, err := time.Parse("2006-01-02", timeTxt[i])
		if err != nil {
			logx.MyAll.Errorf("解析时间失败: %s", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"message": "解析时间失败",
			})
			return
		}

		dayDiff := t2.Sub(t1).Hours() / 24
		valueDiff[i-1] = valueDiff[i-1] / dayDiff
		if valueDiff[i-1] < 0 {
			valueDiff[i-1] = -1
			timeTxt[i] = fmt.Sprintf("%s (近期充值, 无法计算近期耗电量)", timeTxt[i])
			continue
		}
		if dayDiff < 2 {
			continue
		}
		timeTxt[i] = fmt.Sprintf("%s (近 %d 天平均)", timeTxt[i], int(dayDiff))
	}
	// 保留两位小数
	for i := 0; i < len(valueDiff); i++ {
		valueDiff[i] = float64(int(valueDiff[i]*100)) / 100
	}
	c.JSON(http.StatusOK, gin.H{
		"current": valueDiff[len(valueDiff)-1],
		"time":    timeTxt[max(1, len(timeTxt)-31) : len(timeTxt)-1],
		"value":   valueDiff[max(0, len(valueDiff)-31) : len(valueDiff)-1],
	})
}

func readData(filePath string) ([]string, error) {
	// 读取文件到字节数组
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\n"), nil
}
