package entity

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(orderId string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order {
	return &Order{
		ID:            orderId,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        "OPEN",
		Transactions:  []*Transaction{},
	}
}

func (o *Order) executeBuyOrder(book *Book, sellOrder *Order) {
	transaction := NewTransaction(sellOrder, o, o.Shares, sellOrder.Price)
	book.AddTransaction(transaction, book.Wg)

	sellOrder.Transactions = append(sellOrder.Transactions, transaction)
	o.Transactions = append(o.Transactions, transaction)

	book.OrdersChannelOut <- sellOrder
	book.OrdersChannelOut <- o
}

func (o *Order) executeSellOrder(book *Book, buyOrder *Order) {
	transaction := NewTransaction(o, buyOrder, o.Shares, buyOrder.Price)
	book.AddTransaction(transaction, book.Wg)

	buyOrder.Transactions = append(buyOrder.Transactions, transaction)
	o.Transactions = append(o.Transactions, transaction)

	book.OrdersChannelOut <- buyOrder
	book.OrdersChannelOut <- o
}
