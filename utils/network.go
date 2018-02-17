package utils

import (
	"net"
	"strconv"
)

func String2TcpAddr(addr string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", addr)
}

func ToUdpAddr(ip string, port int) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", Concat(ip, strconv.Itoa(port)))
}

func ToTcpAddr(ip string, port int) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", Concat(ip, strconv.Itoa(port)))
}

func LocalTcpAddr(port string) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", Concat(":", port))
	if err != nil {
		return nil
	}
	return addr
}

func RandomLocalTcpAddr() *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", ":")
	if err != nil {
		return nil
	}
	return addr
}

func LocalUdpAddr(port string) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", Concat(":", port))
	if err != nil {
		return nil
	}
	return addr
}

func RandomLocalUdpAddr() *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", ":")
	if err != nil {
		return nil
	}
	return addr
}
