package testutil

import (
	"fmt"
	"net"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()

	if err != nil {
		t.Fatalf("GetFreePort returned an error: %v", err)
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port)))
	if err != nil {
		t.Fatalf("port %d returned by GetFreePort is not bindable: %v", port, err)
	}
	defer listener.Close()
}
