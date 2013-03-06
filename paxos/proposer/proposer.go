package proposer

import (
	"connector"
	"fmt"
	"leaderElection"
	"strconv"
	"strings"
)

var (
	mv                     = map[int]string{} //Round-vote map 
	process                []int
	acceptors              []int
	leader                 int
	currentRound           int = 0 //our ID at first
	systemRound            int = 0
	handlTrustChan             = make(chan int, 20)
	handlSysRoundChan          = make(chan int, 20)
	handlPromiseLeaderChan     = make(chan string, 20)
	valueChan                  = make(chan string, 5)
	maxRound               int = 0
	valueToDecide          string
)

func EntryPoint(p []int, sysRoundChan chan int) (chan int, chan string, chan string) {
	process = p
	handlSysRoundChan = sysRoundChan
	currentRound, _ = connector.GetOwnProcess()
	go listenToValue()
	go loop()
	return handlTrustChan, handlPromiseLeaderChan, valueChan
}

func listenToValue() {
	for {
		preValue := <-valueChan
		str := strings.Split(preValue, "@")
		valueToDecide = str[1]
		fmt.Println("Value to decide received :", valueToDecide)
	}
}

func gotTrust(leader int) {
	currentRound = pickNext(currentRound)
	mv = map[int]string{}
	for pr := range process {
		proc := process[pr]
		message := "Prepare@" + strconv.Itoa(currentRound)
		preSend(message, proc)
	}

}
func pickNext(currentRound int) int {

	for currentRound < systemRound {
		currentRound = currentRound + len(process)
	}
	return currentRound
}

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
				//FIXME: ASK SOMEONE TO ENTER THE VALUE
				proposedValue = valueToDecide
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
		//FIXME: All processes are going inside because at the first time all are leaders ////////DONE
		if isLeader() {
			select {
			case leader := <-handlTrustChan:
				gotTrust(leader)
			case data := <-handlPromiseLeaderChan:
				gotPromise(data)
			case sCR := <-handlSysRoundChan:
				if sCR > systemRound {
					systemRound = sCR
				}
			default: //TODO : check if works correctly
			}
		}
	}
}

func preSend(message string, pr int) {
	connector.Send(message, pr, nil)
}
