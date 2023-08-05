package handler

import (
	"elichika/config"
	"elichika/encrypt"
	// "elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func StartUp(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")

	var mask64 string
	req := gjson.Parse(reqBody)
	req.ForEach(func(key, value gjson.Result) bool {
		if value.Get("mask").String() != "" {
			mask64 = value.Get("mask").String()
			return false
		}
		return true
	})
	// fmt.Println("Request data:", req.String())
	// fmt.Println("Mask:", mask64)

	mask, err := base64.StdEncoding.DecodeString(mask64)
	if err != nil {
		panic(err)
	}
	randomBytes := encrypt.RSA_DecryptOAEP(mask, "privatekey.pem")
	// fmt.Println("Random Bytes:", randomBytes)

	newKey := utils.Xor(randomBytes, []byte(config.SessionKey))
	newKey64 := base64.StdEncoding.EncodeToString(newKey)
	// fmt.Println("Session Key:", newKey64)

	startupBody := GetData("startup.json")
	startupBody, _ = sjson.Set(startupBody, "authorization_key", newKey64)
	resp := SignResp(ctx.GetString("ep"), startupBody, StartUpKey)
	// fmt.Println("Response:", resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func Login(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")

	var mask64 string
	req := gjson.Parse(reqBody)
	req.ForEach(func(key, value gjson.Result) bool {
		if value.Get("mask").String() != "" {
			mask64 = value.Get("mask").String()
			return false
		}
		return true
	})
	// fmt.Println("Request data:", req.String())
	// fmt.Println("Mask:", mask64)

	mask, err := base64.StdEncoding.DecodeString(mask64)
	if err != nil {
		panic(err)
	}
	randomBytes := encrypt.RSA_DecryptOAEP(mask, "privatekey.pem")
	// fmt.Println("Random Bytes:", randomBytes)

	serverEventReceiverKey, err := hex.DecodeString(config.ServerEventReceiverKey)
	if err != nil {
		panic(err)
	}

	jaKey, err := hex.DecodeString(config.JaKey)
	if err != nil {
		panic(err)
	}

	newKey := utils.Xor(randomBytes, []byte(config.SessionKey))
	newKey = utils.Xor(newKey, serverEventReceiverKey)
	newKey = utils.Xor(newKey, jaKey)
	newKey64 := base64.StdEncoding.EncodeToString(newKey)
	// fmt.Println("Session Key:", newKey64)
	serverdb.InitDb(IsGlobal)
	session := serverdb.GetSession(UserID)
	session.UserStatus.LastLoginAt = ClientTimeStamp

	loginBody := GetData("login.json")
	loginBody, _ = sjson.Set(loginBody, "session_key", newKey64)

	loginBody, _ = sjson.Set(loginBody, "user_model.user_status", session.UserStatus)

	/* ======== UserData ======== */
	fmt.Println("User logins: ", UserID)

	// live decks
	dbLiveDecks := session.GetAllLiveDecks()
	if len(dbLiveDecks) == 0 {
		panic("no live deck found")
	}
	userLiveDecks := []any{}
	for _, liveDeckInfo := range dbLiveDecks {
		userLiveDecks = append(userLiveDecks, liveDeckInfo.UserLiveDeckID)
		userLiveDecks = append(userLiveDecks, liveDeckInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_live_deck_by_id", userLiveDecks)

	dbLiveParties := session.GetAllLiveParties()
	if len(dbLiveParties) == 0 {
		panic("no live party")
	}
	userLiveParties := []any{}
	for _, livePartyInfo := range dbLiveParties {
		userLiveParties = append(userLiveParties, livePartyInfo.PartyID)
		userLiveParties = append(userLiveParties, livePartyInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_live_party_by_id", userLiveParties)

	// member settings
	dbMembers := session.GetAllMembers()
	if len(dbMembers) == 0 {
		panic("no member found")
	}
	var userMembers []any
	for _, memberInfo := range dbMembers {
		userMembers = append(userMembers, memberInfo.MemberMasterID)
		userMembers = append(userMembers, memberInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_member_by_member_id", userMembers)

	// member love panel settings
	dbLovePanels := session.GetAllMemberLovePanels()
	if len(dbLovePanels) == 0 {
		panic("no member love panel found")
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.member_love_panels", dbLovePanels)

	// lesson decks
	dbLessonDecks := session.GetAllLessonDecks()
	if len(dbLessonDecks) == 0 {
		panic("no lesson deck")
	}

	userLessonDecks := []any{}
	for _, userLessonDeck := range dbLessonDecks {
		userLessonDecks = append(userLessonDecks, userLessonDeck.UserLessonDeckID)
		userLessonDecks = append(userLessonDecks, userLessonDeck)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_lesson_deck_by_id", userLessonDecks)

	// user cards
	dbCards := session.GetAllCards()
	if len(dbCards) == 0 {
		panic("no card")
	}

	userCards := []any{}
	for _, userCard := range dbCards {
		userCards = append(userCards, userCard.CardMasterID)
		userCards = append(userCards, userCard)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_card_by_card_id", userCards)

	// user suits
	dbSuits := session.GetAllSuits()
	if len(dbSuits) == 0 {
		panic("no suit")
	}

	userSuits := []any{}
	for _, userSuit := range dbSuits {
		userSuits = append(userSuits, userSuit.SuitMasterID)
		userSuits = append(userSuits, userSuit)
	}
	loginBody, err = sjson.Set(loginBody, "user_model.user_suit_by_suit_id", userSuits)
	CheckErr(err)

	// user accessory
	var UserAccessory []any
	decoder := json.NewDecoder(strings.NewReader(
		gjson.Parse(GetUserAccessoryData()).Get("user_accessory_by_user_accessory_id").String()))
	decoder.UseNumber()
	err = decoder.Decode(&UserAccessory)
	CheckErr(err)
	loginBody, _ = sjson.Set(loginBody, "user_model.user_accessory_by_user_accessory_id", UserAccessory)
	/* ======== UserData ======== */
	session.Finalize("{}", "")
	resp := SignResp(ctx.GetString("ep"), loginBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
