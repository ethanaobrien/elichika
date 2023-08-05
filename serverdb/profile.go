package serverdb

import (
	"elichika/model"

	"encoding/json"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func FetchDBProfile(userID int, result interface{}) {
	exists, err := Engine.Table("s_user_info").Where("user_id = ?", userID).Get(result)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("user doesn't exist")
	}
}

func FetchPartnerCards(otherUserID int) []model.UserCard {
	partnerCards := []model.UserCard{}
	err := Engine.Table("s_user_card").
		Where("user_id = ? AND live_partner_categories != 0", otherUserID).
		Find(&partnerCards)
	if err != nil {
		panic(err)
	}
	return partnerCards
}

func GetPartnerCardFromUserCard(card model.UserCard) model.PartnerCardInfo {
	memberId := (card.CardMasterID / 10000) % 1000

	partnerCard := model.PartnerCardInfo{}

	jsonByte, err := json.Marshal(card)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonByte, &partnerCard)
	if err != nil {
		panic(err)
	}

	exists, err := Engine.Table("s_user_member").Where("user_id = ? AND member_master_id = ?", card.UserID, memberId).
		Cols("love_level").Get(&partnerCard.LoveLevel)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("member doesn't exist")
	}

	partnerCard.PassiveSkillLevels = []int{}
	partnerCard.PassiveSkillLevels = append(partnerCard.PassiveSkillLevels, card.PassiveSkillALevel)
	partnerCard.PassiveSkillLevels = append(partnerCard.PassiveSkillLevels, card.PassiveSkillBLevel)
	partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill1ID)
	partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill2ID)
	partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill3ID)
	partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill4ID)
	partnerCard.MemberLovePanels = []int{}

	// filling this for a card of self freeze the game
	// the displayed value still correct for own's card in the guest setup menu with an empty array
	// display value in getOtherUserCard is wrong, but if we fill in for own card then it also freeze
	// TODO: revisit after implmenting friends

	// lovePanel := model.UserMemberLovePanel{}
	// exists, err = Engine.Table("s_user_member").
	// 	Where("user_id = ? AND member_master_id = ?", card.UserID, memberId).Get(&lovePanel)
	// if err != nil {
	// 	panic(err)
	// }
	// if !exists {
	// 	panic("member doesn't exist")
	// }
	// partnerCard.MemberLovePanels = lovePanel.MemberLovePanelCellIDs

	return partnerCard
}

func GetUserCard(userID, cardMasterID int) model.UserCard {
	card := model.UserCard{}
	exists, err := Engine.Table("s_user_card").Where("user_id = ? AND card_master_id = ?", userID, cardMasterID).
		Get(&card)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("user card doesn't exist")
	}
	return card
}

func (sesison *Session) GetOtherUserBasicProfile(otherUserID int) model.UserBasicInfo {
	basicInfo := model.UserBasicInfo{}
	FetchDBProfile(otherUserID, &basicInfo)
	recommendCard := GetUserCard(otherUserID, basicInfo.RecommendCardMasterID)

	basicInfo.RecommendCardLevel = recommendCard.Level
	basicInfo.IsRecommendCardImageAwaken = recommendCard.IsAwakeningImage
	basicInfo.IsRecommendCardAllTrainingActivated = recommendCard.IsAllTrainingActivated

	// friend system, not implemented
	basicInfo.FriendApprovedAt = new(int64)
	*basicInfo.FriendApprovedAt = 0
	basicInfo.RequestStatus = 3
	basicInfo.IsRequestPending = false
	return basicInfo
}

