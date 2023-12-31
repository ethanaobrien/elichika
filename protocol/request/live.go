package request

import (
	"elichika/generic"
	"elichika/model"
)

type LiveStartRequest struct {
	LiveDifficultyID    int                    `json:"live_difficulty_id"`
	DeckID              int                    `json:"deck_id"`
	CellID              *int                   `json:"cell_id"`
	PartnerUserID       int                    `json:"partner_user_id"`
	PartnerCardMasterID int                    `json:"partner_card_master_id"`
	LpMagnification     int                    `json:"lp_magnification"`
	IsAutoPlay          bool                   `json:"is_auto_play"`
	IsReferenceBook     bool                   `json:"is_reference_book"`
	LiveTowerStatus     *model.LiveTowerStatus `json:"live_tower_status"`
}

type LiveFinishCard struct {
	CardMasterID        int `json:"-"`
	GotVoltage          int `json:"got_voltage"`
	SkillTriggeredCount int `json:"skill_triggered_count"`
	AppealCount         int `json:"appeal_count"`
}

func (obj *LiveFinishCard) SetID(id int64) {
	obj.CardMasterID = int(id)
}

type LiveFinishRequest struct {
	LiveID           int64 `json:"live_id"`
	LiveFinishStatus int   `json:"live_finish_status"`
	LiveScore        struct {
		StartInfo                  any                                                `json:"start_info"`
		FinishInfo                 any                                                `json:"finish_info"`
		ResultDict                 []any                                              `json:"result_dict"`
		WaveStatDict               []any                                              `json:"wave_stat_dict"`
		TurnStatDict               []any                                              `json:"turn_stat_dict"`
		CardStatDict               generic.ObjectByObjectIDList[model.LiveFinishCard] `json:"card_stat_dict"`
		TargetScore                int                                                `json:"target_score"`
		CurrentScore               int                                                `json:"current_score"`
		ComboCount                 int                                                `json:"combo_count"`
		ChangeSquadCount           int                                                `json:"change_squad_count"`
		HighestComboCount          int                                                `json:"highest_combo_count"`
		RemainingStamina           int                                                `json:"remaining_stamina"`
		IsPerfectLive              bool                                               `json:"is_perfect_live"`
		IsPerfectFullCombo         bool                                               `json:"is_perfect_full_combo"`
		UseVoltageActiveSkillCount int                                                `json:"use_voltage_active_skill_count"`
		UseHealActiveSkillCount    int                                                `json:"use_heal_active_skill_count"`
		UseDebufActiveSkillCount   int                                                `json:"use_debuf_active_skill_count"`
		UseBufActiveSkillCount     int                                                `json:"use_buf_active_skill_count"`
		UseSpSkillCount            int                                                `json:"use_sp_skill_count"`
		CompleteAppealChanceCount  int                                                `json:"complete_appeal_chance_count"`
		TriggerCriticalCount       int                                                `json:"triggered_critical_count"`
		LivePower                  int                                                `json:"live_power"`
		SpSkillScoreList           []int                                              `json:"sp_skill_score_list"`
	} `json:"live_score"`
	ResumeFinishInfo any `json:"resume_finish_info"`
	RoomID           int `json:"room_id"`
}

type FinishTutorialRequest struct {
}
