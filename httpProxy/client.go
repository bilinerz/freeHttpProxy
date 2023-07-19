package httpProxy

import (
	"log"
	"net"
)

// 创建一个本地端客户端
// 客户端的职责是:
// 1. 监听来自本机浏览器的代理请求
// 2. 转发前加密数据
// 3. 转发http数据到墙代理服务端
// 4. 把代理服务端返回的数据转发给用户的浏览器

// ListenClient 客户端结构体
type ListenClient struct {
	Cipher     *Cipher
	LocalAddr  *net.TCPAddr
	ServerAddr *net.TCPAddr
}

// StartClient 启动客户端
func StartClient(key, clientAddr, serverAddr string) (*ListenClient, error) {
	k, err := ParseBaseKey(key)
	if err != nil {
		return nil, err
	}
	structLocalAddr, err := net.ResolveTCPAddr("tcp", clientAddr)
	if err != nil {
		return nil, err
	}

	structServerAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	return &ListenClient{
		Cipher:     NewCipher(k),
		LocalAddr:  structLocalAddr,
		ServerAddr: structServerAddr,
	}, nil
}

// Listen 本地端启动监听，接收来自本机浏览器的连接
func (client *ListenClient) Listen(didListen func(listenAddr net.Addr)) error {
	return ListenSecureTCP(client.LocalAddr, client.handleConn, client.Cipher, didListen)
}

// 处理http请求
func (client *ListenClient) handleConn(clientConn *SecureTCPConn) {
	defer clientConn.Close()

	log.Println("浏览器地址：", clientConn.RemoteAddr())

	//与代理服务器建立连接
	proxyServerConn, err := DialTCPSecure(client.ServerAddr, client.Cipher)
	if err != nil {
		log.Println("与代理服务器建立连接失败：", err)
		return
	}
	defer proxyServerConn.Close()

	// 进行转发
	// 从 proxyServer 读取数据发送到 clientConn
	go func() {
		err := proxyServerConn.DecodeCopy(clientConn)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			clientConn.Close()
			proxyServerConn.Close()
		}
	}()
	// 从 clientConn 发送加密数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	clientConn.EncodeCopy(proxyServerConn)
}
