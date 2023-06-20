package entity

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}

type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}
func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	// Investor.AssetPosition.append(assetPosition) no js
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  shares,
	}
}

func (i *Investor) updateAssetPosition(assetID string, qntShares int) {
	assetPosition := i.GetAssetPosition(assetID)

	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, qntShares))
	} else {
		assetPosition.Shares += qntShares
	}

}
func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {

	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}

	return nil
}
