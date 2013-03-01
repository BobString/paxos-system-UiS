package acceptor
//@author Remy Pannaud


import (
	"connector"
	"strings"
)

// global variables : 
var (
	learnList []int // list of known learners
	rnd int // current round number
	lvrn int // last voted round number
	lvval int // last voted calue
	inPrepChan = make(chan string,5)
	inAcceptChan = make(chan string,5)
)

// functions :
// 	EntryPoint
// 	prepareListener
// 	acceptListener

// Initialization function
//@parameters :
//	inPrepChan : when the server receives a Prepare message, it is sent to inPrepChan
//	inAcceptChan : idem with Accept messages
// 	sendRoundChan : to send to the Proposer the current system round number
func EntryPoint (list []int, sendRoundChan chan int) (chan string,chan string) {
	learnList = list
	go prepareListener(inPrepChan)
	go acceptListener(inAcceptChan,sendRoundChan)
	return inPrepChan,inAcceptChan
}


func prepareListener (inPrepChan chan string) {
	for {
		v,_ := <-inPrepChan
		strings.Split(v,"@")
		if int(v[1]) > lvrn {
			promise := "Promise@"+string(v[1])+"@"+lvrn+"@"+lvval+"@"
			preSend(promise,int(v[2]))
		}
	}
}

func acceptListener(inAcceptChan chan string, sendRoundChan chan int) {
	for {
		v,_ := <-inAcceptChan		
		strings.Split(v,"@")
		if int(v[1])>= lvrn {
			learn := "Learn@"+string(v[1])+"@"+string(v[2])+"@"
			lvrn = int(v[1])
			sendRoundChan <- lvrn
			lvval = int(v[2])
			for i := range learnList {
				preSend(learn,learnList[i])
			}
		}
	}
}

func preSend(message string, pr int) {
	connector.Send(message, pr, nil)
}