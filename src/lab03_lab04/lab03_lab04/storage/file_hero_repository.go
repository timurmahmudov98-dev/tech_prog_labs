package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"lab03_lab04/hero"
	"lab03_lab04/item"
	"lab03_lab04/progress"
	"os"
	"path/filepath"
)

type FileHeroRepository struct {
	file    string
	catalog *item.ItemCatalog
}

func NewFileHeroRepository(file string, catalog *item.ItemCatalog) *FileHeroRepository {
	return &FileHeroRepository{file, catalog}
}

type itemSave struct {
	Code         string   `json:"code"`
	Enchantments []string `json:"enchantments"`
}
type heroSave struct {
	Name         string     `json:"name"`
	HeroClass    string     `json:"heroClass"`
	Level        int        `json:"level"`
	Experience   int        `json:"experience"`
	Gold         int        `json:"gold"`
	Weapon       itemSave   `json:"weapon"`
	Armor        itemSave   `json:"armor"`
	Inventory    []itemSave `json:"inventory"`
	Achievements []string   `json:"achievements"`
}

func (r *FileHeroRepository) FindByName(name string) (*hero.Hero, error) {
	list, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	for _, h := range list {
		if h.Name() == name {
			return h, nil
		}
	}
	return nil, nil
}
func (r *FileHeroRepository) FindAll() ([]*hero.Hero, error) {
	data, err := os.ReadFile(r.file)
	if errors.Is(err, os.ErrNotExist) {
		return []*hero.Hero{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []*hero.Hero{}, nil
	}
	var saves []heroSave
	if err = json.Unmarshal(data, &saves); err != nil {
		return nil, err
	}
	if saves == nil {
		return nil, fmt.Errorf("в файле должен быть массив героев")
	}
	result := make([]*hero.Hero, 0, len(saves))
	for _, save := range saves {
		h, err := r.restore(save)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}
func (r *FileHeroRepository) Save(h *hero.Hero) error {
	heroes, err := r.FindAll()
	if err != nil {
		return err
	}
	saves := []heroSave{}
	for _, other := range heroes {
		if other.Name() != h.Name() {
			saves = append(saves, snapshot(other))
		}
	}
	saves = append(saves, snapshot(h))
	data, err := json.MarshalIndent(saves, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(r.file), 0755); err != nil {
		return err
	}
	return os.WriteFile(r.file, append(data, '\n'), 0600)
}
func snapshotItem(it item.Item) itemSave {
	codes := []string{}
	for _, e := range it.Enchantments() {
		codes = append(codes, e.Code())
	}
	return itemSave{it.Code(), codes}
}
func snapshot(h *hero.Hero) heroSave {
	inv := []itemSave{}
	for _, it := range h.Inventory() {
		inv = append(inv, snapshotItem(it))
	}
	ach := []string{}
	for _, a := range h.Achievements() {
		ach = append(ach, a.Code())
	}
	return heroSave{h.Name(), h.HeroClass().Code(), h.Level(), h.Experience(), h.Gold(), snapshotItem(h.Weapon()), snapshotItem(h.Armor()), inv, ach}
}
func (r *FileHeroRepository) restoreItem(s itemSave) (item.Item, error) {
	it, err := r.catalog.Create(s.Code)
	if err != nil {
		return nil, err
	}
	for _, code := range s.Enchantments {
		e, err := item.EnchantmentByCode(code)
		if err != nil {
			return nil, err
		}
		if err = it.Enchant(e); err != nil {
			return nil, err
		}
	}
	return it, nil
}
func (r *FileHeroRepository) restore(s heroSave) (*hero.Hero, error) {
	hc, err := hero.HeroClassByCode(s.HeroClass)
	if err != nil {
		return nil, err
	}
	w, err := r.restoreItem(s.Weapon)
	if err != nil {
		return nil, err
	}
	a, err := r.restoreItem(s.Armor)
	if err != nil {
		return nil, err
	}
	weapon, ok := w.(item.Weapon)
	if !ok {
		return nil, fmt.Errorf("сохраненный предмет не оружие")
	}
	armor, ok := a.(*item.Armor)
	if !ok {
		return nil, fmt.Errorf("сохраненный предмет не броня")
	}
	inv := []item.Item{}
	for _, save := range s.Inventory {
		it, err := r.restoreItem(save)
		if err != nil {
			return nil, err
		}
		inv = append(inv, it)
	}
	ach := []*progress.Achievement{}
	for _, code := range s.Achievements {
		v, err := progress.ByCode(code)
		if err != nil {
			return nil, err
		}
		ach = append(ach, v)
	}
	return hero.NewBuilder().SetName(s.Name).SetHeroClass(hc).SetLevel(s.Level).SetExperience(s.Experience).SetGold(s.Gold).SetWeapon(weapon).SetArmor(armor).SetInventory(inv).SetAchievements(ach).Build()
}
