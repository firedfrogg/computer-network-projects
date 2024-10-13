package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("[CRASH] ", r)
	// 	}
	// }()

	fmt.Printf("Name: %s, ID Student: %s\n", Name, IdStudent)
	fmt.Println("========================================")
	fmt.Println("Welcome to Sigmart Point of Sales")
	fmt.Println("Please input your command below")
	fmt.Println("========================================")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		line := scanner.Text()
		err := scanner.Err()
		if err != nil {
			fmt.Println("[CRASH] ", err.Error())
			os.Exit(1)
		}

		spl := strings.Split(line, " ")
		executeCommand(spl[0], spl[1:])
	}
}

func executeCommand(command string, data []string) {
	switch command {
	case "ADD_ITEM":
		if len(data) < 4 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		price := int32(parseInt(data[2]))
		stockQty := int32(parseInt(data[3]))
		msg, err := AddItem(data[0], data[1], price, stockQty)
		PrintMessage(msg, err)
	case "DELETE_ITEM":
		if len(data) < 1 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		msg, err := DeleteItem(data[0])
		PrintMessage(msg, err)
	case "ADD_MEMBER":
		if len(data) < 2 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		msg, err := AddMember(data[0], data[1])
		PrintMessage(msg, err)
	case "DELETE_MEMBER":
		if len(data) < 1 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		msg, err := DeleteMember(data[0])
		PrintMessage(msg, err)
	case "ADD_TRANSACTION":
		if len(data) < 2 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		qty := int32(parseInt(data[0]))
		msg, err := AddTransaction(qty, data[1:]...)
		PrintMessage(msg, err)
	case "RESTOCK_ITEM":
		if len(data) < 2 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		qty := int32(parseInt(data[1]))
		msg, err := RestockItem(data[0], qty)
		PrintMessage(msg, err)
	case "TRANSACTION_ITEM_RECAP":
		if len(data) < 1 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		transactions, err := GetTransactionItem(data[0])
		PrintTransactionRecap(transactions, err)
	case "TRANSACTION_MEMBER_RECAP":
		if len(data) < 1 {
			PrintMessage("", fmt.Errorf("insufficient arguments"))
			return
		}
		transactions, err := GetTransactionMember(data[0])
		PrintTransactionRecap(transactions, err)
	case "EXIT":
		os.Exit(1)
	default:
		PrintMessage("", fmt.Errorf("unknown command: %s", command))
	}
}

func PrintMessage(successMsg string, errMsg error) {
	if errMsg != nil {
		fmt.Printf("[FAILED] %s\n", errMsg.Error())
	} else {
		fmt.Printf("[SUCCESS] %s\n", successMsg)
	}
}

func PrintTransactionRecap(transactions []Transaction, errMsg error) {
	if errMsg != nil {
		fmt.Printf("[FAILED] %s\n", errMsg.Error())
		return
	}

	fmt.Println("-x-x-x-x-x-x-x-x-x-x-x-x-")
	for _, transaction := range transactions {
		totalPrice := transaction.Qty * transaction.Price
		if transaction.IdMember != nil {
			fmt.Printf("SKU: %s, ID Member: %s, Qty: %d, Total Price: %d\n", transaction.SKU, *transaction.IdMember, transaction.Qty, totalPrice)
		} else {
			fmt.Printf("SKU: %s, ID Member: -, Qty: %d, Total Price: %d\n", transaction.SKU, transaction.Qty, totalPrice)
		}
	}
	fmt.Println("-x-x-x-x-x-x-x-x-x-x-x-x-")
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
