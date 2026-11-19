package event

import (
	"lab03_lab04/hero"
	"lab03_lab04/world"
)

type GameEvent struct {
	Type       GameEventType
	Hero       *hero.Hero
	Expedition world.ExpeditionReport
}
