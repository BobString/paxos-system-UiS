package learner

// @authors Remy, Aliaksandr

import (
	//"fmt"
	"strings"
	"strconv"
)

// global variables
var(
	learnChan = make(chan string,5) // the channel of reception
	learnedValue string //
	pairMap = make(map[Pair] int)
	nbProc int
)

type Pair struct {
	nv int
	val string
}

func EntryPoint (count int) (chan string) {
	nbProc = count
	go receivingMsgs()
	return learnChan
}

func receivingMsgs () {	
	for {
	mesg := <- learnChan
	res := strings.Split(mesg,"@")
	a,_ := strconv.Atoi(res[1])
	p := Pair {a, res[2]}
	
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
