package storage

import (
	"lab03_lab04/hero"
)

type HeroRepository interface {
	FindByName(string) (*hero.Hero, error)
	FindAll() ([]*hero.Hero, error)
	Save(*hero.Hero) error
}
