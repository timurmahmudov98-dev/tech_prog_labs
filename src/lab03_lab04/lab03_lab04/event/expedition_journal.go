package event

import (
	"lab03_lab04/world"
	"slices"
)

type ExpeditionJournal struct {
	entries map[string][]world.ExpeditionReport
}

func (j *ExpeditionJournal) OnEvent(e GameEvent) {
	if j.entries == nil {
		j.entries = make(map[string][]world.ExpeditionReport)
	}
	j.entries[e.Hero.Name()] = append(j.entries[e.Hero.Name()], e.Expedition)
}
func (j *ExpeditionJournal) Entries(name string) []world.ExpeditionReport {
	return slices.Clone(j.entries[name])
}
