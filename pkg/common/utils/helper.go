package utils

import (
	"errors"
	"fmt"
	mathrand "math/rand"
	"net"
	"time"
)

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func RandomString(n, alphabet int, numeric bool) string {
	mathrand.Seed(time.Now().UnixNano())

	var strRand string
	var str string
	if alphabet > 0 {
		if numeric {
			str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		} else {
			str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			switch alphabet {
			case 1:
				str = "abcdefghijklmnopqrstuvwxyz"
			case 2:
				str = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			}
		}
	} else if numeric {
		for i := 0; i < n; i++ {
			strRand += fmt.Sprintf("%d", mathrand.Intn(10))
		}
		return strRand
	} else {
		str = "!@#$abcdefghijklmnopqrstuvwxyz!@#$ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$0123456789!@#$"
	}

	for i := 0; i < n; i++ {
		strRand += string(str[mathrand.Intn(len(str))])
	}
	return strRand
}