// fetch profile of another user, from session.UserStatus.UserID's perspective
// it's possible that otherUserID == session.UserStatus.UserID
func (session *Session) FetchProfile(otherUserID int) model.Profile {
	profile := model.Profile{}

	exists, err := Engine.Table("s_user_info").Where("user_id = ?", otherUserID).Get(&profile)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("user doesn't exist")
	}
	// recommend card

	recommendCard := GetUserCard(otherUserID, profile.ProfileInfo.BasicInfo.RecommendCardMasterID)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("card doesn't exist")
	}
	profile.ProfileInfo.BasicInfo.RecommendCardLevel = recommendCard.Level
	profile.ProfileInfo.BasicInfo.IsRecommendCardImageAwaken = recommendCard.IsAwakeningImage
	profile.ProfileInfo.BasicInfo.IsRecommendCardAllTrainingActivated = recommendCard.IsAllTrainingActivated

	// friend system, not implemented
	profile.ProfileInfo.BasicInfo.FriendApprovedAt = new(int64)
	*profile.ProfileInfo.BasicInfo.FriendApprovedAt = 0
	profile.ProfileInfo.BasicInfo.RequestStatus = 3
	profile.ProfileInfo.BasicInfo.IsRequestPending = false

	// other user's members
	members := []model.UserMemberInfo{}
	err = Engine.Table("s_user_member").Where("user_id = ?", otherUserID).OrderBy("love_point DESC").Find(&members)
	if err != nil {
		panic(err)
	}
	profile.ProfileInfo.TotalLovePoint = 0
	for _, member := range members {
		profile.ProfileInfo.TotalLovePoint += member.LovePoint
	}
	for i := 0; i < 3; i++ {
		profile.ProfileInfo.LoveMembers[i].MemberMasterID = members[i].MemberMasterID
		profile.ProfileInfo.LoveMembers[i].LovePoint = members[i].LovePoint
	}

	// need to return this in order, so make the array then write to it
	// this also prevent the game from freezing if 2 cards or more have the same bit
	for i := 0; i < 7; i++ {
		profile.GuestInfo.LivePartnersCards = append(profile.GuestInfo.LivePartnersCards, model.LivePartnerCard{})
		profile.GuestInfo.LivePartnersCards[i].LivePartnerCategoryMasterId = i + 1
	}

	partnerCards := FetchPartnerCards(otherUserID)
	for _, card := range partnerCards {
		partnerCard := GetPartnerCardFromUserCard(card)
		livePartner := model.LivePartnerCard{}
		livePartner.PartnerCard = partnerCard
		for i := 1; i <= 7; i++ {
			if (card.LivePartnerCategories & (1 << i)) != 0 {
				profile.GuestInfo.LivePartnersCards[i-1].PartnerCard = partnerCard
			}
		}
	}
	liveStats := model.DBUserProfileLiveStats{}
	_, err = Engine.Table("s_user_info").Where("user_id = ?", otherUserID).Get(&liveStats)
	if err != nil {
		panic(err)
	}
	jsonByte, err := json.Marshal(liveStats)
	if err != nil {
		panic(err)
	}
	gjson.Parse(string(jsonByte)).ForEach(func(key, value gjson.Result) bool {
		songRarity, _ := strconv.Atoi(key.String()[len(key.String())-2:])
		if strings.Contains(key.String(), "LivePlayCount") {
			profile.PlayInfo.LivePlayCount = append(profile.PlayInfo.LivePlayCount, songRarity)
			profile.PlayInfo.LivePlayCount = append(profile.PlayInfo.LivePlayCount, int(value.Int()))
		} else if strings.Contains(key.String(), "LiveClearCount") {
			profile.PlayInfo.LiveClearCount = append(profile.PlayInfo.LiveClearCount, songRarity)
			profile.PlayInfo.LiveClearCount = append(profile.PlayInfo.LiveClearCount, int(value.Int()))
		}
		return true
	})

	Engine.Table("s_user_card").Where("user_id = ?", otherUserID).
		OrderBy("live_join_count DESC").Limit(3).Find(&profile.PlayInfo.JoinedLiveCardRanking)
	Engine.Table("s_user_card").Where("user_id = ?", otherUserID).
		OrderBy("active_skill_play_count DESC").Limit(3).Find(&profile.PlayInfo.PlaySkillCardRanking)

	// can get from members[] to save sql
	Engine.Table("s_user_member").Where("user_id = ?", otherUserID).Find(&profile.MemberInfo.UserMembers)
	return profile
}
