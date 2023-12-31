package model

import (
	"elichika/utils"

	"encoding/json"
	"fmt"
)

type LiveStage struct {
	LiveDifficultyID int               `json:"live_difficulty_id"`
	LiveNotes        []LiveNote        `json:"live_notes"`
	LiveWaveSettings []LiveWaveSetting `json:"live_wave_settings"`
	NoteGimmicks     []NoteGimmick     `json:"note_gimmicks"`
	StageGimmickDict []any             `json:"stage_gimmick_dict"`
}

func (this *LiveStage) Copy() LiveStage {
	result := LiveStage{
		LiveDifficultyID: this.LiveDifficultyID,
		LiveNotes:        []LiveNote{},
		LiveWaveSettings: []LiveWaveSetting{},
		NoteGimmicks:     []NoteGimmick{},
		StageGimmickDict: []any{},
	}
	result.LiveNotes = append(result.LiveNotes, this.LiveNotes...)
	result.LiveWaveSettings = append(result.LiveWaveSettings, this.LiveWaveSettings...)
	result.NoteGimmicks = append(result.NoteGimmicks, this.NoteGimmicks...)
	result.StageGimmickDict = append(result.StageGimmickDict, this.StageGimmickDict...)
	return result
}

func (this *LiveStage) IsSame(other *LiveStage) bool {
	if this.LiveDifficultyID != other.LiveDifficultyID {
		return false
	}
	if len(this.LiveNotes) != len(other.LiveNotes) {
		return false
	}
	for i := range this.LiveNotes {
		if !this.LiveNotes[i].IsSame(&other.LiveNotes[i]) {
			fmt.Println(this.LiveNotes[i])
			fmt.Println(other.LiveNotes[i])
			return false
		}
	}
	// fmt.Println("Notes OK")
	if len(this.LiveWaveSettings) != len(other.LiveWaveSettings) {
		return false
	}
	for i := range this.LiveWaveSettings {
		if this.LiveWaveSettings[i] != other.LiveWaveSettings[i] {
			fmt.Println(this.LiveWaveSettings[i])
			fmt.Println(other.LiveWaveSettings[i])
			return false
		}
	}
	// fmt.Println("Waves OK")
	if len(this.NoteGimmicks) != len(other.NoteGimmicks) {
		return false
	}
	for i := range this.NoteGimmicks {
		if !this.NoteGimmicks[i].IsSame(&other.NoteGimmicks[i]) {
			return false
		}
	}
	// fmt.Println("Note Gimmicks OK")
	if len(this.StageGimmickDict) != len(other.StageGimmickDict) {
		return false
	}
	if len(this.StageGimmickDict) > 0 {
		thisDict, err := json.Marshal(this.StageGimmickDict[0])
		utils.CheckErr(err)
		thisID := 0
		err = json.Unmarshal(thisDict, &thisID)
		utils.CheckErr(err)
		if thisID != other.StageGimmickDict[0].(int) {
			return false
		}

		thisDict, err = json.Marshal(this.StageGimmickDict[1].([]any)[0])
		utils.CheckErr(err)
		thisGimmick := StageGimmick{}
		err = json.Unmarshal(thisDict, &thisGimmick)
		utils.CheckErr(err)
		otherGimmick := other.StageGimmickDict[1].([]StageGimmick)[0]
		if thisGimmick != otherGimmick {
			fmt.Println(thisGimmick)
			fmt.Println(otherGimmick)
			return false
		}
	}
	// fmt.Println("Stage Gimmick OK")
	return true
}

type LiveNote struct {
	ID                  int `json:"id"`
	CallTime            int `json:"call_time"`
	NoteType            int `json:"note_type"`
	NotePosition        int `json:"note_position"`
	GimmickID           int `json:"gimmick_id"`
	NoteAction          int `json:"note_action"`
	WaveID              int `json:"wave_id"`
	NoteRandomDropColor int `json:"note_random_drop_color"`
	AutoJudgeType       int `json:"auto_judge_type"`
}

func (this *LiveNote) IsSame(other *LiveNote) bool {
	same := true
	same = same && (this.ID == other.ID)
	same = same && (this.CallTime == other.CallTime)
	same = same && (this.NoteType == other.NoteType)
	same = same && (this.NotePosition == other.NotePosition)
	same = same && (this.GimmickID == other.GimmickID)
	same = same && (this.NoteAction == other.NoteAction)
	same = same && (this.WaveID == other.WaveID)
	return same
}

type LiveWaveSetting struct {
	ID            int `json:"id"`
	WaveDamage    int `json:"wave_damage"`
	MissionType   int `json:"mission_type"`
	Arg1          int `json:"arg_1"`
	Arg2          int `json:"arg_2"`
	RewardVoltage int `json:"reward_voltage"`
}

type NoteGimmick struct {
	UniqID          int `json:"uniq_id"`
	ID              int `json:"id"`
	NoteGimmickType int `json:"note_gimmick_type"`
	Arg1            int `json:"arg_1"`
	Arg2            int `json:"arg_2"`
	EffectMID       int `json:"effect_m_id"`
	IconType        int `json:"icon_type"`
}

func (this *NoteGimmick) IsSame(other *NoteGimmick) bool {
	same := true
	same = same && (this.UniqID == other.UniqID)
	same = same && (this.NoteGimmickType == other.NoteGimmickType)
	same = same && (this.Arg1 == other.Arg1)
	same = same && (this.Arg2 == other.Arg2)
	same = same && (this.EffectMID == other.EffectMID)
	if !same {
		return false
	}
	if this.IconType == other.IconType {
		return true
	}
	if this.IconType == 5 && other.IconType == 25 { // there was a db update that change this
		return true
	}
	if this.IconType == 8 && other.IconType == 9 { // there was a db update that change this
		return true
	}
	return false
}

type StageGimmick struct {
	GimmickMasterID    int `json:"gimmick_master_id"`
	ConditionMasterID1 int `json:"condition_master_id_1"`
	ConditionMasterID2 int `json:"condition_master_id_2"`
	SkillMasterID      int `json:"skill_master_id"`
	UniqID             int `json:"uniq_id"`
}
