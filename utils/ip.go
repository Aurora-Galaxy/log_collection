package utils

import (
	"fmt"
	"net"
)

// GetOutboundIp 获取本机对外的ip
func GetOutboundIp() (ip string, err error) {
	//8.8.8.8:80 谷歌的dns地址
	// 只是播报，并没有发送请求
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String()) // 本机对外的ip及端口号
	//ip = strings.Split(localAddr.String(), ":")[0]
	ip = localAddr.IP.String()
	return
}
