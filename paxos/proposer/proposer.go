
package proposer
import (
    "connector"
    "fmt"
    "leaderElection"
    "paxos/slotsManager"
    "strconv"
    "strings"
    //"time"
)
// global variables
var (
    //mv                     = map[int]string{} //Round-vote map
    process []int // list of processes
    leader  int   // currentLeader
    //currentRound           int = 0 //our ID at first
    //systemRound            int = 0 // current system round number
    handlTrustChan = make(chan int, 20) // trust receive chan
    //handlSysRoundChan          = make(chan int, 20) // system round number receive chan
    handlPromiseLeaderChan = make(chan string, 20) // promise receive chan
    valueChan              = make(chan string, 5)  // value to decide chan
    newSlotChan            = make(chan int, 5)
    //maxRound               int = 0 //
    //cptProm             	  int = 0
)
// Init function
func EntryPoint(p []int) (chan int, chan string, chan string, chan int) {
    process = p
    //handlSysRoundChan = sysRoundChan
    //currentRound, _ = connector.GetOwnProcess()
    go loop()
    return handlTrustChan, handlPromiseLeaderChan, valueChan, newSlotChan
}
// what we do when we become the leader : change the round number
func gotTrust(leader int) {
	println("GOT THE TRUST !! STARTING PREPARING SLOTS")
    slots := slotsManager.GetAvailableSlots()
    for i := range slots {
		println("PREPARING SLOT",strconv.Itoa(slots[i]))
        prepareSlot(slots[i])
    }
}
func prepareSlot(slot int) {
    for pr := range process {
        proc := process[pr]
        //////////////////// MESSAGE FORMAT : Prepare@RN@slot
        message := "Prepare@" + strconv.Itoa(pickNext(slot)) + "@" + strconv.Itoa(slot)
		println("PREPARE TO BE SENT :", message,"in slot",strconv.Itoa(slot))
        preSend(message, proc)
    }
}
// the round number increase function
func pickNext(slot int) int {
    currentRound := slotsManager.GetRoundNumber(slot) + 1 // TO TEST
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
    //slotsManager.     FIWME : change here the cptPromise ?
    slotsManager.IncCptProm(slot)
    //processID := res[4]
    slotRN := slotsManager.GetRoundNumber(slot)
    if roundnumber == slotRN { //// FIXME : modify for slotsManager !!!!!!!!!!
        //aux := lastVotedRound//// FIXME : modify for slotsManager !!!!!!!!!!
        //mv[aux] = lastVotedValue//// FIXME : modify for slotsManager !!!!!!!!!!
        slotsManager.AddToPromiseMap(slot, lastVotedRound, lastVotedValue)
        if lastVotedRound > slotsManager.GetMaxRoundInPromises(slot) {
            slotsManager.SetMaxRoundInPromises(slot, lastVotedRound)
        }
        if slotsManager.GetCptPromise(slot) >= len(process)/2 {
            var valueToDecide string // the value we want to decide
            if slotsManager.GetMaxRoundInPromises(slot) == 0 {
                preValue := <-valueChan
                str := strings.Split(preValue, "@")
                valueToDecide = str[1]
                fmt.Println("Value to decide received :", valueToDecide, "in slot", res[4])
                //// FIXME : modify for slotsManager !!!!!!!!!!
            } else {
                //Pick the value form the largest round
                valueToDecide = slotsManager.GetFromPromiseMap(slot, slotsManager.GetMaxRoundInPromises(slot))
            }
            curR := strconv.Itoa(slotRN) //// FIXME : modify for slotsManager !!!!!!!!!!
            sendAll("Accept@" + curR + "@" + valueToDecide)
        }
    }
}
func sendAll(message string) {
    for pr := range process {
        proc := process[pr]
        preSend(message, proc)
    }
}
func isLeader() bool {
    me, _ := connector.GetOwnProcess()
    leader := leaderElection.GetLeader()
    return me == leader
}
func loop() {
    for {
        select {
        case leader := <-handlTrustChan:
            go gotTrust(leader)
        case data := <-handlPromiseLeaderChan:
            go gotPromise(data)
        case newSlot := <-newSlotChan:
            go prepareSlot(newSlot)
            // TODO : add the case when we receive a slot number from the learner, so we have to send prepare messages !!
        }
    }
}
func preSend(message string, pr int) {
    connector.Send(message, pr, nil)
}

