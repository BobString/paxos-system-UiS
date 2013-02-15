package leaderElection

import (
	"connector"

	"net"
)

var (
	pSuspect                 = map[int]bool{}
	leader               int = 0
	handlSuspectChan         = make(chan int, 20)
	handlRecoveryChan        = make(chan int, 20)
	handlNewLeaderChan       = make(chan int, 20)
	handlTrustLeaderChan     = make(chan int, 20)
	process              []int
	ownProcess           int
	pConnections         = map[int]*net.TCPConn{}
	emptyConnection      *net.TCPConn
	acks                 int = 0
)

func EntryPoint(p []int) (chan int, chan int, chan int) {
	process = p
	ownProcess, _ = connector.GetOwnProcess()
	//Take the process and make 2 maps, alive and suspect
	for pr := range process {
		proc := process[pr]
		pSuspect[proc] = false
	}

	go auxiliar()

	return handlSuspectChan, handlRecoveryChan, handlTrustLeaderChan

}

func restore(proc int) {
	pSuspect[proc] = false
}

func auxiliar() {
	for {

		if leader == 0 {
			newLeaderRequest()
		}
		select {
		case pr := <-handlSuspectChan:
			gotSuspectProc(pr)
		case pr := <-handlRecoveryChan:
			gotProcRecov(pr)
		case pr := <-handlTrustLeaderChan:
			trustLeader(pr)
		case <-handlNewLeaderChan:
			newLeaderRequest()
		}
		println("Leader: ", leader)
	}
}
func trustLeader(pr int) {

	//A message from other node came with the leader proposal
	if (leader < pr || pSuspect[pr] || pr == 0 || pr > ownProcess) && leader != 0 {
		//The new leader is bigger or is suspect so I am the leader
		println("I want to be the leader")
		newLeaderRequest()

	} else if (leader == 0 || leader > pr) && !pSuspect[pr] && pr != 0 && pr < ownProcess {
		//Accept leader,
		leader = pr
		println("Accepting leader ", pr)

		//preSend("LeaderACK", pr)
	} else if (leader == 0 || leader > pr) && !pSuspect[pr] && pr != 0 && pr >= ownProcess {
		//We are better leader or we sent the message to ourself
		leader = ownProcess
		println("I am the leader")
		//acks = 0
	}
}

//func leaderACKs() {
//	acks++
//	if acks > (len(process)/2)+1 {
//		leader = ownProcess
//		acks = 0
//	}
//}

func newLeaderRequest() {
	//Send leaderProposal and wait until we have more than 50% ACKs  and then
	for pr := range process {

		proc := process[pr]

		preSend("LeaderRequest", proc)

	}
}

func gotSuspectProc(pr int) {
	if pr != leader {
		pSuspect[pr] = true
	} else {
		leader = 0
		handlNewLeaderChan <- 0
	}
}
func gotProcRecov(pr int) {
	pSuspect[pr] = false
}

func getLeader() int {
	if leader != 0 {
		return leader
	}
	return -1

}

func preSend(message string, pr int) {

	println("LEADER ELECTION: "+message+" message, to the process: ", pr)
	connector.Send(message, pr, nil)

}
