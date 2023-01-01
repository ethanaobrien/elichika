package gamedata

import (
	"elichika/dictionary"
	"elichika/utils"

	"fmt"

	"xorm.io/xorm"
)

type LiveMemberMapping = map[int]LiveMemberMappingMember
type LiveMemberMappingMember struct {
	// from m_live_member_mapping and m_live_override_member_mapping
	// - mapping ID is repeated for the whole mapping
	// - MemberMasterID marked as pk to help loading for Live
	MappingID *int `xorm:"pk 'mapping_id'"`
	Position  int  `xorm:"pk 'position'"`
	// MemberMasterID *int  `xorm:"'member_master_id'"`
	// MemberNonPlayableMasterID *int  `xorm:"'member_non_playable_master_id'"`
	IsCenter     bool `xorm:"'is_center'"`
	CardPosition int  `xorm:"'card_position'"`
	// SuitMasterID   *int `xorm:"'suit_master_id'"`
}

func loadLiveMemberMapping(gamedata *Gamedata, masterdata_db, serverdata_db *xorm.Session, dictionary *dictionary.Dictionary) {
	fmt.Println("Loading LiveMemberMapping")
	gamedata.LiveMemberMapping = make(map[int]LiveMemberMapping)

	tables := []string{"m_live_member_mapping", "m_live_override_member_mapping"}
	for _, table := range tables {
		memberMappings := []LiveMemberMappingMember{}
		err := masterdata_db.Table(table).Find(&memberMappings)
		utils.CheckErr(err)
		for _, memberMapping := range memberMappings {
			_, exist := gamedata.LiveMemberMapping[*memberMapping.MappingID]
			if !exist {
				gamedata.LiveMemberMapping[*memberMapping.MappingID] = make(LiveMemberMapping)
			}
			gamedata.LiveMemberMapping[*memberMapping.MappingID][memberMapping.Position] = memberMapping
		}
	}

}

func init() {
	addLoadFunc(loadLiveMemberMapping)
}
