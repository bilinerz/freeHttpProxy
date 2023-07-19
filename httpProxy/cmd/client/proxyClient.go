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

	//启动客户端
	listClient, err := httpProxy.StartClient(config.Key, config.ClientAddr, config.ServerAddr)
	if err != nil {
		log.Fatalln("启动客户端失败：", err)
	}

	log.Fatalln(
		listClient.Listen(func(listenAddr net.Addr) {
			log.Println(
				fmt.Sprintf(`FreeHTTPProxy-Client:%s 启动成功，服务监听地址：%s`, version, listenAddr))
		}))
}
