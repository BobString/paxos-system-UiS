package slotsManager
import (
    "fmt"
)
// the type of the values inside the map
type mapValueType struct {
    // !!!!!!!!!!!!!!!!!!!!! TODO : do not forget ine the go files to call the values with slosManager.Get...(slot)
    roundNumber, lastVotedRN int //
    lastVotedVal             string
    valueLearned             string
    cptPromise               int
    maxRoundInPromises       int
    promiseMap               map[int]string
    learnMap                 map[LearnPair]int // the map of pairs to decide the value
}
type LearnPair struct {
    Nv  int
    Val string
}
// the map in question
var slotMap = make(map[int]mapValueType)
//////////// GETTERS ////////////
func GetRoundNumber(slot int) int {
    return slotMap[slot].roundNumber
}
func GetLastVotedRN(slot int) int {
    return slotMap[slot].lastVotedRN
}
func GetLastVotedVal(slot int) string {
    return slotMap[slot].lastVotedVal
}
func GetValueLearned(slot int) string {
    if slotMap[slot].valueLearned == "" {
        fmt.Println("Accessing an unlearned value")
    }
    return slotMap[slot].valueLearned
}
func GetCptPromise(slot int) int {
    return slotMap[slot].cptPromise
}
func GetMaxRoundInPromises(slot int) int {
    return slotMap[slot].maxRoundInPromises
}
func GetValueFromLearnPair(p LearnPair) string {
    return p.Val
}
//////////// SETTERS ////////////
func SetRoundNumber(slot int, val int) {
    slotType := slotMap[slot]
    slotType.roundNumber = val
}
func SetLastVotedRN(slot int, val int) {
    slotType := slotMap[slot]
    slotType.lastVotedRN = val
}
func SetLastVotedVal(slot int, val string) {
    slotType := slotMap[slot]
    slotType.lastVotedVal = val
}
func SetValueToLearn(slot int, val string) int {
    slotType := slotMap[slot]
    slotType.valueLearned = val
    // a value has been learned : we create a new entry in the map
    createNewEntry()
    // and we return the value of the added slot (TODO : to be sent to the proposer for a Prepare message sending)
    return len(slotMap)
}
func SetMaxRoundInPromises(slot, maxR int) {
    slotType := slotMap[slot]
    slotType.maxRoundInPromises = maxR
}
func IncCptProm(slot int) {
    slotType := slotMap[slot]
    slotType.cptPromise = slotType.cptPromise + 1
}
/////// FUNCTIONS ON MAPS ///////
//// Promise map
func AddToPromiseMap(slot int, key int, val string) {
    slotMap[slot].promiseMap[key] = val
}
func ClearPromiseMap(slot int) {
    for v, _ := range slotMap[slot].promiseMap {
        delete(slotMap[slot].promiseMap, v)
    }
}
func GetFromPromiseMap(slot int, key int) string {
    return slotMap[slot].promiseMap[key]
}
//// Learn map
func AddToLearnMap(slot int, key LearnPair, val int) {
    slotMap[slot].learnMap[key] = val
}
func ClearLearnMap(slot int) {
    for v, _ := range slotMap[slot].learnMap {
        delete(slotMap[slot].learnMap, v)
    }
}
func GetFromLearnMap(slot int, key LearnPair) int {
    return slotMap[slot].learnMap[key]
}
func BelongsToLearnMap(slot int, key LearnPair) bool {
    _, ok := slotMap[slot].learnMap[key]
    return ok
}
// returns the smallest slot with no learned value yet
func getSmallestUnlearned() int {
    i := 1
    for slotMap[i].valueLearned != "" {
        i = i + 1
    }
    return i
}
// returns the available slots (with no learned value)
func GetAvailableSlots() []int {
    slotMin := getSmallestUnlearned()
    slotMax := len(slotMap)
	length := slotMax-slotMin+1
	index := 0
	res = make([]int,length)
    for index < length {
		res[index] = slotMin + index
        index = index + 1
    }
    return res
}
// adds an empty slot to the map
func createNewEntry() {
    promMap := make(map[int]string)
    leaMap := make(map[LearnPair]int)
    mapValueNil := mapValueType{0, 0, "", "", 0, 0, promMap, leaMap}
    slotMap[len(slotMap)+1] = mapValueNil
}
func EntryPoint() {
    for i := 0; i < 10; i++ {
        createNewEntry()
    }
}


