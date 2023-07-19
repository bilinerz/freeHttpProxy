package httpProxy

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

// 服务端的职责是:
// 1. 监听来自本地代理客户端的请求
// 2. 解密本地代理客户端请求的数据，解析 HTTP 协议，连接用户浏览器真正想要连接的远程服务器
// 3. 转发用户浏览器真正想要连接的远程服务器返回的数据的加密后的内容到本地代理客户端

type ListenServer struct {
	Cipher     *Cipher
	ListenAddr *net.TCPAddr
}

// StartServer 启动服务端
func StartServer(key, listenAddr string) (*ListenServer, error) {
	k, err := ParseBaseKey(key)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &ListenServer{
		Cipher:     NewCipher(k),
		ListenAddr: structListenAddr,
	}, nil
}

// Listen 运行服务端并且监听来自本地代理客户端的请求
func (lsServer *ListenServer) Listen(didListen func(listenAddr net.Addr)) error {
	return ListenSecureTCP(lsServer.ListenAddr, lsServer.handleConn, lsServer.Cipher, didListen)
}

// 处理网络请求
func (lsServer *ListenServer) handleConn(localConn *SecureTCPConn) {
	defer localConn.Close()

	log.Println("客户端地址:", localConn.RemoteAddr())

	// 用来存放客户端HTTP请求数据的缓冲区
	buf := make([]byte, 256)

	//解密数据后读取真实的HTTP请求数据
	n, err := localConn.DecodeRead(buf)
	if err != nil {
		log.Println("从客户端读取数据失败：", err)
		return
	}

	//log.Println("服务器端读取到的解密数据：", buf[:n])
	//log.Println(fmt.Sprintf("服务器端读取到的解密数据字符串：%s\n", string(buf[:n])))

	var method, URL, targetaddr string
	// 从解密的数据读入method，url
	_, err = fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:], '\n')]), "%s%s", &method, &URL)
	if err != nil {
		log.Println("解析失败：", err)
		return
	}
	hostPortURL, err := url.Parse(URL)
	if err != nil {
		log.Println("获取 method 和 url 失败：", err)
		return
	}

	var targetAddr *net.TCPAddr

	// 如果方法是CONNECT，则为https协议
	if method == "CONNECT" {
		targetaddr = hostPortURL.Scheme + ":" + hostPortURL.Opaque
		targetAddr, err = net.ResolveTCPAddr("tcp", targetaddr)
		if err != nil {
			log.Println("与目标建立连接失败(net.ResolveTCPAddr Error)：", err)
			return
		}
	} else { //否则为http协议
		// 如果host不带端口，则默认为80
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			targetaddr = hostPortURL.Host + ":80"
			targetAddr, err = net.ResolveTCPAddr("tcp", targetaddr)
			if err != nil {
				log.Println("与目标建立连接失败(net.ResolveTCPAddr Error)：", err)
				return
			}
		}
	}

	log.Println("目标URL：", targetAddr)

	//获得目标的host和port后，向目标服务器发起tcp连接
	targetConn, err := net.DialTCP("tcp", nil, targetAddr)
	if err != nil {
		log.Println(err)
		return
	} else {
		//响应连接成功
		//如果使用https协议，需先向客户端表示连接建立完毕
		if method == "CONNECT" {
			success := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
			lsServer.Cipher.encode(success) //响应客户端连接成功
			localConn.Write(success)
		} else { //如果使用http协议，需将从客户端得到的http请求转发给目标服务端
			targetConn.Write(buf[:n])
		}
	}

	// 进行转发
	// 从 localConn 读取数据发送到 targetConn
	go func() {
		err := localConn.DecodeCopy(targetConn)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			localConn.Close()
			targetConn.Close()
		}
	}()
	// 从 targetConn 读取数据发送到 localConn，这里因为处在翻墙阶段出现网络错误的概率更大
	(&SecureTCPConn{
		Cipher: localConn.Cipher,
		Conn:   targetConn,
	}).EncodeCopy(localConn)
}
