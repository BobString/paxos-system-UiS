
package proposer
import (
    "connector"
    "fmt"
    "leaderElection"
    "paxos/slotsManager"
    "strconv"
    "strings"
    "time"
)
// global variables
var (
    process []int // list of processes
    leader  int   // currentLeader
    handlTrustChan = make(chan int, 50) // trust receive chan
    handlPromiseLeaderChan = make(chan string, 50) // promise receive chan
    valueChan              = make(chan string, 50)  // value to decide chan
    newSlotChan            = make(chan int, 50)
	//debug chan int
)

// Init function
func EntryPoint(p []int) (chan int, chan string, chan string, chan int) {
    process = p
	//debug = deb
    go loop()
    return handlTrustChan, handlPromiseLeaderChan, valueChan, newSlotChan
}
// what we do when we become the leader : prepare all the next slots
func gotTrust(leader int) {
   	slots := slotsManager.GetAvailableSlots()
    for i := range slots {
		if slotsManager.GetValueLearned(slots[i])=="" {
			//slotsManager.SetInWork(i,true)
			prepareSlot(slots[i])
			time.Sleep(10*time.Millisecond)
		}
    }
}
func prepareSlot(slot int) {
    message := "Prepare@" + strconv.Itoa(pickNext(slot)) + "@" + strconv.Itoa(slot)
    for pr := range process {
        proc := process[pr]
        preSend(message, proc)
    }
}
// the round number increase function
func pickNext(slot int) int {
    currentRound := 1 + slotsManager.GetRoundNumber(slot) // TO TEST
    slotsManager.SetRoundNumber(slot, currentRound)
    return currentRound
}
// function handling the promise message
func gotPromise(data string) {
    res := strings.Split(data, "@")
    roundnumber, _ := strconv.Atoi(res[1])
    lastVotedRound, _ := strconv.Atoi(res[2])
    lastVotedValue := res[3]
    slot, _ := strconv.Atoi(res[4])
    slotRN := slotsManager.GetRoundNumber(slot)
    if roundnumber == slotRN {
        slotsManager.AddToPromiseMap(slot, lastVotedRound, lastVotedValue)
        if lastVotedRound > slotsManager.GetMaxRoundInPromises(slot) {
            slotsManager.SetMaxRoundInPromises(slot, lastVotedRound)
        }
        if slotsManager.GetCptPromise(slot) > len(process)/2 {
			waitForValue(slot)
			slotsManager.ClearPromiseMap(slot)
            curR := strconv.Itoa(slotRN) 
            sendAll("Accept@" + curR + "@" + slotsManager.GetValueToDecide(slot) + "@" + strconv.Itoa(slot))
        }
    }
}
// when we are waiting for a value from a client
func waitForValue(slot int) {	
	if slotsManager.GetMaxRoundInPromises(slot) == 0 {
		preValue := <-valueChan 
		str := strings.Split(preValue, "@")
	    valueToDecide := str[1]
		slotsManager.SetValueToDecide(slot, valueToDecide)
		fmt.Println("Value to decide received :", valueToDecide, "in slot", strconv.Itoa(slot))
		
	} else {
		slotsManager.SetValueToDecide(slot,slotsManager.GetFromPromiseMap(slot, slotsManager.GetMaxRoundInPromises(slot)))
	}
}
func sendAll(message string) {
    for pr := range process {
        proc := process[pr]
        preSend(message, proc)
    }
}
// to check if we are the leader
func isLeader() bool {
    me, _ := connector.GetOwnProcess()
    leader := leaderElection.GetLeader()
    return me == leader
}
func loop() {
    for {
        select {
        case leader := <-handlTrustChan:
            gotTrust(leader)
        case data := <-handlPromiseLeaderChan:
            gotPromise(data)
        case newSlot := <-newSlotChan:
			if isLeader() {
	            prepareSlot(newSlot)										
			}										
        }
    }
}
func preSend(message string, pr int) {
    connector.Send(message, pr)
}

