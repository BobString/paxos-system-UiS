package paxosMain
//@author Remy Pannaud


import (
	"paxos/acceptor"
	"paxos/learner"
	"paxos/proposer"
	"strings"
)

// global variables :
var (
	acceptToPropChan = make(chan int,5)
)
 
// Initialization function
//@parameters :
func EntryPoint () (chan int, chan string, chan string, chan string, chan string) {
	// TODO : get the []int of processes
	
	//
	trustChan,promChan := paxos/proposer.EntryPoint(p,acceptToPropChan)
	prepChan,acceptChan := paxos/acceptor.EntryPoint(acceptToPropChan)
	learnChan =:= paxos/learner.EntryPoint()	
	return trustChan,prepChan,promChan,acceptChan,learnChan
}

