package response

import (
	"elichika/model"
)

type FetchTowerSelectResponse struct {
	TowerIDs      []int      `json:"tower_ids"`
	UserModelDiff *UserModel `json:"user_model_diff"`
}

type FetchTowerTopResponse struct {
	TowerCardUsedCountRows []model.UserTowerCardUsedCount `json:"tower_card_used_count_rows"`
	IsShowUnlockEffect     bool                           `json:"is_show_unlock_effect"`
	UserModelDiff          *UserModel                     `json:"user_model_diff"`
	Order                  *int64                         `json:"order"`
	EachBonusLiveVoltage   []int64                        `json:"each_bonus_live_voltage"`
}
