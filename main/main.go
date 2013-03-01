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
	delay int = 1
)

var (
	process                    = map[int]string{}
	handlHBReplyChan           = make(chan int, 20)
	handlHBRequChan            = make(chan int, 20)
	handlTrustLeaderChan       = make(chan int, 20)
	handlRecoveryChan          = make(chan int, 20)
	handlSuspectChan           = make(chan int, 20)
	endChan                    = make(chan int, 5)
	handlLeaderReqChan         = make(chan int, 20)
	handlTrustChan             = make(chan int, 20)
	handlPromiseLeaderChan     = make(chan string, 20)
	inPrepChan                 = make(chan string, 20)
	inAcceptChan               = make(chan string, 20)
	learnChan                  = make(chan string, 20)
	ownProcess             int = 0
)

func main() {
	//Entry point of the application
	process = connector.GetProcesses()
	ownProcess, _ = connector.GetOwnProcess()
	//Create the connection
	go createServer()

	//Take the ids of all the processess
	keys := make([]int, len(process))
	i := 0
	for k, _ := range process {
		keys[i] = k
		i++
	}

	//Launch Leader Election	
	handlSuspectChan, handlRecoveryChan, handlTrustLeaderChan = leaderElection.EntryPoint(keys)
	//Launch Failure Detector
	handlHBReplyChan, handlHBRequChan = failureDetector.EntryPoint(delay, keys)
	//Call paxos and assign the channels

	<-endChan

}

func createServer() {
	fmt.Println("Starting server...")
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			continue
		}

		//Maintaining the connection per node
		go handleClient(conn)

	}

}

func handleClient(conn net.Conn) {
	for {
		leader := leaderElection.GetLeader()
		println("MAIN LEADER: ", leader)
		buf := make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			//If the client close the connection we get out and start listening again
			break
		}

		//res is where the message is going to be
		var res []string
		string1 := string(buf)
		res = strings.Split(string1, "@")
		stringaux := res[1]

		i, err := strconv.Atoi(stringaux)
		checkError(err)
		print("RECEIVED: ", res[0])
		println(" from ", i)
		switch {
		case res[0] == "Suspect":
			handlSuspectChan <- i
		case res[0] == "Restore":
			handlRecoveryChan <- i
		case res[0] == "HeartbeatReply":
			handlHBReplyChan <- i
		case res[0] == "HeartbeatRequest":
			handlHBRequChan <- i
		case res[0] == "LeaderRequest":
			handlTrustLeaderChan <- i
		case res[0] == "Promise":
			handlPromiseLeaderChan <- string1
		case res[0] == "Prepare":
			inPrepChan <- string1
		case res[0] == "Accept":
			inAcceptChan <- string1
		case res[0] == "Learn":
			learnChan <- string1
		case res[0] == "Trust":
			handlTrustChan <- i
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
