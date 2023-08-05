package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func LiveMvStart(ctx *gin.Context) {
	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("liveMvStart.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LiveMvSaveDeck(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()

	type LiveSaveDeckReq struct {
		LiveMasterID        int   `json:"live_master_id"`
		LiveMvDeckType      int   `json:"live_mv_deck_type"`
		MemberMasterIDByPos []int `json:"member_master_id_by_pos"`
		SuitMasterIDByPos   []int `json:"suit_master_id_by_pos"`
		ViewStatusByPos     []int `json:"view_status_by_pos"`
	}

	req := LiveSaveDeckReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	if err != nil {
		panic(err)
	}

	userLiveMvDeckInfo := model.UserLiveMvDeckInfo{
		LiveMasterID: req.LiveMasterID,
	}
	deckJsonBytes, err := json.Marshal(userLiveMvDeckInfo)
	CheckErr(err)
	deckJson := string(deckJsonBytes)

	for k, v := range req.MemberMasterIDByPos {
		if k%2 == 0 {
			memberId := req.MemberMasterIDByPos[k+1]
			deckJson, err = sjson.Set(deckJson, fmt.Sprintf("member_master_id_%d", v), memberId)
		}
	}
	for k, v := range req.SuitMasterIDByPos {
		if k%2 == 0 {
			suitId := req.SuitMasterIDByPos[k+1]
			deckJson, err = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", v), suitId)
		}
	}
	err = json.Unmarshal([]byte(deckJson), &userLiveMvDeckInfo)
	CheckErr(err)
	session := serverdb.GetSession(UserID)
	for k, _ := range req.ViewStatusByPos {
		if k%2 == 0 {
			memberID := req.MemberMasterIDByPos[k+1]
			// Rina-chan board toggle
			if memberID == 209 {
				RinaChan := session.GetMember(209)
				RinaChan.ViewStatus = req.ViewStatusByPos[k+1]
				session.UpdateMember(RinaChan)
			}
		}
	}

	var userLiveMvDeckCustomByID []any
	userLiveMvDeckCustomByID = append(userLiveMvDeckCustomByID, req.LiveMasterID)
	userLiveMvDeckCustomByID = append(userLiveMvDeckCustomByID, userLiveMvDeckInfo)

	signBody := GetData("userModel.json")
	signBody = session.Finalize(signBody, "user_model")
	if req.LiveMvDeckType == 1 {
		signBody, _ = sjson.Set(signBody, "user_model.user_live_mv_deck_by_id", userLiveMvDeckCustomByID)
	} else {
		signBody, _ = sjson.Set(signBody, "user_model.user_live_mv_deck_custom_by_id", userLiveMvDeckCustomByID)
	}

	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
