package proposer

import (
	"connector"
	"fmt"
	"leaderElection"
	"strconv"
	"strings"
	"time"
)

// global variables
var (
	mv                     = map[int]string{} //Round-vote map 
	process                []int // list of processes
	acceptors              []int // list of acceptors : USEFUL ??? //////
	leader                 int // currentLeader
	currentRound           int = 0 //our ID at first
	systemRound            int = 0 // current system round number
	handlTrustChan             = make(chan int, 20) // trust receive chan
	handlSysRoundChan          = make(chan int, 20) // system round number receive chan
	handlPromiseLeaderChan     = make(chan string, 20) // promise receive chan
	valueChan                  = make(chan string, 5) // value to decide chan
	maxRound               int = 0 // /////////////////////////////////////////// ??
	valueToDecide          string // the value we want to decide
)
// Init function
func EntryPoint(p []int, sysRoundChan chan int) (chan int, chan string, chan string) {
	process = p
	handlSysRoundChan = sysRoundChan
	currentRound, _ = connector.GetOwnProcess()
	go loop()
	return handlTrustChan, handlPromiseLeaderChan, valueChan
}
// what we do when we become the leader : change the round number
func gotTrust(leader int) {
	currentRound = pickNext(currentRound)
	mv = map[int]string{}
	for pr := range process {
		proc := process[pr]
		message := "Prepare@" + strconv.Itoa(currentRound)
		preSend(message, proc)
	}
}
// the round number increase function
func pickNext(currentRound int) int {

	for currentRound < systemRound {
		currentRound = currentRound + len(process)
	}
	return currentRound
}
// function handling the promise message
func gotPromise(data string) {
	res := strings.Split(data, "@")
	roundnumber, _ := strconv.Atoi(res[1])
	lastVotedRound, _ := strconv.Atoi(res[2])
	lastVotedValue := res[3]
	//processID := res[4]
	if roundnumber == currentRound {
		aux := lastVotedRound
		mv[aux] = lastVotedValue
		if aux > maxRound {
			maxRound = aux
		}
		if len(mv) >= len(process)/2 {
			var proposedValue string
			if aux == 0 {
				preValue := <-valueChan
				str := strings.Split(preValue, "@")
				valueToDecide = str[1]
				fmt.Println("Value to decide received :", valueToDecide)
			} else {
				//Pick the value form the largest round 
				proposedValue = mv[maxRound]
			}
			curR := strconv.Itoa(currentRound)
			sendAll("Accept@" + curR + "@" + proposedValue)
		}
	}

}
func sendAll(message string) {
	for pr := range process {
		proc := process[pr]
		preSend(message, proc)
	}
}
func isLeader() bool {
	me, _ := connector.GetOwnProcess()
	leader := leaderElection.GetLeader()
	return me == leader
}
func loop() {
	for {		
		select {
		case leader := <-handlTrustChan:
			gotTrust(leader)
		case data := <-handlPromiseLeaderChan:
			gotPromise(data)
		case sCR := <-handlSysRoundChan:
			if sCR > systemRound {
				systemRound = sCR
			}
		}	
	}
}

func preSend(message string, pr int) {
	connector.Send(message, pr, nil)
}
