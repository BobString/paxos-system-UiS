package main

import (
	"connector"
	"fmt"
	"os"
)

func main() {
	var val,menu,mess string
	var pr int
	for {
		fmt.Println("Choose what you want to send")
		fmt.Println("1 : Value")
		fmt.Println("2 : StopServer")
		fmt.Println("3 : RestoreServer")
		fmt.Println("4 : Exit program")
		_, err := fmt.Scanln(&menu)
		checkErr(err)
		if menu == 4 {
			os.Exit(0)
		} 
		switch menu {
			case 1:
				fmt.Println("Enter the value to send")
				_, err := fmt.Scanln(&val)
				checkErr(err)	
				mess = "Value@" + val
			case 2:
				mess = "StopServer@"
			case 3:
				mess = "RestoreServer@"
		}
		fmt.Println("Please select a process between 1 and 3")
		_, err = fmt.Scanln(&pr)
		checkErr(err)
		connector.Send(mess, pr, nil)
	}
}

func checkErr (err error) {
	if err != nil {
		fmt.Println("read error")
	}
}
