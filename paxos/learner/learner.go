package learner

// @authors Remy, Aliaksandr

import (
	//"fmt"
	"strings"
	"strconv"
	"time"
	"slotsManager"
)

// global variables
var(
	learnChan = make(chan string,5) // the channel of reception
	//learnedValue Pair // the learned value
	//pairMap = make(map[Pair] int) // the map of pairs to decide the value
	nbProc int // the number of precesses, to determine if a quorum has sent Learn
)

/*type Pair struct {
	Nv int
	Val string
}*/
// Init function
func EntryPoint (count int) (chan string) {
	nbProc = count
	go receivingMsgs()
	return learnChan
}

func clearMap(slot int) {
	/*for v := range pairMap {
		delete(pairMap,v)
	}*/
	slotsManager.ClearLearnMap(slot)
}


func receivingMsgs () {	
	for {
		mesg := <- learnChan // wait for Learn message
		res := strings.Split(mesg,"@")
		a,_ := strconv.Atoi(res[1]) // to get the int for the round number
		p := slotsManager.LearnPair {a, res[2]} // creation of the pair to store in the map
		slot := res[3]
		ok := slotsManager.BelongToLearnMap(slot, p)
		// if the pair is in the map, we increase the count by 1, otherwise we put 1
		if ok {
			slotsManager.AddToLearnMap(slot,p,slotsManager.GetFromLearnMap(slot,p) + 1)
			//pairMap[p] = pairMap[p]+1
		} else {
			slotsManager.AddToLearnMap(slot,p,1)
			//pairMap[p] = 1
		}
		// we then check if the a quorum of acceptors has sent the same Learn message
		if v := slotsManager.GetFromLearnMap(slot,p); v>(nbProc/2) {
			nextSlot := slotsManager.SetValueToLearn(slotsManager.GetValueFromLearnPair(p))
			//learnedValue = p
			//TODO : send to the proposer nextSlot, so he can send prepare !!!!!
			println("["+time.Now().String()+"]","NEW VALUE LEARNED :",p.Val,"in slot",res[3])
			clearMap(slot)
		}
	}
}
