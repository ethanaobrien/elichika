package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/tidwall/sjson"
)

func FetchBootstrap(ctx *gin.Context) {
	session := serverdb.GetSession(UserID)
	session.UserStatus.BootstrapSifidCheckAt = ClientTimeStamp
	signBody := session.Finalize(GetData("fetchBootstrap.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func GetClearedPlatformAchievement(ctx *gin.Context) {
	signBody := GetData("getClearedPlatformAchievement.json")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
