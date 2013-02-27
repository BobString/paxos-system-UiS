package acceptor
//@author Remy Pannaud


import (
	"paxos/signals"
)

// global variables : 
var (
	propList []string // list of known proposers
	learnList []string // list of known learners
	rnd int // current round number
	lvrn int // last voted round number
	lvval int // last voted calue
	sendPromiseChan = make(chan string,5) //intern channel
	sendLearnChan = make(chan string,5)	// idem
)

// functions :
// 	init
// 	prepareListener
//	promiseSender
// 	acceptListener
// 	learnSender

func init (inPrepChan chan PrepareType, inAcceptChan chan AcceptType) {
	go prepareListener(inPrepChan)
	go acceptListener(inAcceptChan)
}


func prepareListener () {
	for {
		v,_ := <-inPrepChan
		if v.>lvren {
			
		}
	}
}

func promiseSender () {
}

func acceptListener() {
}

func learnSender() {
}