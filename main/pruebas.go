package main

import (
	"connector"
	"fmt"
)

var (
	ownProcess int = 0
)

func main() {
	ownProcess, _ = connector.GetOwnProcess()
	_, err := connector.Send("HeartbeatRequest", ownProcess, nil)
	//TODO: If connection closed, error we have to treat the error.
	if err != nil {
		fmt.Printf("Error sending HeartbeatRequest to %d", ownProcess)
	}
}
