package world

import (
	"slices"

	"lab03/core"
)

// LocationRegistry - все локации игры
type LocationRegistry struct {
	locations []Location
}

func NewLocationRegistry(locations ...Location) *LocationRegistry {
	return &LocationRegistry{locations: locations}
}

func (r *LocationRegistry) Get(code string) (Location, error) {
	for _, location := range r.locations {
		if location.Code() == code {
			return location, nil
		}
	}

	return nil, &core.UnknownCodeError{What: "Локация", Code: code}
}

func (r *LocationRegistry) All() []Location {
	return slices.Clone(r.locations)
}
