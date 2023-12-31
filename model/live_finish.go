package model

type LiveFinishCard struct {
	CardMasterID        int `json:"-"`
	GotVoltage          int `json:"got_voltage"`
	SkillTriggeredCount int `json:"skill_triggered_count"`
	AppealCount         int `json:"appeal_count"`
}

func (obj *LiveFinishCard) SetID(id int64) {
	obj.CardMasterID = int(id)
}
