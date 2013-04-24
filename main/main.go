package main

import (
	"connector"
	"failureDetector"
	"fmt"
	"leaderElection"
	"net"
	"os"
	"paxos"
	"strconv"
	"strings"
	//"time"
	"accountManager"
)

const (
	delay int = 1
)

var (
	process                     = map[int]string{}
	handlHBReplyChan            = make(chan int, 50)
	handlHBRequChan             = make(chan int, 50)
	handlTrustLeaderChan        = make(chan int, 50)
	handlRecoveryChan           = make(chan int, 50)
	handlSuspectChan            = make(chan int, 50)
	endChan                     = make(chan int, 50)
	handlLeaderReqChan          = make(chan int, 50)
	handlTrustChan              = make(chan int, 50)
	handlPromiseLeaderChan      = make(chan string, 50)
	inPrepChan                  = make(chan string, 50)
	inAcceptChan                = make(chan string, 50)
	learnChan                   = make(chan string, 50)
	valueChan                   = make(chan string, 50)
	stopChan					= make(chan bool, 50)
	learnerToAccountManager = make(chan string, 50)
	askForProcChan chan int
	procMapChan chan string
	//debug 						= make(chan int, 50)
	ownProcess             int  = 0
	stopFlag				bool = false
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
	//Launch Account Manager
	learnerToAccountManager = accountManager.EntryPoint()
	//Call paxos and assign the channels
	handlTrustChan, inPrepChan, handlPromiseLeaderChan, inAcceptChan, learnChan, valueChan = paxosMain.EntryPoint(learnerToAccountManager)
	//Launch Leader Election	
	handlSuspectChan, handlRecoveryChan, handlTrustLeaderChan,askForProcChan,procMapChan = leaderElection.EntryPoint(keys, handlTrustChan)
	//Launch Failure Detector
	handlHBReplyChan, handlHBRequChan = failureDetector.EntryPoint(delay, keys)

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
		//leader := leaderElection.GetLeader()
		//println("MAIN LEADER: ", leader)
		buf := make([]byte, 4096)
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

		if stopFlag {
			if res[0] == "RestoreServer" {
				stopFlag = false
				connector.Stopped = false
			}
			continue

		}
		stringaux := res[1]
		if !strings.Contains(res[0],"Heartbeat"){
			//println("["+time.Now().String()+"]", "RECEIVED: ", res[0])
		}
		//println(" from ", i)
		switch res[0] {
		case "Suspect":
			i, err := strconv.Atoi(stringaux)
			checkError(err)
			handlSuspectChan <- i
		case "Restore":
			i, err := strconv.Atoi(stringaux)
			checkError(err)
			handlRecoveryChan <- i
		case "HeartbeatReply":
			i, err := strconv.Atoi(stringaux)
			checkError(err)
			handlHBReplyChan <- i
		case "HeartbeatRequest":
			i, err := strconv.Atoi(stringaux)
			checkError(err)
			handlHBRequChan <- i
		case "LeaderRequest":
			i, err := strconv.Atoi(stringaux)
			checkError(err)
			handlTrustLeaderChan <- i
		case "ProcMap":
		println("Receiving a MAP !!!!")
		 	procMapChan <- string1
		case "AskForProcMap":
		println("Receiving a ProcMap Ask !")
			i, err := strconv.Atoi(stringaux)
			checkError(err)
			askForProcChan <- i
		case "Promise":
			handlPromiseLeaderChan <- string1
		case "Prepare":
			inPrepChan <- string1
		case "Accept":
			inAcceptChan <- string1
		case "Learn":
			learnChan <- string1
		case "Value":
			lead := leaderElection.GetLeader()
			if lead == ownProcess {
				valueChan <- string1
			} else {
				connector.Send(string1, lead)
			}
		case "StopServer":
			stopFlag = true
			connector.Stopped = true
		//case "Debug":
			//debug <- 0
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
