package combat

// BattleEvent - одно событие боя: в каком раунде, что произошло, кто кого,
// на сколько и сколько здоровья осталось у цели
type BattleEvent struct {
	Round             int
	Type              EventType
	Source            Side
	Target            Side
	Amount            int
	TargetHealthAfter int
}
