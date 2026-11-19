package tactic

type Action struct{ AttackMultiplier, IncomingMultiplier float64 }

var Attack = Action{1, 1}
var PowerStrike = Action{1.5, 1.5}
var Defend = Action{0.5, 0.3}
