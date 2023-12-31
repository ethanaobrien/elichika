package model

type UserLiveDeck struct {
	UserID         int `xorm:"pk 'user_id'" json:"-"`
	UserLiveDeckID int `xorm:"pk 'user_live_deck_id'" json:"user_live_deck_id"`
	Name           struct {
		DotUnderText string `xorm:"name" json:"dot_under_text"`
	} `xorm:"extends" json:"name"` // deck name
	CardMasterID1 int `xorm:"'card_master_id_1'" json:"card_master_id_1"`
	CardMasterID2 int `xorm:"'card_master_id_2'" json:"card_master_id_2"`
	CardMasterID3 int `xorm:"'card_master_id_3'" json:"card_master_id_3"`
	CardMasterID4 int `xorm:"'card_master_id_4'" json:"card_master_id_4"`
	CardMasterID5 int `xorm:"'card_master_id_5'" json:"card_master_id_5"`
	CardMasterID6 int `xorm:"'card_master_id_6'" json:"card_master_id_6"`
	CardMasterID7 int `xorm:"'card_master_id_7'" json:"card_master_id_7"`
	CardMasterID8 int `xorm:"'card_master_id_8'" json:"card_master_id_8"`
	CardMasterID9 int `xorm:"'card_master_id_9'" json:"card_master_id_9"`
	SuitMasterID1 int `xorm:"'suit_master_id_1'" json:"suit_master_id_1"`
	SuitMasterID2 int `xorm:"'suit_master_id_2'" json:"suit_master_id_2"`
	SuitMasterID3 int `xorm:"'suit_master_id_3'" json:"suit_master_id_3"`
	SuitMasterID4 int `xorm:"'suit_master_id_4'" json:"suit_master_id_4"`
	SuitMasterID5 int `xorm:"'suit_master_id_5'" json:"suit_master_id_5"`
	SuitMasterID6 int `xorm:"'suit_master_id_6'" json:"suit_master_id_6"`
	SuitMasterID7 int `xorm:"'suit_master_id_7'" json:"suit_master_id_7"`
	SuitMasterID8 int `xorm:"'suit_master_id_8'" json:"suit_master_id_8"`
	SuitMasterID9 int `xorm:"'suit_master_id_9'" json:"suit_master_id_9"`
}

func (uld *UserLiveDeck) ID() int64 {
	return int64(uld.UserLiveDeckID)
}

type UserLiveParty struct {
	UserID         int `xorm:"pk 'user_id'" json:"-"`
	PartyID        int `xorm:"pk 'party_id'" json:"party_id"`
	UserLiveDeckID int `xorm:"'user_live_deck_id'" json:"user_live_deck_id"`
	Name           struct {
		DotUnderText string `xorm:"name" json:"dot_under_text"`
	} `xorm:"extends" json:"name"` // deck name
	IconMasterID     int    `xorm:"'icon_master_id'" json:"icon_master_id"`
	CardMasterID1    int    `xorm:"'card_master_id_1'" json:"card_master_id_1"`
	CardMasterID2    int    `xorm:"'card_master_id_2'" json:"card_master_id_2"`
	CardMasterID3    int    `xorm:"'card_master_id_3'" json:"card_master_id_3"`
	UserAccessoryID1 *int64 `xorm:"'user_accessory_id_1'" json:"user_accessory_id_1"` // null for empty
	UserAccessoryID2 *int64 `xorm:"'user_accessory_id_2'" json:"user_accessory_id_2"`
	UserAccessoryID3 *int64 `xorm:"'user_accessory_id_3'" json:"user_accessory_id_3"`
}

func (uld *UserLiveParty) ID() int64 {
	return int64(uld.PartyID)
}

// DeckSquadDict ...
type DeckSquadDict struct {
	CardMasterIDs    []int    `json:"card_master_ids"`
	UserAccessoryIDs []*int64 `json:"user_accessory_ids"`
}
