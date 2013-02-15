package connector

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	process = map[int]string{
		1: "152.94.0.116:1200",
		2: "152.94.0.115:1200",
		3: "152.94.0.114:1200",
		//4: "111.111.11.1:1200",
		//5: "111.111.11.1:1200",
	}
	ownProcess int    = 0
	ownIP      string = ""
)

func Send(message string, pr int, connect *net.TCPConn) (*net.TCPConn, error) {
	if connect == nil {
		service := process[pr]
		tcpAddr, err := net.ResolveTCPAddr("tcp", service)
		if err != nil {
			println("Error resolving the TCP addrs")
			return nil, err
		}
		connect, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			println("Error dialing the TCP addrs")
			return nil, err
		}
	}

	//print("SEND: ", message)
	//println(" to ", pr)
	if message == "HeartbeatRequest" || message == "HeartbeatReply" || message == "LeaderRequest" {
		ownProcess, _ := GetOwnProcess()
		_, err := connect.Write([]byte(message + "@" + strconv.Itoa(ownProcess) + "@"))
		return connect, err
	}
	_, err := connect.Write([]byte(message + "@" + strconv.Itoa(pr) + "@"))

	return connect, err

}
func GetProcesses() map[int]string {
	return process
}

func GetOwnProcess() (int, string) {
	if ownProcess != 0 {
		return ownProcess, ownIP
	}

	ownIP = getLocalIp()

	for k, ip := range process {
		var res []string
		res = strings.Split(ip, ":")
		if res[0] == ownIP {
			ownProcess = k
			ownIP = ip
			break
		}

	}
	if ownProcess == 0 {
		fmt.Println("Fatal error, not own process found ")
		os.Exit(1)
	}
	return ownProcess, ownIP
}

func getLocalIp() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return ""
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return ""
	}

	return addrs[0]

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
