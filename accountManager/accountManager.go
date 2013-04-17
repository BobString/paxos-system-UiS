package accountManager

import (
	"strings"
	//"fmt"
	"strconv"
	"connector"
)

var (
	accounts    = make(map[int]int)
	actionsChan = make(chan string, 20)
	message string
)

func EntryPoint() chan string {
	go loop()
	return actionsChan
}

func loop() {
	for {
		select {

		case action := <-actionsChan:
			manageMessage(action)
		}

	}
}

func manageMessage(action string) {
	res := strings.Split(action,"@")
	string1 := res[1]
	res = strings.Split(string1, ",")
	switch res[0] {
	case "Balance":
		balance(res[1], res[2])
	case "Withdraw":
		withdraw(res[1], res[2], res[3])
	case "Deposit":
		deposit(res[1], res[2], res[3])
	case "Transfer":
		transfer(res[1], res[2], res[3], res[4])
	}

}

func balance(account string, client string) {
	accountNum, _ := strconv.Atoi(account)
	if _, exists := accounts[accountNum]; exists {
		res := accounts[accountNum]
		message = "Account number "+ account+ " has "+ strconv.Itoa(res)+ " credits"
		//Send message to the client with the answer
	} else {
		message = "Sorry the asked account ("+account+") does not exist. Please check the account number"
		//Send message to the client saying that that account doesn't exist
	}
	sendToClient(message,client)
}
func withdraw(account string, amount string, client string) {
	accountNum, _ := strconv.Atoi(account)
	amountNum, _ := strconv.Atoi(amount)
	if _, exists := accounts[accountNum]; exists {
		res := accounts[accountNum]
		if amountNum > res {
			//Asking more money than available, send ERROR
			message = "Error : The account number "+account+" does not have enough money to withdraw the required amount ("+amount+")"
		} else {
			accounts[accountNum] = res - amountNum
			//Send message to the client with the new ammotun
			message = amount+" credits have been withdrawn successfully from account "+account+"."
		}

	} else {
		//Send message to the client saying that that account doesn't exist
		message = "Sorry the asked account ("+account+") does not exist. Please check the account number"
	}
	sendToClient(message,client)
}
func deposit(account string, amount string, client string) {
	accountNum, _ := strconv.Atoi(account)
	amountNum, _ := strconv.Atoi(amount)
	if _, exists := accounts[accountNum]; exists {
		res := accounts[accountNum]
		accounts[accountNum] = res + amountNum
		//Send message to the client with the new ammotun
		message = amount+" credits have been added to the account "+ account+". The new balance is "+strconv.Itoa(accounts[accountNum])

	} else {
		accounts[accountNum] = amountNum
		//Send message to the client saying that the account is created with the money
		message = "The account did not exist. It has been created with the following amount : "+amount +" credits. If there is any complaint, please ask our support service"
	}
}
func transfer(accFrom string, accTo string, amount string, client string) {
	accountNumTo, _ := strconv.Atoi(accTo)
	accountNumFrom, _ := strconv.Atoi(accFrom)
	amountNum, _ := strconv.Atoi(amount)
	_, existsTo := accounts[accountNumTo]
	_, existsFrom := accounts[accountNumFrom]
	if existsTo && existsFrom {
		resTo := accounts[accountNumTo]
		resFrom := accounts[accountNumFrom]
		if amountNum > resFrom {
			//Asking more money than available, send ERROR
			message = "Transfer error : The account of origin does not have the required amount"
		} else {
			accounts[accountNumFrom] = resFrom - amountNum
			accounts[accountNumTo] = resTo + amountNum
			//Send message to the client with the new ammotun
			message = "Transfer successful : account "+accFrom+" has now "+strconv.Itoa(accounts[accountNumFrom])+" credits; and account "+ accTo+ " has "+strconv.Itoa(accounts[accountNumTo])+" credits."
		}
		//Send message to the client with the new ammotun

	} else {
		//Send message to the client saying that that account doesn't exist
		message = "Transfer error : one of the accounts does not exist"
	}
}


func sendToClient(message string, client string){
	connector.SendByAddr(message,client)
}