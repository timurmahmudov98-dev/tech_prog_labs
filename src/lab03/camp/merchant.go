package camp

import (
	"slices"

	"lab03/core"
	"lab03/hero"
	"lab03/item"
)

const sellPricePercent = 50

// Merchant - торговец в лагере
type Merchant struct {
	catalog    *item.ItemCatalog
	assortment []string
}

func NewMerchant(catalog *item.ItemCatalog, assortment ...string) *Merchant {
	return &Merchant{catalog: catalog, assortment: assortment}
}

// Showcase возвращает образцы товаров, чтобы показать их покупателю
func (m *Merchant) Showcase() ([]item.Item, error) {
	showcase := make([]item.Item, 0, len(m.assortment))
	for _, code := range m.assortment {
		it, err := m.catalog.Create(code)
		if err != nil {
			return nil, err
		}
		showcase = append(showcase, it)
	}

	return showcase, nil
}

func (m *Merchant) Buy(h *hero.Hero, code string) (item.Item, error) {
	if !slices.Contains(m.assortment, code) {
		return nil, &core.UnknownCodeError{What: "Товар", Code: code}
	}

	// покупатель получает собственную копию предмета из каталога
	it, err := m.catalog.Create(code)
	if err != nil {
		return nil, err
	}
	if err := h.SpendGold(it.Price()); err != nil {
		return nil, err
	}

	h.AddItem(it)
	return it, nil
}

func (m *Merchant) Sell(h *hero.Hero, it item.Item) (int, error) {
	price := m.SellPrice(it)
	if err := h.RemoveItem(it); err != nil {
		return 0, err
	}

	h.AddGold(price)
	return price, nil
}

func (m *Merchant) SellPrice(it item.Item) int {
	return it.Price() * sellPricePercent / 100
}
