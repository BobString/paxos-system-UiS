package proposer

import (
	"connector"
	"fmt"
	"leaderElection"
	"strconv"
	"strings"
	"slotsManager"
	//"time"
)

// global variables
var (
	mv                     = map[int]string{} //Round-vote map 
	process                []int // list of processes
	leader                 int // currentLeader
	currentRound           int = 0 //our ID at first
	systemRound            int = 0 // current system round number
	handlTrustChan             = make(chan int, 20) // trust receive chan
	handlSysRoundChan          = make(chan int, 20) // system round number receive chan
	handlPromiseLeaderChan     = make(chan string, 20) // promise receive chan
	valueChan                  = make(chan string, 5) // value to decide chan
	maxRound               int = 0 // 
	cptProm				   int = 0
	
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
	currentRound = pickNext(currentRound)  //// FIXME : modify for slotsManager !!!!!!!!!!
	mv = map[int]string{}
	for pr := range process {
		proc := process[pr]
		cptProm = 0
		message := "Prepare@" + strconv.Itoa(currentRound)
		preSend(message, proc)
	}
}
// the round number increase function
func pickNext(currentRound int) int {
//FIXME : no need for systemRound.. just add 1 to slotsManager.GetRoundNumber ?
/*	for currentRound < systemRound {//// FIXME : modify for slotsManager !!!!!!!!!!
		currentRound = currentRound + len(process)//// FIXME : modify for slotsManager !!!!!!!!!!
	}*/
// FIXING
	currentRound := slotsManager.GetRoundNumber() +1// TO TEST
	slotsManager.SetRoundNumber(currentRound)
	return currentRound
}
// function handling the promise message
func gotPromise(data string) {
	// add variable cptProm to slotsManager ??
	cptProm = cptProm+1 //// FIXME : modify for slotsManager !!!!!!!!!!
	res := strings.Split(data, "@")
	roundnumber, _ := strconv.Atoi(res[1])
	lastVotedRound, _ := strconv.Atoi(res[2])
	lastVotedValue := res[3]
	slot = strconv.Atoi(res[4])
	//processID := res[4]
	if roundnumber == currentRound {//// FIXME : modify for slotsManager !!!!!!!!!!
		aux := lastVotedRound//// FIXME : modify for slotsManager !!!!!!!!!!
		mv[aux] = lastVotedValue//// FIXME : modify for slotsManager !!!!!!!!!!
		if aux > maxRound {
			maxRound = aux
		}
		if cptProm >= len(process)/2 {
			var valueToDecide string // the value we want to decide
			if maxRound == 0 {
				preValue := <-valueChan
				str := strings.Split(preValue, "@")
				valueToDecide = str[1]
				fmt.Println("Value to decide received :", valueToDecide, "in slot", res[4])
				//// FIXME : modify for slotsManager !!!!!!!!!!
			} else {
				//Pick the value form the largest round 
				valueToDecide = mv[maxRound]
			}
			curR := strconv.Itoa(currentRound)//// FIXME : modify for slotsManager !!!!!!!!!!
			sendAll("Accept@" + curR + "@" + valueToDecide)
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
			if sCR > systemRound {//// FIXME : modify for slotsManager !!!!!!!!!!
				systemRound = sCR
			}
		// TODO : add the case when we receive a slot number from the learner, so we have to send prepare messages !!
		}	
	}
}

func preSend(message string, pr int) {
	connector.Send(message, pr, nil)
}
