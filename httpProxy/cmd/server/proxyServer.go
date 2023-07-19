package main

import (
	"fmt"
	"github.com/bilinerz/freeHttpProxy/httpProxy"
	"github.com/bilinerz/freeHttpProxy/httpProxy/cmd"
	"log"
	"net"
)

var version = "测试版"

func main() {
	//读取默认配置
	config := &cmd.Config{}
	config.ReadConfig()

	//启动服务端
	listServer, err := httpProxy.StartServer(config.Key, config.ServerAddr)
	if err != nil {
		log.Fatalln("NewServer Error：", err)
	}

	log.Fatalln(
		listServer.Listen(func(listenAddr net.Addr) {
			log.Println(
				fmt.Sprintf(`FreeHTTPProxy-Server:%s 启动成功，服务监听地址：%s`, version, listenAddr))
		}))
}
