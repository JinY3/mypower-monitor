# mypower-monitor

ucasnj 电量监控

## usage

```shell
cd {project_path}
# ls
# bin  checkdaily  cmd  default.log  go.mod  go.sum  LICENSE  README.md  server  static  time.txt  value.txt

# 执行文件放到{project_path}/bin下
go build -o bin/run cmd/main.go

# config/userlist.yaml 存放用户信息
`example
Token: xxxxxxxxxxxxxxx  # pushplus token
Users:
  - Account: 2023xxxxxxx
    Password: "mypassword1"
    Homeid: b905
  - Account: 2022xxxxxxx
    Password: "mypassword2"
    To: xxxxxxxxxxxxxxxx  # pushplus friend id
    Homeid: a801
`

# use
bin/run
```

## token

使用 [pushplus](https://www.pushplus.plus/) 推送消息，需要传入token, 若不需要可以置空

## 后续

改用 github action 运行, 实现无服务器使用
