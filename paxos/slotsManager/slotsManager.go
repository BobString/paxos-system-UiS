// the type of the values inside the map
type mapValueType struct {
	// !!!!!!!!!!!!!!!!!!!!! TODO : do not forget ine the go files to call the values with slosManager.Get...(slot)
	roundNumber,lastVotedRN int // 
	lastVotedVal string
	valueLearned string
	cptPromise int
	maxRoundInPromises int
	promiseMap = make(map[int] string)
	learnMap = make(map[LearnPair] int) // the map of pairs to decide the value
}
type LearnPair struct {
	nv int
	val string
}
// the map in question 
var slotMap = make(map[int] string)


//////////// GETTERS ////////////
func GetSlots () []int {
	var res []int = make([]int, len(slotMap))
	i := 0
	for v,_ := range slotMap {
		res[i] = v
		i = i+1
	}
}
func GetRoundNumber (slot int) int {
	return slotMap[slot].roundNumber	
}
func GetLastVotedRN (slot int) int {
	return slotMap[slot].lastVotedRN		
}
func GetLastVotedVal (slot int) string {
	return slotMap[slot].lastVotedVal
}
func GetValueLearned(slot int) string {
	if slotMap[slot].valueLearned=="" {
		fmt.Println("Accessing an unlearned value")
	}
	return slotMap[slot].valueLearned
}
func GetCptPromise (slot int) int {
	return slotMap[slot].cptPromise
}
func GetMaxRoundInPromises (slot int) int {
	return slotMap[slot].maxRoundInPromises
}
func GetValueFromLearnPair (p LearnPair) {
	return p.val
}

//////////// SETTERS ////////////
func SetRoundNumber (slot int, val int) {
	slotMap[slot].roundNumber = val	
}
func SetLastVotedRN (slot int, val int) {
	slotMap[slot].lastVotedRN = val		
}
func SetLastVotedVal (slot int, val string) {
	slotMap[slot].lastVotedVal = val
}
func SetValueToLearn (slot int, val string) int {
	slotMap[i].valueLearned = val
	// a value has been learned : we create a new entry in the map
	createNewEntry()
	// and we return the value of the added slot (TODO : to be sent to the proposer for a Prepare message sending)
	return len(slotMap)
}
func SetMaxRoundInPromises (slot,maxR int) {
	slotMap[slot].maxRoundInPromises = maxR
}
func IncCptProm (slot int) {
	slotMap[slot].cptPromises = slotMap[slot].cptPromises +1
}


/////// FUNCTIONS ON MAPS ///////
//// Promise map
func AddToPromiseMap(slot int, key int, val string) {
	slotMap[slot].promiseMap[key] = val
}
func ClearPromiseMap(slot int) {
	for v,_ := range slotMap[slot].promiseMap {
		delete(slotMap[slot].promiseMap,v)		
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
	for v,_ := range slotMap[slot].learnMap {
		delete(slotMap[slot].learnMap,v)
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
func getSmallestUnlearned int {
	i := 0
	for	getValueLearned(i)!="" {
		i = i+1
	}
	return i
}

// returns the available slots (with no learned value)
func GetAvailableSlots () []int {
	slot :=getSmallestUnlearned()
	index := 0
	res []int
	for slot<len(slotMap) {
		if GetValueLearned(slot)=="" {
			res[index] = slot
			index = index +1
		}
		slot = slot +1
	}
	return res
}

// adds an empty slot to the map
func createNewEntry() {
	var mapValueNil := mapValueType {false,0,0,"",""}
	slotMap[len(slotMap)+1] = mapValueNil 	
}


func EntryPoint() {
	for i:=0;i<10;i++ {
		createNewEntry()
	}
}
