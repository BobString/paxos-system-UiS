package connector

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	//"time"
	//"errors"
)

var (
	process = map[int]string{
		1: "152.94.0.120:1200", //pitter28
		2: "152.94.0.117:1200", // pitter25
		//2: "152.94.0.121:1200", //pitter29		
		3: "152.94.0.114:1200", //pitter22
		//3: "152.94.0.124:1200", //pitter32
		//3: "152.94.0.115:1200", //pitter23
		
	}
	ownProcess int    = 0
	ownIP      string = ""
	Stopped bool = false
	connMap = make (map[int] *net.TCPConn)
)

func Send(message string, pr int) (*net.TCPConn, error) {
	var err error
	var connect *net.TCPConn
	err = nil
	connect = nil
	if !Stopped {
		if _,ok := connMap[pr]; !ok {
			service := process[pr]
			tcpAddr, err := net.ResolveTCPAddr("tcp", service)
			if err != nil {
				return nil, err
			}
			connect, err = net.DialTCP("tcp", nil, tcpAddr)
			connMap[pr] = connect
			if err != nil {
				return nil, err
			}
		} else {
			connect = connMap[pr]
		}
		ownProcess := 0
		if !strings.Contains(message,"Value") {
			ownProcess, _ = GetOwnProcess()
		} else {
			//println(connect.LocalAddr().String())
			addr := connect.LocalAddr().String()
			aux := strings.Split(addr,":")
			addr = aux[0]+":1201"
			message = message + addr + ","
		}
		_, err := connect.Write([]byte(message + "@" + strconv.Itoa(ownProcess) + "@"))
		if err != nil {
			delete(connMap,pr)
			//println("Error dialing the TCP addrs")
			return nil, err
		}
	}
	return connect, err
}

func SendByAddr(message string, remAddr string) (error) {
	var err error
	var connect *net.TCPConn
	err = nil
	connect = nil
	if !Stopped {
		tcpAddr, err := net.ResolveTCPAddr("tcp", remAddr)
		if err != nil {
			return err
		}
		connect, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return err
		}
		ownProcess := 0
		if !strings.Contains(message,"Value") {
			ownProcess, _ = GetOwnProcess()
		} 
		_, err = connect.Write([]byte(message + "@" + strconv.Itoa(ownProcess) + "@"))
		if err != nil {
			//println("Error dialing the TCP addrs")
			return err
		}
		connect.Close()
	}
	return err
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
		fmt.Fprintf(os.Stderr, "Fatal error connector ", err.Error())
		os.Exit(1)
	}
}
