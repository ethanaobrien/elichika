package handler

import (
	"elichika/config"
	// "elichika/enum"
	"elichika/model"
	"elichika/protocol/request"
	"elichika/protocol/response"
	"elichika/userdata"
	"elichika/utils"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func FetchTowerSelect(ctx *gin.Context) {
	// there's no request body
	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	// no need to return anything, the same use database for this
	respObj := response.FetchTowerSelectResponse{
		TowerIDs:      []int{},
		UserModelDiff: &session.UserModel,
	}

	respBytes, _ := json.Marshal(respObj)
	resp := SignResp(ctx, string(respBytes), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchTowerTop(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	fmt.Println(reqBody)
	req := request.FetchTowerTopRequest{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)

	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	respObj := response.FetchTowerTopResponse{
		TowerCardUsedCountRows: []model.UserTowerCardUsedCount{},
		IsShowUnlockEffect:     false, // TODO: check if user has this tower to set this correctly
		UserModelDiff:          &session.UserModel,
		// other fields are for DLP with voltage ranking
	}

	respBytes, _ := json.Marshal(respObj)
	resp := SignResp(ctx, string(respBytes), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
