package main

import (
	"log"
	"net"

	"golang.org/x/net/proxy"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	// 创建SOCKS5代理
	server, err := proxy.SOCKS5("tcp", "target-server:1080", nil, proxy.Direct)
	if err != nil {
		log.Println(err)
		return
	}

	// 连接目标服务器
	targetConn, err := server.Dial("tcp", "example.com:80")
	if err != nil {
		log.Println(err)
		return
	}
	defer targetConn.Close()

	// 在两者之间进行数据转发
	go func() {
		_, err := io.Copy(clientConn, targetConn)
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(targetConn, clientConn)
	if err != nil {
		log.Println(err)
	}
}
