package progress

import (
	"lab03_lab04/core"
)

type Achievement struct{ code, displayName, description string }

var (
	FirstBlood   = &Achievement{"first-blood", "Первая кровь", "одержать первую победу"}
	Lucky        = &Achievement{"lucky", "Счастливчик", "найти зачарованный лут"}
	Veteran      = &Achievement{"veteran", "Ветеран", "достичь 5 уровня"}
	DragonSlayer = &Achievement{"dragon-slayer", "Драконоборец", "победить дракона"}
)

func (a *Achievement) Code() string        { return a.code }
func (a *Achievement) DisplayName() string { return a.displayName }
func (a *Achievement) Description() string { return a.description }
func ByCode(code string) (*Achievement, error) {
	for _, a := range []*Achievement{FirstBlood, Lucky, Veteran, DragonSlayer} {
		if a.code == code {
			return a, nil
		}
	}
	return nil, &core.UnknownCodeError{What: "Достижение", Code: code}
}
