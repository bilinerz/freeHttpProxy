# freeHttpProxy

练手之作🙃

基于 [Lightsocks](https://github.com/gwuhaolin/lightsocks) 实现一个简单的 HTTP加密代理服务。

```json
{
	"server": ":1992",
	"client": ":1080",
	"key": ""
}
```

```go build 编译好 Server端 和 Client端```
服务端和客户端的 key 必须使用生成的一致密钥

在服务器运行 server端程序产生配置文件后，修改server监听地址 :1992 为自己需要的地址，再重新运行 server端

在客户端运行 Client端程序产生配置文件后，修改server监听地址 :1992 为中转服务器监听地址，修改client监听地址为自己客户端监听的地址，再重新运行 client端
