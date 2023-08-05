package serverdb

// card grade up trigger is responsible for showing the pop-up animation when openning a card after getting a new copy
// or right after performing a limit break using items

func (session *Session) AddCardGradeUpTrigger(id int64, trigger any) {
	session.CardGradeUpTriggers = append(session.CardGradeUpTriggers, id)
	session.CardGradeUpTriggers = append(session.CardGradeUpTriggers, trigger)
}

func (session *Session) FinalizeCardGradeUpTrigger() []any {
	// probably need to stores the trigger
	return session.CardGradeUpTriggers
}
