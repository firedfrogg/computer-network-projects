package main

import (
	"errors"
	"fmt"
)

var (
	Name      string   = "Julian Alex Joshua" // please insert your name here
	IdStudent string   = "2206082606"         // please insert your id student here
	Items     []Item                          // contain array of item pointer
	Members   []Member                        // contain array of member pointer
)

func AddItem(SKU string, itemName string, price int32, stockQty int32) (string, error) {
	for _, item := range Items {
		if item.SKU == SKU {
			return "", fmt.Errorf("item %s is already in list of items", SKU)
		}
	}
	newItem := Item{
		SKU:      SKU,
		ItemName: itemName,
		Price:    price,
		StockQty: stockQty,
	}
	Items = append(Items, newItem)
	return fmt.Sprintf("successfully added item %s to list of items", SKU), nil
}

func DeleteItem(SKU string) (string, error) {
	for i, item := range Items {
		if item.SKU == SKU {
			if len(item.Transactions) > 0 {
				return "", fmt.Errorf("there is at least one transaction taking item %s", SKU)
			}
			Items = append(Items[:i], Items[i+1:]...)
			return fmt.Sprintf("successfully deleted item %s from list of items", SKU), nil
		}
	}
	return "", fmt.Errorf("item %s is not in list of items", SKU)
}

func AddMember(idMember string, memberName string) (string, error) {
	for _, member := range Members {
		if member.IdMember == idMember {
			return "", fmt.Errorf("member %s is already in list of members", idMember)
		}
	}
	newMember := Member{
		IdMember:   idMember,
		MemberName: memberName,
	}
	Members = append(Members, newMember)
	return fmt.Sprintf("successfully added member %s to list of members", idMember), nil
}

func DeleteMember(idMember string) (string, error) {
	for i, member := range Members {
		if member.IdMember == idMember {
			if len(member.Transactions) > 0 {
				return "", fmt.Errorf("there is at least one transaction taking member %s", idMember)
			}
			Members = append(Members[:i], Members[i+1:]...)
			return fmt.Sprintf("successfully deleted member %s from list of members", idMember), nil
		}
	}
	return "", fmt.Errorf("member %s is not in list of members", idMember)
}

func AddTransaction(qty int32, data ...string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no SKU provided")
	}

	SKU := data[0]
	var idMember *string
	if len(data) > 1 {
		idMember = &data[1]
	}

	var item *Item
	for i := range Items {
		if Items[i].SKU == SKU {
			item = &Items[i]
			break
		}
	}
	if item == nil {
		return "", fmt.Errorf("item %s is not in list of items", SKU)
	}

	if idMember != nil {
		var member *Member
		for i := range Members {
			if Members[i].IdMember == *idMember {
				member = &Members[i]
				break
			}
		}
		if member == nil {
			return "", fmt.Errorf("member %s is not in list of members", *idMember)
		}
	}

	if item.StockQty < qty {
		return "", fmt.Errorf("stock qty for item %s is not sufficient", SKU)
	}

	item.StockQty -= qty
	transaction := Transaction{
		IdMember: idMember,
		SKU:      SKU,
		Qty:      qty,
		Price:    item.Price,
	}
	item.Transactions = append(item.Transactions, transaction)

	if idMember != nil {
		for i := range Members {
			if Members[i].IdMember == *idMember {
				Members[i].Transactions = append(Members[i].Transactions, transaction)
				break
			}
		}
		return fmt.Sprintf("successfully added transaction item %s for member %s", SKU, *idMember), nil
	}

	return fmt.Sprintf("successfully added transaction item %s", SKU), nil
}

func RestockItem(SKU string, qty int32) (string, error) {
	for i := range Items {
		if Items[i].SKU == SKU {
			Items[i].StockQty += qty
			return fmt.Sprintf("successfully restock qty for item %s", SKU), nil
		}
	}
	return "", fmt.Errorf("item %s is not in list of items", SKU)
}

func GetTransactionItem(SKU string) ([]Transaction, error) {
	for _, item := range Items {
		if item.SKU == SKU {
			return item.Transactions, nil
		}
	}
	return nil, fmt.Errorf("item %s is not in list of items", SKU)
}

func GetTransactionMember(idMember string) ([]Transaction, error) {
	for _, member := range Members {
		if member.IdMember == idMember {
			return member.Transactions, nil
		}
	}
	return nil, fmt.Errorf("member %s is not in list of members", idMember)
}
