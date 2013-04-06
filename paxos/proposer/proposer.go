
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
	debug chan int
    //maxRound               int = 0 //
    //cptProm             	  int = 0
)
// Init function
func EntryPoint(p []int, deb chan int) (chan int, chan string, chan string, chan int) {
    process = p
	debug = deb
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
	/*println("WAITING FOR DEBUG#######, slot",strconv.Itoa(slot))
	<- debug*/
    //////////////////// MESSAGE FORMT : Prepare@RN@slot
    message := "Prepare@" + strconv.Itoa(pickNext(slot)) + "@" + strconv.Itoa(slot)
	println("PREPARE TO BE SENT :", message,"in slot",strconv.Itoa(slot))
    for pr := range process {
        proc := process[pr]
        preSend(message, proc)
    }
}
// the round number increase function
func pickNext(slot int) int {
    currentRound := 1 + slotsManager.GetRoundNumber(slot) // TO TEST
    slotsManager.SetRoundNumber(slot, currentRound)
	println("*************************************************************0")
	println("Slot ",strconv.Itoa(slot),":")
	println("Round number set :", strconv.Itoa(currentRound), "and get :", strconv.Itoa(slotsManager.GetRoundNumber(slot)))
	println("*************************************************************0")
    return currentRound
}
// function handling the promise message
func gotPromise(data string) {
    res := strings.Split(data, "@")
    roundnumber, _ := strconv.Atoi(res[1])
    lastVotedRound, _ := strconv.Atoi(res[2])
    lastVotedValue := res[3]
    slot, _ := strconv.Atoi(res[4])
    //processID := res[4]
    slotRN := slotsManager.GetRoundNumber(slot)
	//println("############PROMISE DECRYPTED", strconv.Itoa(roundnumber),strconv.Itoa(slotRN))
    if roundnumber == slotRN {
		//println("$$$$$$$$ADDING TO PROMISE MAP") 
        slotsManager.AddToPromiseMap(slot, lastVotedRound, lastVotedValue)
        if lastVotedRound > slotsManager.GetMaxRoundInPromises(slot) {
            slotsManager.SetMaxRoundInPromises(slot, lastVotedRound)
        }
        if slotsManager.GetCptPromise(slot) > len(process)/2 {
   			//println("WAITING FOR VALUE ##########")         
			waitForValue(slot)
			//println("VALUE DECIDED !!!!!!!!!!!!!!!!!")
			slotsManager.ClearPromiseMap(slot)
            curR := strconv.Itoa(slotRN) //// FIXME : modify for slotsManager !!!!!!!!!!
            sendAll("Accept@" + curR + "@" + slotsManager.GetValueToDecide(slot) + "@" + strconv.Itoa(slot))
        }
    }
}
func waitForValue(slot int) {	
	if slotsManager.GetMaxRoundInPromises(slot) == 0 {
		for {
			select{
			case preValue := <-valueChan :
				str := strings.Split(preValue, "@")
			    valueToDecide := str[1]
				slotsManager.SetValueToDecide(slot, valueToDecide)
				fmt.Println("Value to decide received :", valueToDecide, "in slot", strconv.Itoa(slot))
			default:
				if slotsManager.HasLearned(slot) {
					break
				}
			}
		}
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
            prepareSlot(newSlot)
        }
    }
}
func preSend(message string, pr int) {
    connector.Send(message, pr, nil)
}

