package failureDetector

import (
	"connector"
	"net"
	"strconv"
	"time"
)

var (
	pAlive           = map[int]bool{}
	pSuspect         = map[int]bool{}
	pConnections     = map[int]*net.TCPConn{}
	emptyConnection  *net.TCPConn
	process          []int
	delay            int = 0
	ownProcess       int
	actualTimeout    int = 0
	handlHBReplyChan     = make(chan int, 20)
	handlHBRequChan      = make(chan int, 20)
)

func EntryPoint(d int, p []int) (chan int, chan int) {
	delay = d
	actualTimeout = d
	ownProcess, _ = connector.GetOwnProcess()
	process = p

	//Take the process and make 2 maps, alive and suspect
	for aux := range p {
		proc := process[aux]
		pAlive[proc] = true
		pSuspect[proc] = false

	}

	go startTimer(actualTimeout)
	return handlHBReplyChan, handlHBRequChan
}

func timeout(ticker *time.Ticker) bool {
	//First part to manage the delay
	var flag bool = false
	for aux := range process {
		pr := process[aux]
		if pSuspect[pr] && pAlive[pr] {
			//If there any project that is suspect, increase the timeout
			actualTimeout = actualTimeout + delay
			//We stop the last ticker
			ticker.Stop()
			flag = true
			break
		}
	}
	
	//Second part to manage the suspect process
	for aux := range process {
		pr := process[aux]
		if pr == ownProcess {
			continue
		}
		if !pAlive[pr] {
			pSuspect[pr] = true
			mess := "Suspect@" + strconv(pr)
			preSend(mess, pr)

		} else if pSuspect[pr] {
			pSuspect[pr] = false
			mess := "Restore@" + strconv(pr)
			preSend(mess, pr)
		}

		//Sent HeartbeatRequest to process
		preSend("HeartbeatRequest", pr)
	}

	//Put pAlive all to false
	for pr,_ := range pAlive {
		pAlive[pr] = false
	}

	//Start timer again only if the delay is changed
	if flag {
		startTimer(actualTimeout)

	}
	return flag
}

func gotHeartBeatRequest(pr int) {
	//Sent HeartbeatReply to pr
	preSend("HeartbeatReply", pr)

}

func gotHeartBeatReply(pr int) {
	pAlive[pr] = true
	//pSuspect[pr] = false
}

func startTimer(sec int) {
	//Time in seconds
	c := time.NewTicker(time.Duration(sec) * time.Second)
	var flag bool = false
	for {
		select {
		case pr := <-handlHBReplyChan:
			gotHeartBeatReply(pr)
		case pr := <-handlHBRequChan:
			gotHeartBeatRequest(pr)
		case <-c.C:
			flag = timeout(c)
		}
		if flag {
			break
		}
	}
}

func preSend(message string, pr int) {
	if message == "Suspect" || message == "Restore" {
		//If is Suspect or Restore we are going to send the message to ourself to our leader election 
		aux := message + "@" + strconv.Itoa(pr)
		connector.Send(aux, ownProcess)

	} else if message == "HeartbeatReply" || message == "HeartbeatRequest" {

		if pr == ownProcess {
			return
		}
		_, err := connector.Send(message, pr)
		if err != nil {
			pSuspect[pr] = true
			preSend("Suspect", pr)
		}
	}

}

//Presend function to keep in cache the connections
//func preSend(message string, pr int) {
//	if message == "Suspect" || message == "Restore" {
//		aux := message + "@" + strconv.Itoa(pr)
//		println("Suspected message " + aux)
//		connector.Send(aux, ownProcess, nil)
//	} else if pConnections[pr] == emptyConnection {
//		conn, err := connector.Send(message, pr, nil)
//		//If connection closed, error we have to treat the error.
//		if err != nil {
//			fmt.Printf("Error sending ", message, " to %d", pr)
//			pSuspect[pr] = true
//			preSend("Suspect", pr)
//		}
//		pConnections[pr] = conn
//	} else {
//		_, err := connector.Send(message, pr, pConnections[pr])
//		if err != nil {
//			fmt.Printf("Error sending ", message, " to %d", pr)
//			fmt.Printf("Trying to reconnect....")
//			conn, err := connector.Send(message, pr, nil)
//			if err != nil {
//				fmt.Printf("Error sending ", message, " to %d", pr)
//				pSuspect[pr] = true
//				preSend("Suspect", ownProcess)
//			}
//			pConnections[pr] = conn

//		}
//	}
//}
