package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders           []*Order
	Transactions     []*Transaction
	OrdersChannel    chan *Order     // all orders received we will take from this channel // from kafka
	OrdersChannelOut chan *Order     // all orders sent we will send to this channel // to kafka
	Wg               *sync.WaitGroup // help to sync the threads // the wait group will wait all necessary transactions finish
}

func NewBook(orderChan chan *Order, ordersChannelOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders:           []*Order{},
		Transactions:     []*Transaction{},
		OrdersChannel:    orderChan,
		OrdersChannelOut: ordersChannelOut,
		Wg:               wg,
	}
}

func (b *Book) Trade() {
	// buyOrders := NewOrderQueue()
	// sellOrders := NewOrderQueue()
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	// heap.Init(buyOrders)
	// heap.Init(sellOrders)

	// loop infinito pq o tempo todo pode ficar caindo orders aqui
	for order := range b.OrdersChannel {
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		b.AddOrder(order, buyOrders[asset], sellOrders[asset])
	}
}

func (b *Book) AddOrder(order *Order, buyOrders *OrderQueue, sellOrders *OrderQueue) {
	if order.OrderType == "BUY" {
		buyOrders.Push(order)

		if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
			lastSellOrder := sellOrders.Pop().(*Order)

			if lastSellOrder.PendingShares > 0 {
				order.executeBuyOrder(b, lastSellOrder)

				if lastSellOrder.PendingShares > 0 {
					sellOrders.Push(lastSellOrder)
				}
			}
		}
	} else if order.OrderType == "SELL" {
		sellOrders.Push(order)

		if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
			lastBuyOrder := buyOrders.Pop().(*Order)

			if lastBuyOrder.PendingShares > 0 {
				order.executeSellOrder(b, lastBuyOrder)

				if lastBuyOrder.PendingShares > 0 {
					buyOrders.Push(lastBuyOrder)
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done() // defer -> Everything under this line will be executed and after that, this line will run

	transaction.UpdateInvestorAssetPosition()

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.updateStatus()

	b.Transactions = append(b.Transactions, transaction)
}
