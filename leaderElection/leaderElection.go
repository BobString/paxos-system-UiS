package leaderElection

import (
	"connector"
	"net"
	"time"
	"strings"
	"strconv"
)

var (
	pSuspect                 = map[int]bool{}
	leader               int = 0
	//countLeader          int = 0
	handlSuspectChan         = make(chan int, 20)
	handlRecoveryChan        = make(chan int, 20)
	handlTrustLeaderChan     = make(chan int, 20)
	handlLeaderReqChan       = make(chan int, 20)
	askForProcMapChan  = make(chan int, 20)
	procMapChan  = make(chan string, 20)
	handlTrustChan       chan int
	process              []int
	ownProcess           int
	pConnections         = map[int]*net.TCPConn{}
	emptyConnection      *net.TCPConn
	acks                 int = 0
	timeLeaderRequest    int = 1
	processMap = make (map[int]int)
)

func EntryPoint(p []int, trust chan int) (chan int, chan int, chan int, chan int, chan int) {
	process = p
	ownProcess, _ = connector.GetOwnProcess()
	handlTrustChan = trust
	//Take the process and make one map, suspect
	for pr := range process {
		proc := process[pr]
		pSuspect[proc] = false
	}
	//leader = ownProcess
	go auxiliar()

	return handlSuspectChan, handlRecoveryChan, handlTrustLeaderChan, askForProcMapChan, procMapChan

}

func GetLeader() int {
	return leader
}

func initProcMap() map[int]int {
	// send the AskForProcMap to everyone
	for pr := range process {
		preSend("AskForProcMap",process[pr])
	}
	// if within an amount of time we did not receive a ProcMap message, then we consider ourself as the leader, and create the map
	//start timer
	timer := time.NewTicker(time.Duration(100) * time.Millisecond)
	select {
	case <-timer.C:	
		processMap[ownProcess] = 0
		leader = ownProcess
		println("Init leader : we are first")
		//creation of map
	case mess:=<-procMapChan
		// decrypting message
		aux := strings.Split(mess,"@")	
		for i:=1;i<len(aux)-2;i++ {
			p,_ := strconv.Atoi(aux[i])
			c,_ := strconv.Atoi(aux[i+1])
			if c==0 {
				leader = c
				println("Init leader : "+strconv.Itoa(leader)+" is leader")
			}
			processMap[p]=c
			i = i+1 // that way each time i is increased of 2			
		}
	}	
}

func auxiliar() {
	//newLeaderRequest()
	c := time.NewTicker(time.Duration(timeLeaderRequest) * time.Second)
	for {

		select {
		case pr := <-handlSuspectChan:
			gotSuspectProc(pr)
		case pr := <-handlRecoveryChan:
			gotProcRecov(pr)
		case pr := <-handlTrustLeaderChan:
			trustLeader(pr)
		case pr := <-askForProcMapChan :
			if leader==ownProcess {
				encryptMap(pr)
			}
		case <-c.C:
			if leader == ownProcess {
				//countLeader = countLeader + 1
				//newLeaderRequest()
				//if countLeader == 2 {
					//We are a stable leader
					handlTrustChan <- ownProcess
					//println("["+time.Now().String()+"]","We are the leader !!")
				//	countLeader = countLeader + 1
				//}
			//} else {
			//	countLeader = 0
			//}
			}
		//println("Leader: ", leader)
		}
	}
}

func encryptMap (pr int) {
	mess := "ProcMap"
	for p := range processMap {
		mess = mess + "@" + strconv.Itoa(p) + "@" + strconv.Itoa(processMap[p])
	}
	preSend(mess,pr)
}

func trustLeader(pr int) {
	switch {
	case pr == leader:
		//We do nothing, we already now that this is a the leader
	case processMap[pr] < processMap[leader]:
		//If we are the leader, we tell this confused process
		if leader == ownProcess {
			//sendLeaderRquest(pr)
			println("LEADER ELECTION: Process confused sending candidature")
		}
	case processMap[pr] > processMap[leader]:
		// if we are the leader, we say : hey no bro i'm the leader
		newLeaderRequest()
		//println("LEADER ELECTION: Accepting leader ", pr)
	}

}

func newLeaderRequest() {
	//Send Leader request to everyone
	for pr := range process {
		proc := process[pr]
		preSend("LeaderRequest", proc)
	}
}

/*func sendLeaderRquest(pr int) {
	//Send leader request to process pr
	preSend("LeaderRequest", pr)
}*/

func gotSuspectProc(pr int) {
	pSuspect[pr] = true
	old = processMap[pr]
	delete(processMap,pr)
	for p := range processMap {
		if processMap[pr]> old {
			processMap[pr] = processMap[pr] - 1
			if processMap[pr]==0 {
				leader = pr
			}
		}
	}
	println(strconv.Itoa(pr)+" crashed. Leader is "+strconv.Itoa(leader))
	/*if pr == leader {
		leader = ownProcess
		//println("Leader: ", leader)
		newLeaderRequest()
	}*/
}
func gotProcRecov(pr int) {
	pSuspect[pr] = false
	max := 0
	for p := range processMap{
		if processMap[p] > max{
			max = processMap[p]
		}
	}
	processMap[pr] = max + 1	
}

func preSend(message string, pr int) {
	_, err := connector.Send(message, pr)
	if err != nil {
		gotSuspectProc(pr)
	} else {
		gotProcRecov(pr)
	}
}
