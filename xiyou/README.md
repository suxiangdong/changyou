# 客户 畅游
## 项目 西游

## 编译
```
go1.19

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/xiyou -ldflags="-s -w" main.go
```

# 文件说明
## parser 业务解析
## prot protobuf文件
## real 线上正式运行配置文件