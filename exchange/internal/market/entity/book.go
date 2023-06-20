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
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	// loop infinito pq o tempo todo pode ficar caindo orders aqui
	for order := range b.OrdersChannel {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)

			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {

				sellOrder := sellOrders.Pop().(*Order)

				if sellOrder.PedingShares > 0 {

					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)

					b.AddTransaction(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					b.OrdersChannelOut <- sellOrder
					b.OrdersChannelOut <- order

					if sellOrder.PedingShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}

			}
		} else if order.OrderType == "SELL" {
			sellOrders.Push(order)

			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {

				buyOrder := buyOrders.Pop().(*Order)

				if buyOrder.PedingShares > 0 {

					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)

					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					b.OrdersChannelOut <- buyOrder
					b.OrdersChannelOut <- order

					if buyOrder.PedingShares > 0 {
						buyOrders.Push(buyOrder)
					}
				}

			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done() // defer -> Everything under this line will be executed and after that, this line will run

	sellingShares := transaction.SellingOrder.PedingShares
	buyingShares := transaction.BuyingOrder.PedingShares

	minShares := sellingShares
	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.updateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellingOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.updateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyingOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.updateStatus()

	b.Transactions = append(b.Transactions, transaction)
}
