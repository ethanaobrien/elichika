package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func UpdateCardNewFlag(ctx *gin.Context) {
	// reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())
	session := serverdb.GetSession(UserID)

	signBody := session.Finalize(GetData("updateCardNewFlag.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ChangeIsAwakeningImage(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.CardAwakeningReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)
	userCard := session.GetUserCard(req.CardMasterID)
	userCard.IsAwakeningImage = req.IsAwakeningImage
	session.UpdateUserCard(userCard)

	cardResp := session.Finalize(GetData("changeIsAwakeningImage.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ChangeFavorite(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.CardFavoriteReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)
	userCard := session.GetUserCard(req.CardMasterID)
	userCard.IsFavorite = req.IsFavorite
	session.UpdateUserCard(userCard)

	cardResp := session.Finalize(GetData("changeFavorite.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func GetOtherUserCard(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())
	type OtherUserCardReq struct {
		UserID       int `json:"user_id"`
		CardMasterID int `json:"card_master_id"`
	}
	req := OtherUserCardReq{}
	// userCardReq := model.UserCardReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	partnerCard := serverdb.GetPartnerCardFromUserCard(serverdb.GetUserCard(req.UserID, req.CardMasterID))
	userCardResp, _ := sjson.Set("{}", "other_user_card", partnerCard)
	resp := SignResp(ctx.GetString("ep"), userCardResp, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
