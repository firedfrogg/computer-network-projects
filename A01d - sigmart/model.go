package main

type Tool interface {
	AddTransaction(data any)
	GetData() any
}

type Transaction struct {
	IdMember *string
	SKU      string
	Qty      int32
	Price    int32
}

type Member struct {
	IdMember     string
	MemberName   string
	Transactions []Transaction
}

func (m *Member) AddTransaction(data any) {
	transaction := data.(Transaction)
	m.Transactions = append(m.Transactions, transaction)
}

func (m *Member) GetData() any {
	return m.Transactions
}

type Item struct {
	SKU          string
	ItemName     string
	StockQty     int32
	Transactions []Transaction
	Price        int32
}

func (it *Item) AddTransaction(data any) {
	transaction := data.(Transaction)
	it.Transactions = append(it.Transactions, transaction)
}

func (it *Item) GetData() any {
	return it.Transactions
}
