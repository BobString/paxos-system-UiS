package learner

import (
	"fmt"
	"connector"
)

var(
	learnChan = make(chan string,5)
	learnedValue int
	pairMap = make(map[Pair] int)
)

type Pair struct {
	nv int
	val string
}

func EntryPoint () (chan string) {
	go receivingMsgs()
	return learnChan
}

func receivingMsgs ()  
{	
	procList := connector.GetProcesses()
	nbProc := len(procList)
	for {
	mesg := <- learnChan
	res = strings.Split(mesg,"@")
	
	p = Pair {res[1], res[2]}
	
	_,ok := pairMap [p]
	
	if ok {
		pairMap[p] = pairMap[p]+1
	} else {
		pairMap[p] = 1
	}
	
	if v,_ := pairMap[p]; v>(nbProc/2) {
		learnedValue = v 
	}
	}
}
