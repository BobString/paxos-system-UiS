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
	learnedValue Pair // the learned value
	pairMap = make(map[Pair] int) // the map of pairs to decide the value
	nbProc int // the number of precesses, to determine if a quorum has sent Learn
)

type Pair struct {
	nv int
	val string
}
// Init function
func EntryPoint (count int) (chan string) {
	nbProc = count
	go receivingMsgs()
	return learnChan
}

func receivingMsgs () {	
	for {
		mesg := <- learnChan // wait for Learn message
		res := strings.Split(mesg,"@")
		a,_ := strconv.Atoi(res[1]) // to get the int for the round number
		p := Pair {a, res[2]} // creation of the pair to store in the map
		
		_,ok := pairMap [p] 
		// if the pair is in the map, we increase the count by 1, otherwise we put 1
		// IMPROVEMENT : create a global variable and use it to show the highest round number, and when receiving a higher round number, we delete the old ones
		if ok {
			pairMap[p] = pairMap[p]+1
		} else {
			pairMap[p] = 1
		}
		// we then checl if the a quorum of acceptors has sent the same Learn message
		if v,_ := pairMap[p]; v>(nbProc/2) {
			learnedValue = p
		}
	}
}
