package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetIp() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		os.Exit(1)
	}
	for _, iface := range interfaces {
		// 跳过无效的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取该接口的所有地址
		addresses, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error getting addresses:", err)
			continue
		}

		// 遍历每个地址，检查是否是 IPv4 地址
		for _, addr := range addresses {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}

			// 检查地址是否是以 "192" 开头
			if ipNet.IP.To4() != nil && strings.HasPrefix(ipNet.IP.String(), "192.") {
				return ipNet.IP.String()
			}
		}
	}
	return "localhost"
}
