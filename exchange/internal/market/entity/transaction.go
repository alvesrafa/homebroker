package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price

	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (t *Transaction) CalculateTotal(shares int, price float64) {

	t.Total = float64(shares) * price

}

func (t *Transaction) UpdateInvestorAssetPosition() {
	minShares := t.GetMinShares()

	t.SellingOrder.Investor.updateAssetPosition(t.SellingOrder.Asset.ID, -minShares)
	t.AddSellingOrderPendingShares(-minShares)

	t.BuyingOrder.Investor.updateAssetPosition(t.BuyingOrder.Asset.ID, minShares)
	t.AddBuyingOrderPendingShares(-minShares)
}

func (t *Transaction) GetMinShares() int {
	sellingShares := t.SellingOrder.PendingShares
	buyingShares := t.BuyingOrder.PendingShares
	minShares := sellingShares
	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	return minShares
}

func (t *Transaction) CloseBuyingOrder() {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}
}
func (t *Transaction) CloseSellingOrder() {
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}
func (t *Transaction) updateStatus() {
	t.CloseBuyingOrder()
	t.CloseSellingOrder()
}

func (t *Transaction) AddBuyingOrderPendingShares(shares int) {
	t.BuyingOrder.PendingShares += shares
}
func (t *Transaction) AddSellingOrderPendingShares(shares int) {
	t.SellingOrder.PendingShares += shares
}
