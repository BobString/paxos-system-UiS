package learner
// @authors Remy, Aliaksandr
import (
    //"fmt"
    "paxos/slotsManager"
    "strconv"
    "strings"
    "time"
)
// global variables
var (
    learnChan = make(chan string, 5) // the channel of reception
    slotChan  chan int
    //learnedValue Pair // the learned value
    //pairMap = make(map[Pair] int) // the map of pairs to decide the value
    nbProc int // the number of precesses, to determine if a quorum has sent Learn
)
/*type Pair struct {
    Nv int
    Val string
}*/
// Init function
func EntryPoint(count int, newSlotChan chan int) chan string {
    nbProc = count
    slotChan = newSlotChan
    go learnListener()
    return learnChan
}
func clearMap(slot int) {
    /*for v := range pairMap {
        delete(pairMap,v)
    }*/
    slotsManager.ClearLearnMap(slot)
}
func learnListener () {
	for {
		learn := <-learnChan 
		go learnHandler(learn)
	}
}

func learnHandler(learn) {
	res := strings.Split(learn, "@")
    a, _ := strconv.Atoi(res[1])           // to get the int for the round number
    p := slotsManager.LearnPair{a, res[2]} // creation of the pair to store in the map
    slot, _ := strconv.Atoi(res[3])
    ok := slotsManager.BelongsToLearnMap(slot, p)
    // if the pair is in the map, we increase the count by 1, otherwise we put 1
    if ok {
 	   slotsManager.AddToLearnMap(slot, p, slotsManager.GetFromLearnMap(slot, p)+1)
	    //pairMap[p] = pairMap[p]+1
    } else {
		 slotsManager.AddToLearnMap(slot, p, 1)
		//pairMap[p] = 1
	}
    // we then check if the a quorum of acceptors has sent the same Learn message
    if v := slotsManager.GetFromLearnMap(slot, p); v > (nbProc / 2) {
        nextSlot := slotsManager.SetValueToLearn(slot, p.Val)
        //learnedValue = p
        //TODO : send to the proposer nextSlot, so he can send prepare !!!!!
        slotChan <- nextSlot
        println("["+time.Now().String()+"]", "NEW VALUE LEARNED :", p.Val, "in slot", res[3])
        clearMap(slot)
	}
}

