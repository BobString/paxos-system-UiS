package failureDetector

import (
	"net"
	"time"
)

var (
	//process   []int = make([]int, 20)
	pAlive        = map[int]bool{}
	pSuspect      = map[int]bool{}
	process       []int
	delay         int = 0
	actualTimeout int = 0
)

func init(d int, p []int) {
	delay = d
	actualTimeout = d
	process = p
	//Take the process and make 2 maps alive and suspect
	for proc := range p {
		pAlive[proc] = true
		pSuspect[proc] = false

	}
	startTimer(actualTimeout)
}

func timeout(ticker *Ticker) {
	//First part to manage the delay
	var flag bool = false
	for pr := range process {
		if pSuspect[pr] {
			//If there any project that is suspect, increase the timeout
			actualTimeout = actualTimeout + delay
			//We stop the last ticker
			ticker.Stop()
			flag = true
			break
		}
	}

	//Second part to manage the suspect process
	for pr := range process {
		if !pAlive[pr] && !pSuspect[pr] {
			pSuspect[pr] = true

			//Trigger Suspect, pr
			//FIXME: Send to all process??
			//err := send("Suspect", pr)
			//if err != nil {
			//	fmt.Printf("Error sending Suspect to ", pr)

			//} else {
			//	fmt.Printf("Suspect!!", pr)
			//}
		} else if pAlive[pr] && pSuspect[pr] {
			pSuspect[pr] = false

			//Trigger Restore, pr
			//FIXME: Send to all process??
			//err := send("Restore", pr)
			//if err != nil {
			//	fmt.Printf("Error sending Restore to ", pr)

			//} else {
			//	fmt.Printf("Restored!!", pr)
			//}

		}
		//Sent HeartbeatRequest to pr	
		err := send("HeartbeatRequest", pr)
		if err != nil {
			fmt.Printf("Error sending HeartbeatRequest to ", pr)

		}

	}

	//Put pAlive all to false
	for pr := range process {
		pAlive[proc] = true
	}

	//Start timer again only if the delay is changed
	if flag {
		startTimer(actualTimeout)
	}
}

func gotHeartBeatRequest(pr int) {
	//Sent HeartbeatReply to pr
	err := send("HeartbeatReply", pr)
	if err != nil {
		fmt.Printf("Error sending HeartbeatReply to ", pr)

	}
}

func gotHeartBeatReply(pr int) {
	pAlive[pr] = true
}

func send(message string, pr int) {

	//Send a message through TCP
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	service := process[pr]
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	//Just in case the message has an @
	message = strings.Replace(message, "@", "(at)", -1)
	_, err = conn.Write([]byte(message + "@" + pr))
	conn.Close()
	return err

}
func startTimer(time) {
	//Time in seconds
	c := time.NewTicker(time * time.Second)
	for now := range c.C {
		fmt.Printf("Timeout!")
		//Call timeout function every time seconds
		timeout(s)
	}
}
