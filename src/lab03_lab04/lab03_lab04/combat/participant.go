package combat

// Participant - участник боя в том виде, в котором он попадает в отчет
type Participant struct {
	Code      string
	Name      string
	MaxHealth int
}

func ParticipantOf(combatant Combatant) Participant {
	return Participant{Code: combatant.Code(), Name: combatant.Name(), MaxHealth: combatant.MaxHealth()}
}
