package request

type FetchTowerTopRequest struct {
	TowerID int `json:"tower_id"`
}

type ClearedTowerFloorRequest struct {
	TowerID    int  `json:"tower_id"`
	FloorNo    int  `json:"floor_no"`
	IsAutoMode bool `json:"is_auto_mode"`
}

type RecoveryTowerCardUsedRequest struct {
	TowerID       int   `json:"tower_id"`
	CardMasterIDs []int `json:"card_master_ids"`
}

type RecoveryTowerCardUsedAllRequest = FetchTowerTopRequest
type FetchTowerRankingRequest = FetchTowerTopRequest
