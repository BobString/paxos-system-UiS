package slotsManager
import (
    //"fmt"
	//"strconv"
)
// the type of the values inside the map
type MapValueType struct {
    RoundNumber, LastVotedRN int // round numbers (current and last voted)
    LastVotedVal             string // the last voted value
    ValueLearned             string // the learned value
	ValueToDecide			 string // the value to be decided
    CptPromise               int // the promise counter
    MaxRoundInPromises       int // the maximum round retrieved from the promises
	InWork					 bool
    PromiseMap               map[int]string // the map of promises to check for last voted values
    LearnMap                 map[LearnPair]int // the map of pairs to decide the value
}
type LearnPair struct {
    Nv  int
    Val string
}
// the map in question
var index int
const sizeMax = 100
var slotMap map[int] MapValueType

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
        //fmt.Println("Accessing an unlearned value")
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
/*func IsInWork (slot int) bool {
	return slotMap[slot].InWork
}*/


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
	ClearLearnMap(slot)
    slotType := slotMap[slot]
    slotType.ValueLearned = val
	slotMap[slot]=slotType
	ClearPromiseMap(slot)
    // a value has been learned : we re initialise an entry in the map
	initSlot()
    // and we return the value of the added slot 
	nextSlot := index - 3
	if nextSlot < 1 {
		nextSlot = nextSlot + sizeMax
	}
    return nextSlot
}
func SetMaxRoundInPromises(slot, maxR int) {
    slotType := slotMap[slot]
    slotType.MaxRoundInPromises = maxR
	slotMap[slot]=slotType
}
/*func SetInWork (slot int, val bool) {
	slotType := slotMap[slot]
    slotType.InWork = val
	slotMap[slot]=slotType
}*/
// a local function to increase the counter of promises
func incCptProm(slot int) {
    slotType := slotMap[slot]
    slotType.CptPromise = slotType.CptPromise + 1
	slotMap[slot]=slotType
}


/////// FUNCTIONS ON MAPS ///////
//// PROMISE MAP
// to add to the map
func AddToPromiseMap(slot int, key int, val string) {
	slotType := slotMap[slot]
	auxMap := slotType.PromiseMap
	auxMap[key] = val
	slotType.PromiseMap = auxMap
	slotMap[slot] = slotType
	incCptProm(slot)
}
// to clear the map
func ClearPromiseMap(slot int) {
    var mapAux map[int]string
	slotType := slotMap[slot]
	mapAux = slotType.PromiseMap
    for v, _ := range mapAux {
		delete(mapAux, v)
    }
	slotType.PromiseMap = mapAux
	slotType.CptPromise = 0
	slotMap[slot] = slotType
	
}
// to get an element from the map 
func GetFromPromiseMap(slot int, key int) string {
    return slotMap[slot].PromiseMap[key]
}
//// LEARN MAP
// to add to the map
func AddToLearnMap(slot int, key LearnPair, val int) {
    slotType := slotMap[slot]
	auxMap := slotType.LearnMap
	auxMap[key] = val
	slotType.LearnMap = auxMap
	slotMap[slot] = slotType
}
// To clear the map
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
// a getter for the map elements
func GetFromLearnMap(slot int, key LearnPair) int {
	var result int
	if c,is := slotMap[slot].LearnMap[key] ; is {
		result = c
	} else {
		result = 0
	}
    return result
}
// A classic belongsTo function
func BelongsToLearnMap(slot int, key LearnPair) bool {
    _, ok := slotMap[slot].LearnMap[key]
    return ok
}


// returns the smallest slot with no learned value yet
func getSmallestUnused() int {
    i := index 	
    for GetValueLearned(i)=="" {
        i = i + 1
		if i>sizeMax {
			break
		}
    }
    return i
}
// returns the biggest slot number greater with no learned value
func getBiggestUnused() int {
	i := getSmallestUnused()
	stop := i
	for !(GetValueLearned(i+1)=="") {
        i = i + 1
		if i == stop {
			if i!=1 {
				i = i-1
			} else {
				i = sizeMax
			}			
			break
		}
		if i>=sizeMax {
			i = 0
		}
    }
	return i
}
// returns the available slots (with no learned value) as a slice of int
func GetAvailableSlots() []int {
	// verified
    slotMin := getSmallestUnused()
    slotMax := getBiggestUnused()
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
    }/*
	res := make(map[int]int)
	ind := 0
	for i:=1; i<=sizeMax ; i++ {
		if GetValueLearned(i)=="" {
			res[ind] = i
			ind = ind +1
		}
	}*/
    return res
}

func initSlot () {
	index = index +1
	if index > sizeMax {
		index = index - sizeMax
	}	
	_,ok := slotMap[index]
	if ok {	
		ClearPromiseMap(index)
		ClearLearnMap(index)
		slotType := slotMap[index]
		slotType.RoundNumber = 0
		slotType.LastVotedRN = 0
		slotType.LastVotedVal = ""
		slotType.ValueLearned = ""
		slotType.ValueToDecide = ""
		slotType.CptPromise = 0
		slotType.MaxRoundInPromises = 0	
		slotType.InWork = false
		slotMap[index] = slotType
	} else {
		promMap := make(map[int]string)
   		leaMap := make(map[LearnPair]int)
   		mapValueNil := MapValueType{0, 0, "", "", "", 0, 0,false, promMap, leaMap}
		slotMap[index] = mapValueNil
	}	
}
func EntryPoint() {
	index = 1
	slotMap = make(map[int] MapValueType,sizeMax)
    for i := 1; i <= sizeMax; i++ {
        initSlot()
    }
	index = 3
}


