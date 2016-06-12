package test

import (
	"fmt"
	"net"
	"strings"
	"sync/atomic"
)

var counter uint64

func mountString(host string, client string) string {
	toSlash := strings.Replace(host, "\\", "/", -1)
	winFix := strings.Replace(toSlash, "C:", "//c", 1)
	//winFix = strings.Replace(winFix, " ", "\\ ", -1)
	return winFix + ":" + client
}

func createID() string {
	i := atomic.AddUint64(&counter, 1)
	return fmt.Sprintf("cluster_%d", i)
}

// Address gets the local address
func Address() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}
