package leaderElection

import (
	"connector"
	"net"
	"time"
)

var (
	pSuspect                 = map[int]bool{}
	leader               int = 0
	handlSuspectChan         = make(chan int, 20)
	handlRecoveryChan        = make(chan int, 20)
	handlTrustLeaderChan     = make(chan int, 20)
	handlLeaderReqChan       = make(chan int, 20)
	process              []int
	ownProcess           int
	pConnections         = map[int]*net.TCPConn{}
	emptyConnection      *net.TCPConn
	acks                 int = 0
	timeLeaderRequest    int = 2
)

func EntryPoint(p []int) (chan int, chan int, chan int) {
	process = p
	ownProcess, _ = connector.GetOwnProcess()

	//Take the process and make one map, suspect
	for pr := range process {
		proc := process[pr]
		pSuspect[proc] = false
	}
	leader = ownProcess
	go auxiliar()

	return handlSuspectChan, handlRecoveryChan, handlTrustLeaderChan

}

func GetLeader () {
	return leader
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
		case <-c.C:
			if leader == ownProcess {
				newLeaderRequest()
			}
		}
		println("Leader: ", leader)
	}
}

func trustLeader(pr int) {
	switch {
	case pr == leader:
		//We do nothing, we already now that this is a the leader
	case pr > leader:
		//We have a better leader, if we are the leader, we tell this confused process
		if leader == ownProcess {
			newLeaderRequest()
			//sendLeaderRquest(pr)
			println("LEADER ELECTION: Process confused sending candidature")
		}
	case pr < leader:
		leader = pr
		println("LEADER ELECTION: Accepting leader ", pr)
	}

}

func newLeaderRequest() {
	//Send Leader request to everyone
	for pr := range process {
		proc := process[pr]
		preSend("LeaderRequest", proc)
	}
}

func sendLeaderRquest(pr int) {
	//Send leader request to process pr
	preSend("LeaderRequest", pr)
}

func gotSuspectProc(pr int) {
	pSuspect[pr] = true
	if pr == leader {
		leader = ownProcess
		println("Leader: ", leader)
		newLeaderRequest()
	}
}
func gotProcRecov(pr int) {
	pSuspect[pr] = false
}

func preSend(message string, pr int) {
	_, err := connector.Send(message, pr, nil)
	if err != nil {
		gotSuspectProc(pr)
	} else {
		gotProcRecov(pr)
	}
}
