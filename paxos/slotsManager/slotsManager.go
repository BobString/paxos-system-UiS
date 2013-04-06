package slotsManager
import (
    "fmt"
)
// the type of the values inside the map
type MapValueType struct {
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
var index int
const sizeMax = 100

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
	slotMap[slot]=slotType
}
func SetLastVotedRN(slot int, val int) {
    slotType := slotMap[slot]
    slotType.LastVotedRN = val
	slotMap[slot]=slotType
}
func SetLastVotedVal(slot int, val string) {
    slotType := slotMap[slot]
    slotType.LastVotedVal = val
	slotMap[slot]=slotType
}
func SetValueToDecide (slot int, val string){
	slotType := slotMap[slot]
	slotType.ValueToDecide = val
	slotMap[slot]=slotType
}
func SetValueToLearn(slot int, val string) int {
    slotType := slotMap[slot]
    slotType.ValueLearned = val
	slotMap[slot]=slotType
    // a value has been learned : we create a new entry in the map
    if index >= sizeMax {
		index = index - sizeMax
	}
	createNewEntry()
    // and we return the value of the added slot 
    return index
}
func SetMaxRoundInPromises(slot, maxR int) {
    slotType := slotMap[slot]
    slotType.MaxRoundInPromises = maxR
	slotMap[slot]=slotType
}
func IncCptProm(slot int) {
    slotType := slotMap[slot]
    slotType.CptPromise = slotType.CptPromise + 1
	slotMap[slot]=slotType
}


/////// FUNCTIONS ON MAPS ///////
//// Promise map
func AddToPromiseMap(slot int, key int, val string) {
	slotType := slotMap[slot]
    slotType.PromiseMap[key] = val
	slotMap[slot] = slotType
}
func ClearPromiseMap(slot int) {
    var mapAux map[int]string
	slotType := slotMap[slot]
	mapAux = slotType.PromiseMap
    for v, _ := range mapAux {
		delete(mapAux, v)
    }
	slotType.PromiseMap = mapAux
	slotMap[slot] = slotType
}
func GetFromPromiseMap(slot int, key int) string {
    return slotMap[slot].PromiseMap[key]
}
//// Learn map
func AddToLearnMap(slot int, key LearnPair, val int) {
    slotMap[slot].LearnMap[key] = val
}
func ClearLearnMap(slot int) {
	var mapAux map[LearnPair]int
	slotType := slotMap[slot]
	mapAux = slotType.LearnMap
    for v, _ := range mapAux {
		delete(mapAux, v)
    }
	slotType.LearnMap = mapAux
	slotMap[slot] = slotType
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
	// verified
    slotMin := getSmallestUnlearned()
    slotMax := index
	length := 0
	if slotMax < slotMin {
		length = sizeMax-slotMin+slotMax +1
	} else {
		length = slotMax-slotMin+1
	}	
	ind := 0
	res := make([]int,length)
    for i:=0;i<length;i++ {
		if slotMin+ind>sizeMax {
			ind = ind - sizeMax
		}
		res[i] = slotMin + ind
        ind = ind + 1
    }
    return res
}


// adds an empty slot to the map
func createNewEntry() {
    promMap := make(map[int]string)
    leaMap := make(map[LearnPair]int)
    mapValueNil := MapValueType{0, 0, "", "", "", 0, 0, promMap, leaMap}
	index = index+1
    slotMap[index] = mapValueNil
}
func EntryPoint() {
	index = 1
    for i := 1; i < 10; i++ {
        createNewEntry()
    }
}


