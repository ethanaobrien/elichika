package handler

import (
	// "elichika/config"
	"elichika/serverdb"

	// "fmt"
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func UpdateConsentState(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	session := serverdb.GetSession(UserID)
	session.UserStatus.GdprVersion = int(reqBody.Get("version").Int())
	session.Finalize(GetData("userModel.json"), "user_model")
	// Don't know the format of this response, but we can set gdpr version to 4 to skip it.
	// resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// ctx.Header("Content-Type", "application/json")
	// ctx.String(500, resp)
}
