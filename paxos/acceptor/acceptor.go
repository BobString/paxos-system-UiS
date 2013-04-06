package acceptor

//@author Remy Pannaud

import (
	"connector"
	"strconv"
	"strings"
	"paxos/slotsManager"
)

// global variables : 
var (
	learnList    []int  // list of known learners
	inPrepChan   = make(chan string, 5) // chan for reception of prepare 
	inAcceptChan = make(chan string, 5) // chan for reception of accept 
)

// functions :
// 	EntryPoint : init function
// 	prepareListener : prepare messages handler
// 	acceptListener : accept messages handler
//  preSend : to send messages easily

// Initialization function
//@parameters :
//	inPrepChan : when the server receives a Prepare message, it is sent to inPrepChan
//	inAcceptChan : idem with Accept messages
// 	sendRoundChan : to send to the Proposer the current system round number
func EntryPoint(list []int) (chan string, chan string) {
	learnList = list
	go readyListener ()
	return inPrepChan, inAcceptChan
}

func readyListener () {
	for {
		select {
		case prepare := <-inPrepChan:
			prepareHandler (prepare)
		case accept := <-inAcceptChan:
			acceptHandler (accept)
		}
	}
}


func prepareHandler(prepare string) {
		// wait for inPrepChan
		println("################### PREPARE HANDLER", prepare)
		s := strings.Split(prepare, "@")
		prepRN, _ := strconv.Atoi(s[1]) // get the int value for the if
		prepSender, _ := strconv.Atoi(s[3]) // get the int value for the preSend func
		slot, _ := strconv.Atoi(s[2]) // the slot number (instance of paxos)
		lvrn := slotsManager.GetLastVotedRN(slot)
		if prepRN > lvrn {
			lvval := slotsManager.GetLastVotedVal(slot)
			promise := "Promise@" + strconv.Itoa(prepRN) + "@" + strconv.Itoa(lvrn) + "@" + lvval + "@" + strconv.Itoa(slot)
			preSend(promise, prepSender)
		}
}

func acceptHandler(accept string) {
	// wait for inAcceptChan
	s := strings.Split(accept, "@")
	rN, _ := strconv.Atoi(s[1]) // get the int value for the if
	val := s[2] // the value
	slot, _ := strconv.Atoi(s[3]) // the slot number (instance of paxos)
	lvrn := slotsManager.GetLastVotedRN(slot)
	if rN >= lvrn {
		learn := "Learn@" + s[1] + "@" + s[2] + "@" + s[3]
		//lvrn = rN
		slotsManager.SetLastVotedRN(slot,rN)
		//sendRoundChan <- lvrn
		slotsManager.SetLastVotedVal(slot,val)
		for i := range learnList {
			preSend(learn, learnList[i])
		}
	}
	
}

func preSend(message string, pr int) {
	connector.Send(message, pr, nil)
}
