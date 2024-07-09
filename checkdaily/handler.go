package checkdaily

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/JinY3/gopkg/filex"
	"github.com/JinY3/gopkg/logx"
	"github.com/chromedp/chromedp"
)

var Token string

func Check(username, password string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	url := "http://ehall.njc.ucas.ac.cn/qljfwapp/sys/lwPsXykApp/index.do?#/dledcx"

	var dataValue string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#username"),
		chromedp.WaitVisible("#password"),
		chromedp.SendKeys("#username", username),
		chromedp.SendKeys("#password", password),
		chromedp.Click("#login_submit", chromedp.NodeVisible),
		chromedp.AttributeValue(`//*[@name="REMAINEQ"]`, "data-value", &dataValue, nil),
	)
	if err != nil {
		logx.MyAll.Errorf("查询电量失败: %s", err)
		sendEmail("查询电量失败")
	} else {
		logx.MyAll.Infof("查询电量成功: %s", dataValue)
		sendEmail(fmt.Sprintf("当前电量: %s", dataValue))
		appendFile("./value.txt", fmt.Sprintf("%s\n", dataValue))
		appendFile("./time.txt", fmt.Sprintf("%s\n", time.Now().Format("2006-01-02")))
	}

}

// 向指定文件追加写入内容
func appendFile(filename string, content string) {
	filex.CreateFile(filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logx.MyAll.Errorf("打开文件失败: %s", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		logx.MyAll.Errorf("写入文件失败: %s", err)
		return
	}
}

func sendEmail(msg string) {
	if Token == "" {
		logx.MyAll.Errorf("未设置pushplus token")
		return
	}
	data := []byte(fmt.Sprintf("{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\"}", Token, "查询电量监控", msg))
	response, err := http.Post("http://www.pushplus.plus/send", "application/json", bytes.NewBuffer(data))
	if err != nil {
		logx.MyAll.Errorf("发送邮件失败: %s", err)
		return
	}
	defer response.Body.Close()
}
