package accountManager

var (
	accounts    = make(map[int]int)
	actionsChan = make(chan string, 20)
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
	res = strings.Split(string1, ",")
	select {
	case res[0] == "balance":
		balance(res[1], res[2])
	case res[0] == "withdraw":
		withdraw(res[1], res[2], res[3])
	case res[0] == "deposit":
		deposit(res[1], res[2], res[3])
	case res[0] == "transfer":
		transfer(res[1], res[2], res[3], res[4])
	}

}

func balance(account string, client string) {
	accountNum, _ := strconv.Atoi(account)
	if _, exists := accounts[accountNum]; exists {
		res := accounts[accountNum]
		//Send message to the client with the answer
	} else {
		//Send message to the client saying that that account doesn't exist
	}
}
func withdraw(account string, amount string, client string) {
	accountNum, _ := strconv.Atoi(account)
	amountNum, _ := strconv.Atoi(amount)
	if _, exists := accounts[accountNum]; exists {
		res := accounts[accountNum]
		if amountNum > res {
			//Asking more money than available, send ERROR
		} else {
			accounts[accountNum] = res - amountNum
			//Send message to the client with the new ammotun
		}

	} else {
		//Send message to the client saying that that account doesn't exist
	}
}
func deposit(account string, amount string, client string) {
	accountNum, _ := strconv.Atoi(account)
	amountNum, _ := strconv.Atoi(amount)
	if _, exists := accounts[accountNum]; exists {
		res := accounts[accountNum]
		accounts[accountNum] = res + amountNum
		//Send message to the client with the new ammotun

	} else {
		accounts[accountNum] = amountNum
		//Send message to the client saying that the account is created with the money
	}
}
func transfer(accFrom string, accTo string, amount string, client string) string {
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
		} else {
			accounts[accountNumFrom] = resFrom - amountNum
			accounts[accountNumTo] = resTo + amountNum
			//Send message to the client with the new ammotun
		}
		//Send message to the client with the new ammotun

	} else {
		//Send message to the client saying that that account doesn't exist
	}
}
