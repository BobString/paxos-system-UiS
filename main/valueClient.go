package main

import (
	"connector"
	"fmt"
	"net"
	//"os"
	//"time"
	"strconv"
)

func main() {
	var val,menu,mess string
	var pr,number int
	go createServer()
	for {
		fmt.Println("Choose what you want to send")
		fmt.Println("1 : Deposit")
		fmt.Println("2 : Withdraw")
		fmt.Println("3 : Transfer")
		fmt.Println("4 : Balance")
		/*fmt.Println("2 : StopServer")
		fmt.Println("3 : RestoreServer")
		fmt.Println("4 : Exit program")
		fmt.Println("5 : 500 values")*/
		_, err := fmt.Scanln(&menu)
		checkErr(err)
		valmenu,_ := strconv.Atoi(menu)
		number = 1
		switch valmenu {
			case 1:
				fmt.Println("Please enter the account number (D)")
				_, err := fmt.Scanln(&val)
				checkErr(err)
				mess = "Value@" + "Deposit,"+ val + ","
				fmt.Println("Enter the amount to deposit")
				_, err = fmt.Scanln(&val)
				checkErr(err)	
				mess = mess + val + ","
			case 2:
				fmt.Println("Please enter the account number (W)")
				_, err := fmt.Scanln(&val)
				checkErr(err)
				mess = "Value@" + "Withdraw,"+ val + ","
				fmt.Println("Enter the amount to withdraw")
				_, err = fmt.Scanln(&val)
				checkErr(err)	
				mess = mess + val + ","
			case 3:
				fmt.Println("Please enter the account number of origin")
				_, err := fmt.Scanln(&val)
				checkErr(err)
				mess = "Value@" + "Transfer,"+ val + ","
				fmt.Println("Please enter the account number of destination")
				_, err = fmt.Scanln(&val)
				checkErr(err)
				mess = mess + val + ","
				fmt.Println("Enter the amount to transfer")
				_, err = fmt.Scanln(&val)
				checkErr(err)
				mess = mess + val + ","
			case 4:
				fmt.Println("Please enter the account number (B)")
				_, err := fmt.Scanln(&val)
				checkErr(err)
				mess = "Value@" + "Balance,"+ val + ","
			default:
				fmt.Println("Please enter a correct value")
				continue
		}
		fmt.Println("Please select a process between 1 and 3")
		_, err = fmt.Scanln(&pr)
		checkErr(err)
		/*for i:=1;i<=number;i++ {
			if number>1 {
				mess = "Value@" + strconv.Itoa(i)
			}*/
			connector.Send(mess, pr)
			// a small time to let the paxos machine decide values and not slow the network down
			//time.Sleep(100*time.Millisecond)
		//}
	}
}

func createServer() {
	fmt.Println("Starting server...")
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			continue
		}

		//Maintaining the connection per node
		go handleClient(conn)

	}

}

func handleClient(conn net.Conn) {
	buf := make([]byte, 4096)
	_, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		//If the client close the connection we get out and start listening again
		break
	}
	println(string(buf))
}

func checkErr (err error) {
	if err != nil {
		fmt.Println("read error")
	}
}
