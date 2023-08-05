package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"
	"net/http"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func FetchCommunicationMemberDetail(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	var memberId int
	gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
		if value.Get("member_id").String() != "" {
			memberId = int(value.Get("member_id").Int())
			return false
		}
		return true
	})

	session := serverdb.GetSession(UserID)
	lovePanelCellIds := session.GetLovePanelCellIDs(memberId)

	now := time.Now()
	year, month, day := now.Year(), now.Month(), now.Day()
	tomorrow := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location()).Unix()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	signBody := GetData("fetchCommunicationMemberDetail.json")
	signBody, _ = sjson.Set(signBody, "member_love_panels.0.member_id", memberId)
	signBody, _ = sjson.Set(signBody, "member_love_panels.0.member_love_panel_cell_ids", lovePanelCellIds)
	signBody, _ = sjson.Set(signBody, "weekday_state.weekday", weekday)
	signBody, _ = sjson.Set(signBody, "weekday_state.next_weekday_at", tomorrow)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

// where is this called?
func UpdateUserCommunicationMemberDetailBadge(ctx *gin.Context) {
	panic("UpdateUserCommunicationMemberDetailBadge")
	// reqBody := ctx.GetString("reqBody")
	// var memberMasterId int64
	// gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
	// 	if value.Get("member_master_id").String() != "" {
	// 		memberMasterId = value.Get("member_master_id").Int()
	// 		return false
	// 	}
	// 	return true
	// })

	// userDetail := []any{}
	// userDetail = append(userDetail, memberMasterId)
	// userDetail = append(userDetail, model.UserCommunicationMemberDetailBadgeByID{
	// 	MemberMasterID: int(memberMasterId),
	// })

	// signBody := GetData("updateUserCommunicationMemberDetailBadge.json")
	// signBody, _ = sjson.Set(signBody, "user_model.user_status", GetUserStatus())
	// signBody, _ = sjson.Set(signBody, "user_model.user_communication_member_detail_badge_by_id", userDetail)
	// resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	// ctx.Header("Content-Type", "application/json")
	// ctx.String(http.StatusOK, resp)
}

// seems like called in setting featured song, shouldn't be in this file
func UpdateUserLiveDifficultyNewFlag(ctx *gin.Context) {
	panic("UpdateUserLiveDifficultyNewFlag")
	// signBody, _ := sjson.Set(GetData("updateUserLiveDifficultyNewFlag.json"),
	// 	"user_model.user_status", GetUserStatus())
	// resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	// ctx.Header("Content-Type", "application/json")
	// ctx.String(http.StatusOK, resp)
}

func FinishUserStorySide(ctx *gin.Context) {
	// need to award items / mark as read
	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("finishUserStorySide.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FinishUserStoryMember(ctx *gin.Context) {
	// need to award items / mark as read
	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("finishUserStoryMember.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetTheme(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	session := serverdb.GetSession(UserID)

	var memberMasterID, suitMasterID, backgroundMasterID int
	gjson.Parse(reqBody).ForEach(func(key, value gjson.Result) bool {
		if value.Get("member_master_id").String() != "" {
			memberMasterID = int(value.Get("member_master_id").Int())
			suitMasterID = int(value.Get("suit_master_id").Int())
			backgroundMasterID = int(value.Get("custom_background_master_id").Int())

			member := session.GetMember(memberMasterID)
			member.SuitMasterID = suitMasterID
			member.CustomBackgroundMasterID = backgroundMasterID
			session.UpdateMember(member)
			return false
		}
		return true
	})

	userSuitRes := []any{}
	userSuitRes = append(userSuitRes, suitMasterID)
	userSuitRes = append(userSuitRes, model.SuitInfo{
		SuitMasterID: int(suitMasterID),
		IsNew:        false,
	})

	signBody := session.Finalize(GetData("setTheme.json"), "user_model")
	signBody, _ = sjson.Set(signBody, "user_model.user_suit_by_suit_id", userSuitRes)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetFavoriteMember(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	session := serverdb.GetSession(UserID)
	session.UserStatus.FavoriteMemberID = int(gjson.Parse(reqBody).Array()[0].Get("member_master_id").Int())
	signBody := session.Finalize(GetData("setFavoriteMember.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
