package serverdb

import (
	"elichika/model"

	"fmt"
)

func (session *Session) GetMember(memberMasterID int) model.UserMemberInfo {
	member := model.UserMemberInfo{}
	exists, err := Engine.Table("s_user_member").
		Where("user_id = ? AND member_master_id = ?", session.UserStatus.UserID, memberMasterID).Get(&member)
	// inserted at login if not exist
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("member not found")
	}
	return member
}

func (session *Session) GetAllMembers() []model.UserMemberInfo {
	members := []model.UserMemberInfo{}
	err := Engine.Table("s_user_member").Where("user_id = ?", session.UserStatus.UserID).Find(&members)
	if err != nil {
		panic(err)
	}
	return members
}

func (session *Session) UpdateMember(member model.UserMemberInfo) {
	session.UserMemberDiffs[member.MemberMasterID] = member
}

func (session *Session) InsertMembers(members []model.UserMemberInfo) {
	affected, err := Engine.Table("s_user_member").Insert(&members)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " members")
}

func (session *Session) FinalizeUserMemberDiffs() []any {
	userMemberByMemberID := []any{}
	for memberMasterID, member := range session.UserMemberDiffs {
		userMemberByMemberID = append(userMemberByMemberID, memberMasterID)
		userMemberByMemberID = append(userMemberByMemberID, member)
		affected, err := Engine.Table("s_user_member").
			Where("user_id = ? AND member_master_id = ?", session.UserStatus.UserID, memberMasterID).AllCols().Update(member)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userMemberByMemberID
}

func (session *Session) GetLovePanelCellIDs(memberID int) []int {
	userMemberLovePanel := model.UserMemberLovePanel{}
	userMemberLovePanel.MemberLovePanelCellIDs = make([]int, 0) // return empty if empty
	exists, err := Engine.Table("s_user_member").
		Where("user_id = ? AND member_master_id = ?", session.UserStatus.UserID, memberID).
		Get(&userMemberLovePanel)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("not exists")
	}
	return userMemberLovePanel.MemberLovePanelCellIDs
}

func (session *Session) GetAllMemberLovePanels() []model.UserMemberLovePanel {
	lovePanels := []model.UserMemberLovePanel{}
	err := Engine.Table("s_user_member").
		Where("user_id = ?", session.UserStatus.UserID).Find(&lovePanels)
	if err != nil {
		panic(err)
	}
	return lovePanels
}
