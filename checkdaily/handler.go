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

func (user *User) Check(token string) {
	tryCounter := 3

	for i := 0; i < tryCounter; i++ {
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
			chromedp.SendKeys("#username", user.Account),
			chromedp.SendKeys("#password", user.Password),
			chromedp.Click("#login_submit", chromedp.NodeVisible),
			chromedp.AttributeValue(`//*[@name="REMAINEQ"]`, "data-value", &dataValue, nil),
		)
		if err != nil {
			logx.MyAll.WithField("user", user.Homeid).Errorf("查询电量失败: %s", err)
		} else {
			logx.MyAll.WithField("user", user.Homeid).Infof("查询电量成功: %s", dataValue)
			user.send(token, fmt.Sprintf("当前电量: %s", dataValue), fmt.Sprintf("http://157.0.19.2:10063/mypower/%s", user.Homeid))
			appendFile(fmt.Sprintf("data/%s/value.txt", user.Homeid), fmt.Sprintf("%s\n", dataValue))
			appendFile(fmt.Sprintf("data/%s/time.txt", user.Homeid), fmt.Sprintf("%s\n", time.Now().Format("2006-01-02")))
			return
		}
	}

	time.Sleep(5 * time.Minute)

	user.send(token, "查询电量失败", fmt.Sprintf("http://157.0.19.2:10063/mypower/%s", user.Homeid))
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

func (user *User) send(token, title, msg string) {
	if token == "" {
		logx.MyAll.Debugf("未设置pushplus token")
		return
	}
	var data []byte
	if user.To == "" {
		data = []byte(fmt.Sprintf("{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\"}", token, title, msg))
	} else {
		data = []byte(fmt.Sprintf("{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\", \"to\": \"%s\"}", token, title, msg, user.To))
	}
	response, err := http.Post("http://www.pushplus.plus/send", "application/json", bytes.NewBuffer(data))
	if err != nil {
		logx.MyAll.WithField("user", user.Homeid).Errorf("发送邮件失败: %s", err)
		return
	}
	defer response.Body.Close()
}
