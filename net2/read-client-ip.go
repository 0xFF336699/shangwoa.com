package net2

import (
	"net"
	"net/http"
	"strings"
)

func GetClientPublicIP(r *http.Request) string {
	var ip string
	//fmt.Println("x forwarded for is", r.Header.Get("X-Forwarded-For"))
	//fmt.Println("x real ip is", r.Header.Get("X-Real-Ip"))
	//fmt.Println("remote addr is", r.RemoteAddr)
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}

	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPddr(ip) {
			return ip
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPddr(ip) {
			return ip
		}
	}

	return ""
}

func HasLocalIPddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {

	//for _, network := range localNetworks {
	//	if network.Contains(ip) {
	//		return true
	//	}
	//}

	return ip.IsLoopback()
}

