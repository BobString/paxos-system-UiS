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
func EntryPoint(list []int, sendRoundChan chan int) (chan string, chan string) {
	learnList = list
	go prepareListener(inPrepChan)
	go acceptListener(inAcceptChan, sendRoundChan)
	return inPrepChan, inAcceptChan
}

func prepareListener(inPrepChan chan string) {
	for {
		// wait for inPrepChan
		v := <-inPrepChan
		s := strings.Split(v, "@")
		v1, _ := strconv.Atoi(s[1]) // get the int value for the if
		v2, _ := strconv.Atoi(s[2]) // get the int value for the preSend func
		v3, _ := strconv.Atoi(s[3]) // the slot number (instance of paxos)
		lvrn := slotsManager.GetLastVotedRN(v3)
		if v1 > lvrn {
			lvval := slotsManager.GetLastVotedVal(v3)
			promise := "Promise@" + strconv.Itoa(v1) + "@" + strconv.Itoa(lvrn) + "@" + lvval + "@" + strconv.Itoa(v3)
			preSend(promise, v2)
		}
	}
}

func acceptListener(inAcceptChan chan string, sendRoundChan chan int) {
	for {
		// wait for inAcceptChan
		v := <-inAcceptChan
		s := strings.Split(v, "@")
		s1, _ := strconv.Atoi(s[1]) // get the int value for the if
		s3, _ := strconv.Atoi(s[3]) // the slot number (instance of paxos)
		lvrn := slotsManager.GetLastVotedRN(s3)
		if s1 >= lvrn {
			learn := "Learn@" + s[1] + "@" + s[2] + "@" + strconv.Itoa(s3)
			lvrn = s1
			slotsManager.SetLastVotedRN(s3,lvrn)
			//sendRoundChan <- lvrn
			lvval := s[2]
			slotsManager.SetLastVotedVal(s3,lvval)
			for i := range learnList {
				preSend(learn, learnList[i])
			}
		}
	}
}

func preSend(message string, pr int) {
	connector.Send(message, pr, nil)
}
