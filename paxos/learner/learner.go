package learner
// @authors Remy, Aliaksandr
import (
    //"fmt"
    "paxos/slotsManager"
    "strconv"
    "strings"
    "time"
	//"connector"
)

// global variables
var (
    learnChan = make(chan string, 50) // the channel of reception
    slotChan  chan int
    nbProc int // the number of precesses, to determine if a quorum has sent Learn
	toAccManager chan string
)

// Init function
func EntryPoint(count int, newSlotChan chan int, toAccountManager chan string) chan string {
    nbProc = count
    slotChan = newSlotChan
	toAccManager = toAccountManager
    go learnListener()
    return learnChan
}

func clearMap(slot int) {
    slotsManager.ClearLearnMap(slot)
}

func learnListener () {
	for {
		learn := <-learnChan 
		learnHandler(learn)
	}
}

func learnHandler(learn string) {
	res := strings.Split(learn, "@")
    a, _ := strconv.Atoi(res[1])           // to get the int for the round number
    p := slotsManager.LearnPair{a, res[2]} // creation of the pair to store in the map
    slot, _ := strconv.Atoi(res[3])
 	count := slotsManager.GetFromLearnMap(slot, p)
	println("learn count BEFORE ", strconv.Itoa(count+1)
  	slotsManager.AddToLearnMap(slot, p, count+1)
	println("learn count AFTER (for test)", strconv.Itoa(count+1)
    // we then check if the a quorum of acceptors has sent the same Learn message
	v := slotsManager.GetFromLearnMap(slot, p)
    if (v > (nbProc / 2)) /*&& (slotsManager.GetValueLearned(slot)=="") */{
        nextSlot := slotsManager.SetValueToLearn(slot, p.Val)
        slotChan <- nextSlot
        println("["+time.Now().String()+"]", "NEW VALUE LEARNED :", p.Val, "in slot", res[3])
		toAccManager <- p.Val
	}
}

