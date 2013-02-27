package acceptor
//@author Remy Pannaud


import (
	"connector"
)

// global variables : 
var (
	propList []string // list of known proposers
	learnList []string // list of known learners
	rnd int // current round number
	lvrn int // last voted round number
	lvval int // last voted calue
	sendPromiseChan = make(chan PromiseType,5) //intern channel
	sendLearnChan = make(chan LearnType,5)	// idem
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
		if v.RoundNum > lvren {
			promise := "Promise@"+v.RoundNum+"@"+lvrn+"@"+lvval
			preSend(prom,v.From)
		}
	}
}

func promiseSender (promise PromiseType, to int) {
	
}

func acceptListener() {
	for {
		v,_ := <-inAcceptChan
		
	}
}

func learnSender() {
}


func preSend(message string, pr int) {
	_, err := connector.Send(message, pr, nil)
	if err != nil {
		gotSuspectProc(pr)
	} else {
		gotProcRecov(pr)
	}
}