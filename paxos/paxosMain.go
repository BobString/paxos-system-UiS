package paxosMain
//@author Remy Pannaud


import (
	"connector"
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
func EntryPoint () (chan int, chan string, chan string, chan string, chan string, chan string) {
	// TODO : get the []int of processes
	pMap := connector.GetProcesses()
	var  p = make([]int,len(pMap))
	i := 0
	for v,_ := range pMap {
		p[i] = v 
		i++
	}
	//
	trustChan,promChan,valueChan := paxos/proposer.EntryPoint(p,acceptToPropChan)
	prepChan,acceptChan := paxos/acceptor.EntryPoint(p,acceptToPropChan)
	learnChan := paxos/learner.EntryPoint(len(p))	
	return trustChan,prepChan,promChan,acceptChan,learnChan,valueChan
}

