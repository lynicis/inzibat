package testUtils

import (
	"net"
	"strconv"
)

func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer tcpListener.Close()

	return strconv.Itoa(tcpListener.Addr().(*net.TCPAddr).Port), nil
}
