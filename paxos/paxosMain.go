package paxosMain

//@author Remy Pannaud

import (
	"connector"
	"paxos/acceptor"
	"paxos/learner"
	"paxos/proposer"
)

// global variables :
var (
	acceptToPropChan = make(chan int, 5)
)

// Initialization function
//@parameters :
func EntryPoint() (chan int, chan string, chan string, chan string, chan string, chan string) {
	pMap := connector.GetProcesses()
	var p = make([]int, len(pMap))
	i := 0
	for v, _ := range pMap {
		p[i] = v
		i++
	}
	//
	trustChan, promChan, valueChan, slotChan := proposer.EntryPoint(p, acceptToPropChan)
	prepChan, acceptChan := acceptor.EntryPoint(p, acceptToPropChan)
	learnChan := learner.EntryPoint(len(p),slotChan)
	return trustChan, prepChan, promChan, acceptChan, learnChan, valueChan
}
