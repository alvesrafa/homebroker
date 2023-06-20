package entity

type Asset struct {
	ID           string
	Name         string
	MarketVloume int
}

func NewAsset(id string, name string, marketVloume int) *Asset {
	return &Asset{
		ID:           id,
		Name:         name,
		MarketVloume: marketVloume,
	}
}
