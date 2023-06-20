package entity

type OrderQueue struct {
	Orders []*Order
}

// Less - is less?
func (oq *OrderQueue) Less(i, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}

// Swap - reverse positions
func (oq *OrderQueue) Swap(i, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

// Len - length
func (oq *OrderQueue) Len() int {
	return len(oq.Orders)
}

// Push - add more
func (oq *OrderQueue) Push(someValue interface{}) {
	oq.Orders = append(oq.Orders, someValue.(*Order))
}

// Pop - remove the last
func (oq *OrderQueue) Pop() interface{} {
	// save orders
	old := oq.Orders

	// get quantity
	quantity := len(old)

	// get last order
	order := old[quantity-1]

	// remove last order
	oq.Orders = old[0 : quantity-1]

	return order
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}
