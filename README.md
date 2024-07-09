# mypower-monitor

ucasnj 电量监控

## usage

```shell
cd {project_path}
# ls
# bin  checkdaily  cmd  default.log  go.mod  go.sum  LICENSE  README.md  server  static  time.txt  value.txt

# 也可以使用release，二进制文件放到{project_path}/bin下
go build -o bin/run cmd/main.go

# help
bin/run -h

# example
bin/run -account 2023*********** -pwd ********* -token **************
```

## token

使用 [pushplus](https://www.pushplus.plus/) 推送消息，需要传入token

## 后续

改用 github action 运行, 实现无服务器使用