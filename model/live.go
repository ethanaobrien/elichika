package model

import (
	"encoding/json"
	"reflect"
)

// LiveDaily ...
type LiveDaily struct {
	LiveDailyMasterID      int `json:"live_daily_master_id" xorm:"id"`
	LiveMasterID           int `json:"live_master_id" xorm:"live_id"`
	EndAt                  int `json:"end_at"`
	RemainingPlayCount     int `json:"remaining_play_count"`
	RemainingRecoveryCount int `json:"remaining_recovery_count"`
}

// LivePartnerInfo ...
type LivePartnerInfo UserBasicInfo

// guests before live start
type LiveStartLivePartner struct {
	UserID int `xorm:"'user_id' "json:"user_id"`
	Name   struct {
		DotUnderText string `xorm:"'name'" json:"dot_under_text"`
	} `xorm:"extends" json:"name"`
	Rank                int   `json:"rank"`
	LastLoginAt         int64 `json:"last_login_at"`
	CardByCategory      []any `xorm:"-" json:"card_by_category"`
	EmblemID            int   `xorm:"'emblem_id'" json:"emblem_id"`
	IsFriend            bool  `xorm:"-" json:"is_friend"`
	IntroductionMessage struct {
		DotUnderText string `xorm:"'message'" json:"dot_under_text"`
	} `xorm:"extends" json:"introduction_message"`
}

// the live being played
type UserLive struct {
	UserID          int             `xorm:"pk 'user_id'" json:"-"`
	LiveID          int64           `xorm:"'live_id'" json:"live_id"`
	LiveType        int             `xorm:"'live_type'" json:"live_type"`
	DeckID          int             `xorm:"'deck_id'" json:"deck_id"`
	LiveStage       LiveStage       `xorm:"-" json:"live_stage"`
	PartnerUserID   int             `xorm:"'partner_user_id'" json:"-"`
	IsAutoplay      bool            `xorm:"'is_autoplay'" json:"-"`
	LivePartnerCard PartnerCardInfo `xorm:"extends" json:"live_partner_card"`
	IsPartnerFriend bool            `xorm:"'is_partner_friend'" json:"is_partner_friend"`
	CellID          *int            `xorm:"'cell_id' "json:"cell_id"`
	TowerLive       TowerLive       `xorm:"extends" json:"tower_live"`
}

func (this UserLive) MarshalJSON() ([]byte, error) {
	bytes := []byte("{")
	rType := reflect.TypeOf(this)
	isFirst := true
	for i := 0; i < rType.NumField(); i++ {
		rField := rType.Field(i)
		if (rField.Name == "TowerLive") && (this.TowerLive.TowerID == nil) {
			continue
		}
		key := rField.Tag.Get("json")
		if key == "-" {
			continue
		} else if key == "" {
			panic("empty key")
		}
		if isFirst {
			isFirst = false
		} else {
			bytes = append(bytes, []byte(",")...)
		}
		bytes = append(bytes, []byte("\"")...)
		bytes = append(bytes, []byte(key)...)
		bytes = append(bytes, []byte("\":")...)
		if (rField.Name == "LivePartnerCard") && (this.LivePartnerCard.CardMasterID == 0) {
			bytes = append(bytes, []byte("null")...)
			continue
		}
		fieldBytes, err := json.Marshal(reflect.ValueOf(this).Field(i).Interface())
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, fieldBytes...)
	}

	bytes = append(bytes, []byte("}")...)
	return bytes, nil
}

// MemberLovePanels ...
type MemberLovePanels struct {
	MemberID               int   `json:"member_id"`
	MemberLovePanelCellIds []int `json:"member_love_panel_cell_ids"`
}

type LiveUpdatePlayListReq struct {
	LiveMasterID int  `json:"live_master_id"`
	GroupNum     int  `json:"group_num"`
	IsSet        bool `json:"is_set"`
}

type UserPlayListItem struct {
	UserID         int  `xorm:"pk 'user_id'" json:"-"`
	UserPlayListID int  `xorm:"pk 'user_play_list_id'" json:"user_play_list_id"`
	GroupNum       int  `xorm:"'group_num'" json:"group_num"` // UserPlayListID % 10
	LiveID         int  `xorm:"'live_id'" json:"live_id"`     // UserPlayListID / 10
	IsNull         bool `xorm:"-" json:"-"`
}

func (item UserPlayListItem) ID() int64 {
	return int64(item.UserPlayListID)
}

func init() {
	if TableNameToInterface == nil {
		TableNameToInterface = make(map[string]interface{})
	}
	TableNameToInterface["u_live_deck"] = UserLiveDeck{}
	TableNameToInterface["u_live_party"] = UserLiveParty{}
	// TableNameToInterface["u_live_state"] = LiveState{}
	TableNameToInterface["u_play_list"] = UserPlayListItem{}
}
