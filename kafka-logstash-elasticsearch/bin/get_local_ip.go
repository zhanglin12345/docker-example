package main

import (
	"fmt"
	"net"
	"os"
)

func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	local_ip := ""

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// check if the ip address is loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				local_ip = ipnet.IP.String()
				break
			}

		}
	}
	return local_ip
}


func main() {
	local_ip := GetIntranetIp()
	fmt.Printf("%s", local_ip)
}
