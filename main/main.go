package main

import (
	"connector"
	"failureDetector"
	"fmt"
	"leaderElection"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	delay int = 2
)

var (
	process              = map[int]string{}
	handlHBReplyChan     = make(chan int, 20)
	handlHBRequChan      = make(chan int, 20)
	handlTrustLeaderChan = make(chan int, 20)
	handlRecoveryChan    = make(chan int, 20)
	handlSuspectChan     = make(chan int, 20)
	endChan              = make(chan int, 5)

	ownProcess int = 0
)

func main() {
	//Entry point of the application
	process = connector.GetProcesses()
	ownProcess, _ = connector.GetOwnProcess()
	//Create the connection
	go createServer()

	//Launch Failure Detector		
	keys := make([]int, len(process))
	i := 0

	for k, _ := range process {
		keys[i] = k
		i++
	}
	handlSuspectChan, handlRecoveryChan, handlTrustLeaderChan = leaderElection.EntryPoint(keys)
	handlHBReplyChan, handlHBRequChan = failureDetector.EntryPoint(delay, keys)
	<-endChan

}

func createServer() {
	fmt.Println("Starting server...")
	//TODO: Make it available from outside
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listen(tcpAddr)

}

func listen(tcp *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", tcp)
	checkError(err)
	for {
		fmt.Println("Listening.......")
		conn, err := listener.Accept()
		if err != nil {

			continue

		}

		//Maintaining the connection per node
		go handleClient(conn)

	}

}

func handleClient(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			//If the client close the connection we get out and start listening again
			break
		}
		var res []string
		string1 := string(buf)

		res = strings.Split(string1, "@")
		stringaux := res[1]

		i, err := strconv.Atoi(stringaux)
		checkError(err)
		print("RECEIVED: ", res[0])
		println(" for ", i)

		switch {
		case res[0] == "Suspect":
			//TODO: Send to the leader election with a channel

			handlSuspectChan <- i
		case res[0] == "Restore":

			handlRecoveryChan <- i
		case res[0] == "HeartbeatReply":
			handlHBReplyChan <- i
		case res[0] == "HeartbeatRequest":
			handlHBRequChan <- i
		//case res[0] == "LeaderACK":
		case res[0] == "LeaderRequest":

			handlTrustLeaderChan <- i

		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
