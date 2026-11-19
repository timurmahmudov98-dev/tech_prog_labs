package storage

import (
	"lab03_lab04/hero"
	"slices"
)

type CachingHeroRepository struct {
	target HeroRepository
	cache  []*hero.Hero
	loaded bool
}

func NewCachingHeroRepository(target HeroRepository) *CachingHeroRepository {
	return &CachingHeroRepository{target: target}
}
func (r *CachingHeroRepository) FindAll() ([]*hero.Hero, error) {
	if !r.loaded {
		list, err := r.target.FindAll()
		if err != nil {
			return nil, err
		}
		r.cache = list
		r.loaded = true
	}
	return slices.Clone(r.cache), nil
}
func (r *CachingHeroRepository) FindByName(name string) (*hero.Hero, error) {
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
func (r *CachingHeroRepository) Save(h *hero.Hero) error {
	if err := r.target.Save(h); err != nil {
		return err
	}
	if _, err := r.FindAll(); err != nil {
		return err
	}
	for i, other := range r.cache {
		if other.Name() == h.Name() {
			r.cache[i] = h
			return nil
		}
	}
	r.cache = append(r.cache, h)
	return nil
}
