package gamedata

import (
	"elichika/model"
	"elichika/utils"

	"strconv"
)

type SimpleLiveStage struct {
	LiveDifficultyID int                     `json:"live_difficulty_id"`
	LiveNotes        []model.LiveNote        `json:"live_notes"`
	LiveWaveSettings []model.LiveWaveSetting `json:"live_wave_settings"`
	Original         *int                    `json:"original"`
}

type LiveDifficultyGimmick struct {
	ID                     int `xorm:"pk 'id'"`
	LiveDifficultyMasterID int `xorm:"live_difficulty_master_id"`
	TriggerType            int `xorm:"'trigger_type'"`
	ConditionMasterID1     int `xorm:"'condition_master_id1'"`
	ConditionMasterID2     int `xorm:"'condition_master_id2'"`
	SkillMasterID          int `xorm:"'skill_master_id'"`
}

type LiveDifficultyNoteGimmick struct {
	LiveDifficultyID    int    `xorm:"pk 'live_difficulty_id'"`
	NoteID              int    `xorm:"pk 'note_id'"`
	NoteGimmickType     int    `xorm:"note_gimmick_type"`
	NoteGimmickIconType int    `xorm:"note_gimmick_icon_type"`
	SkillMasterID       int    `xorm:"skill_master_id"`
	Name                string `xorm:"name"`
	ID                  int    `xorm:"-"` // must be extracted from name
}

type LiveNoteWaveGimmickGroup struct {
	LiveDifficultyID int `xorm:"pk 'live_difficulty_id'"`
	WaveID           int `xorm:"pk 'wave_id'"`
	State            int
	SkillID          int
}

func (this *LiveDifficultyNoteGimmick) populate() {
	var err error
	this.ID, err = strconv.Atoi(this.Name[len(this.Name)-8:])
	utils.CheckErr(err)
}
