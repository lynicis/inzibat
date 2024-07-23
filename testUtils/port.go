package testUtils

import "net"

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer tcpListener.Close()

	return tcpListener.Addr().(*net.TCPAddr).Port, nil
}
