package model

type LiveMVDeck struct {
	// unused, LiveMVDeck is not saved in server db
	UserID              int   `xorm:"pk 'user_id'" json:"-"`
	LiveMasterID        int   `xorm:"pk 'live_master_id'" json:"live_master_id"`
	LiveMvDeckType      int   `json:"live_mv_deck_type"` // 1 for original deck, 2 for custom deck
	MemberMasterIDByPos []int `xorm:"'member_master_id_by_pos'" json:"member_master_id_by_pos"`
	SuitMasterIDByPos   []int `xorm:"'suit_master_id_by_pos'" json:"suit_master_id_by_pos"`
	ViewStatusByPos     []int `json:"view_status_by_pos"`
}
