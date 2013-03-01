package main

import (
	"connector"
	"fmt"
)

func main() {
	var val string
	var pr int
	fmt.Println("Enter the value to send")
	_, err := fmt.Scanln(&val)
	if err != nil {
		fmt.Println("read error")
	}
	valueMess := "Value@" + val
	//fmt.Println("Message to be sent :", valueMess)
	fmt.Println("Please select a process between 1 and 3")
	_, err = fmt.Scanln(&pr)
	if err != nil {
		fmt.Println("read error")
	}
	connector.Send(valueMess, pr, nil)
}
