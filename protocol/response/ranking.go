package response

type RankingUser struct {
	UserID                 int           `json:"user_id"`
	UserName               LocalizedText `json:"user_name"`
	UserRank               int           `json:"user_rank"`
	CardMasterID           int           `json:"card_master_id"`
	Level                  int           `json:"level"`
	IsAwakening            bool          `json:"is_awakening"`
	IsAllTrainingActivated bool          `json:"is_all_training_activated"`
	EmblemMasterID         int           `json:"emblem_master_id"`
}
