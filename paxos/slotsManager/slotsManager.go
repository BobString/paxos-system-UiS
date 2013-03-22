// the type of the values inside the map
type mapValueType struct {
	// !!!!!!!!!!!!!!!!!!!!! TODO : do not forget ine the go files to call the values with slosManager.Get...(slot)
	roundNumber,lastVotedRN int
	lastVotedVal string
	valueLearned string
}
// the map in question 
var slotMap = make(map[int] string)


//////////// GETTERS ////////////
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
