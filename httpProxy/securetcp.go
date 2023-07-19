package httpProxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	bufSize = 1024
)

// SecureTCPConn 加密传输的 TCP Socket
type SecureTCPConn struct {
	net.Conn
	Cipher *Cipher
}

// DialTCPSecure see net.DialTCP
func DialTCPSecure(serverAddr *net.TCPAddr, cipher *Cipher) (*SecureTCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return nil, err
	}
	return &SecureTCPConn{
		Conn:   remoteConn,
		Cipher: cipher,
	}, nil
}

// ListenSecureTCP 建立安全的TCP隧道
func ListenSecureTCP(listenaddr *net.TCPAddr, handleConn func(localConn *SecureTCPConn), cipher *Cipher, didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", listenaddr)
	if err != nil {
		return errors.New(fmt.Sprintf("net.ListenTCP Error：%s\n", err))
	}
	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("与 %s 建立连接失败：%s\n", localConn.RemoteAddr(), err)
			continue
		}
		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		go handleConn(&SecureTCPConn{
			Conn:   localConn,
			Cipher: cipher,
		})
	}
}

// DecodeRead 从输入流里读取加密过的数据，解密后把原数据放到bs里
func (secureTCP *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureTCP.Read(bs)
	if err != nil {
		return
	}

	secureTCP.Cipher.decode(bs[:n])
	return
}

// DecodeCopy 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (secureTCP *SecureTCPConn) DecodeCopy(dst io.Writer) error {
	buf := make([]byte, bufSize)
	for {
		readCount, readErr := secureTCP.DecodeRead(buf)
		if readErr != nil {
			if readErr != io.EOF {
				return readErr
			} else {
				return nil
			}
		}

		if readCount > 0 {
			writeCount, writeErr := dst.Write(buf[0:readCount])
			if writeErr != nil {
				return writeErr
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// EncodeWrite 把放在bs里的数据加密后立即全部写入输出流
func (secureTCP *SecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	secureTCP.Cipher.encode(bs)
	return secureTCP.Write(bs)
}

// EncodeCopy 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (secureTCP *SecureTCPConn) EncodeCopy(dst net.Conn) error {
	buf := make([]byte, bufSize)
	for {
		readCount, readErr := secureTCP.Read(buf)
		if readErr != nil {
			if readErr != io.EOF {
				return readErr
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, writeErr := (&SecureTCPConn{
				Conn:   dst,
				Cipher: secureTCP.Cipher,
			}).EncodeWrite(buf[0:readCount])
			if writeErr != nil {
				return writeErr
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}
