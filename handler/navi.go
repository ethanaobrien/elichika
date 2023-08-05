package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"encoding/json"
	"net/http"
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func SaveUserNaviVoice(ctx *gin.Context) {
	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("saveUserNaviVoice.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func TapLovePoint(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type TapLovePointReq struct {
		MemberMasterID int `json:"member_master_id"`
	}

	req := TapLovePointReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}
	session := serverdb.GetSession(UserID)
	member := session.GetMember(req.MemberMasterID)
	member.LovePoint += 20
	if member.LovePoint > member.LovePointLimit {
		member.LovePoint = member.LovePointLimit
	}
	member.LoveLevel = GetBondLevel(member.LovePoint)
	session.UpdateMember(member)

	signBody := session.Finalize(GetData("saveUserNaviVoice.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
