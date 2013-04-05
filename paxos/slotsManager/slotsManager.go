package slotsManager
import (
    "fmt"
)
// the type of the values inside the map
type MapValueType struct {
    // !!!!!!!!!!!!!!!!!!!!! TODO : do not forget ine the go files to call the values with slosManager.Get...(slot)
    RoundNumber, LastVotedRN int //
    LastVotedVal             string
    ValueLearned             string
	ValueToDecide			 string
    CptPromise               int
    MaxRoundInPromises       int
    PromiseMap               map[int]string
    LearnMap                 map[LearnPair]int // the map of pairs to decide the value
}
type LearnPair struct {
    Nv  int
    Val string
}
// the map in question
var slotMap = make(map[int] MapValueType)
//////////// GETTERS ////////////
func GetRoundNumber(slot int) int {
    return slotMap[slot].RoundNumber
}
func GetLastVotedRN(slot int) int {
    return slotMap[slot].LastVotedRN
}
func GetLastVotedVal(slot int) string {
    return slotMap[slot].LastVotedVal
}
func GetValueLearned(slot int) string {
    if slotMap[slot].ValueLearned == "" {
        fmt.Println("Accessing an unlearned value")
    }
    return slotMap[slot].ValueLearned
}
func GetValueToDecide (slot int) string {
	return slotMap[slot].ValueToDecide
}
func GetCptPromise(slot int) int {
    return slotMap[slot].CptPromise
}
func GetMaxRoundInPromises(slot int) int {
    return slotMap[slot].MaxRoundInPromises
}
func GetValueFromLearnPair(p LearnPair) string {
    return p.Val
}
//////////// SETTERS ////////////
func SetRoundNumber(slot int, val int) {
    slotType := slotMap[slot]
    slotType.RoundNumber = val
}
func SetLastVotedRN(slot int, val int) {
    slotType := slotMap[slot]
    slotType.LastVotedRN = val
}
func SetLastVotedVal(slot int, val string) {
    slotType := slotMap[slot]
    slotType.LastVotedVal = val
}
func SetValueToDecide (slot int, val string){
	slotType := slotMap[slot]
	slotType.ValueToDecide = val
}
func SetValueToLearn(slot int, val string) int {
    slotType := slotMap[slot]
    slotType.ValueLearned = val
    // a value has been learned : we create a new entry in the map
    createNewEntry()
    // and we return the value of the added slot (TODO : to be sent to the proposer for a Prepare message sending)
    return len(slotMap)
}
func SetMaxRoundInPromises(slot, maxR int) {
    slotType := slotMap[slot]
    slotType.MaxRoundInPromises = maxR
}
func IncCptProm(slot int) {
    slotType := slotMap[slot]
    slotType.CptPromise = slotType.CptPromise + 1
}
/////// FUNCTIONS ON MAPS ///////
//// Promise map
func AddToPromiseMap(slot int, key int, val string) {
    slotMap[slot].PromiseMap[key] = val
}
func ClearPromiseMap(slot int) {
    for v, _ := range slotMap[slot].PromiseMap {
        delete(slotMap[slot].PromiseMap, v)
    }
}
func GetFromPromiseMap(slot int, key int) string {
    return slotMap[slot].PromiseMap[key]
}
//// Learn map
func AddToLearnMap(slot int, key LearnPair, val int) {
    slotMap[slot].LearnMap[key] = val
}
func ClearLearnMap(slot int) {
    for v, _ := range slotMap[slot].LearnMap {
        delete(slotMap[slot].LearnMap, v)
    }
}
func GetFromLearnMap(slot int, key LearnPair) int {
    return slotMap[slot].LearnMap[key]
}
func BelongsToLearnMap(slot int, key LearnPair) bool {
    _, ok := slotMap[slot].LearnMap[key]
    return ok
}
// returns the smallest slot with no learned value yet
func getSmallestUnlearned() int {
    i := 1
    for slotMap[i].ValueLearned != "" {
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
	res := make([]int,length)
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
    mapValueNil := mapValueType{0, 0, "", "", "", 0, 0, promMap, leaMap}
    slotMap[len(slotMap)+1] = mapValueNil
}
func EntryPoint() {
    for i := 0; i < 10; i++ {
        createNewEntry()
    }
}


