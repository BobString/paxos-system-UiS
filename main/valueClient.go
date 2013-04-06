package main

import (
	"connector"
	"fmt"
	"os"
	"time"
	"strconv"
)

func main() {
	var val,menu,mess string
	var pr,number int
	for {
		fmt.Println("Choose what you want to send")
		fmt.Println("1 : Value")
		fmt.Println("2 : StopServer")
		fmt.Println("3 : RestoreServer")
		fmt.Println("4 : Exit program")
		fmt.Println("5 : 500 values")
		_, err := fmt.Scanln(&menu)
		checkErr(err)
		valmenu,_ := strconv.Atoi(menu)
		number = 1
		switch valmenu {
			case 1:
				fmt.Println("Enter the value to send")
				_, err := fmt.Scanln(&val)
				checkErr(err)	
				mess = "Value@" + val
			case 2:
				mess = "StopServer@"
			case 3:
				mess = "RestoreServer@"
			case 4:
				os.Exit(0)
			case 5:
				mess = "Value@"
				number = 500
		}
		fmt.Println("Please select a process between 1 and 3")
		_, err = fmt.Scanln(&pr)
		checkErr(err)
		for i:=1;i<=number;i++ {
			if number>1 {
				mess = "Value@" + strconv.Itoa(i)
			}
			connector.Send(mess, pr, nil)
			time.Sleep(100*time.Millisecond)
		}
	}
}

func checkErr (err error) {
	if err != nil {
		fmt.Println("read error")
	}
}
