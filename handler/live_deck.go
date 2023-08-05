package handler

import (
	"bytes"
	"elichika/config"
	// "elichika/database"
	"elichika/model"
	"elichika/serverdb"

	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func SaveDeckAll(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())
	type SaveDeckAllReq struct {
		DeckID       int   `json:"deck_id"`
		CardWithSuit []int `json:"card_with_suit"`
		SquadDict    []any `json:"squad_dict"`
	}

	req := SaveDeckAllReq{}
	decoder := json.NewDecoder(strings.NewReader(reqBody.String()))
	decoder.UseNumber()
	err := decoder.Decode(&req)
	CheckErr(err)

	session := serverdb.GetSession(UserID)
	deckInfo := session.GetUserLiveDeck(req.DeckID)

	for i := 0; i < 9; i++ {
		if req.CardWithSuit[i*2+1] == 0 {
			req.CardWithSuit[i*2+1] = GetMemberDefaultSuitByCardMasterId(req.CardWithSuit[i*2])
		}
	}

	deckByte, _ := json.Marshal(deckInfo)
	deckJson := string(deckByte)
	for i := 0; i < 9; i++ {
		deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("card_master_id_%d", i+1), req.CardWithSuit[i*2])
		deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", i+1), req.CardWithSuit[i*2+1])
	}

	if err := json.Unmarshal([]byte(deckJson), &deckInfo); err != nil {
		panic(err)
	}
	// fmt.Println(deckInfo)

	session.UpdateUserLiveDeck(deckInfo)

	for k, v := range req.SquadDict {
		if k%2 == 0 {
			partyId, err := v.(json.Number).Int64()
			if err != nil {
				panic(err)
			}
			// fmt.Println("Party ID:", partyId)

			rDictInfo, err := json.Marshal(req.SquadDict[k+1])
			CheckErr(err)

			dictInfo := model.DeckSquadDict{}
			decoder := json.NewDecoder(bytes.NewReader(rDictInfo))
			decoder.UseNumber()
			err = decoder.Decode(&dictInfo)
			CheckErr(err)
			// fmt.Println("Party Info:", dictInfo)

			roleIds := []int{}
			err = MainEng.Table("m_card").
				Where("id IN (?,?,?)", dictInfo.CardMasterIds[0], dictInfo.CardMasterIds[1], dictInfo.CardMasterIds[2]).
				Cols("role").Find(&roleIds)
			CheckErr(err)
			// fmt.Println("roleIds:", roleIds)

			partyInfo := model.UserLiveParty{}
			partyInfo.UserID = UserID
			partyInfo.PartyID = int(partyId)
			partyIcon, partyName := GetPartyInfoByRoleIds(roleIds)
			partyInfo.Name.DotUnderText = GetRealPartyName(partyName)
			partyInfo.UserLiveDeckID = req.DeckID
			partyInfo.IconMasterID = partyIcon
			partyInfo.CardMasterID1 = dictInfo.CardMasterIds[0]
			partyInfo.CardMasterID2 = dictInfo.CardMasterIds[1]
			partyInfo.CardMasterID3 = dictInfo.CardMasterIds[2]
			partyInfo.UserAccessoryID1 = dictInfo.UserAccessoryIds[0]
			partyInfo.UserAccessoryID2 = dictInfo.UserAccessoryIds[1]
			partyInfo.UserAccessoryID3 = dictInfo.UserAccessoryIds[2]
			session.UpdateUserLiveParty(partyInfo)
		}
	}

	respBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), respBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchLiveDeckSelect(ctx *gin.Context) {
	// return last deck for this song
	signBody := GetData("fetchLiveDeckSelect.json")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SaveSuit(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type SaveSuitReq struct {
		DeckID       int `json:"deck_id"`
		CardIndex    int `json:"card_index"`
		SuitMasterID int `json:"suit_master_id"`
		ViewStatus   int `json:"view_status"` // 2 for Rina-chan board off, 1 for everyone else
	}
	// fmt.Println(reqBody)

	req := SaveSuitReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	CheckErr(err)

	session := serverdb.GetSession(UserID)
	deck := session.GetUserLiveDeck(req.DeckID)
	deckJsonByte, err := json.Marshal(deck)
	deckJson := string(deckJsonByte)
	deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", req.CardIndex), req.SuitMasterID)
	err = json.Unmarshal([]byte(deckJson), &deck)
	session.UpdateUserLiveDeck(deck)

	// Rina-chan board toggle
	if (req.SuitMasterID/10000)%1000 == 209 {
		RinaChan := session.GetMember(209)
		RinaChan.ViewStatus = req.ViewStatus
		session.UpdateMember(RinaChan)
	}

	signBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SaveDeck(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	// fmt.Println(reqBody)
	type SaveDeckReq struct {
		DeckID        int    `json:"deck_id"`
		CardMasterIDs [2]int `json:"card_master_ids"`
	}
	req := SaveDeckReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	CheckErr(err)

	position := req.CardMasterIDs[0]
	newCardMasterID := req.CardMasterIDs[1]
	newSuitMasterID := newCardMasterID

	session := serverdb.GetSession(UserID)

	// fetch the deck and parties affected
	deck := session.GetUserLiveDeck(req.DeckID)
	deckJsonByte, err := json.Marshal(deck)
	CheckErr(err)
	deckJson := string(deckJsonByte)

	oldCardMasterID := int(gjson.Get(deckJson, fmt.Sprintf("card_master_id_%d", position)).Int())
	oldSuitMasterID := int(gjson.Get(deckJson, fmt.Sprintf("suit_master_id_%d", position)).Int())
	if newCardMasterID == oldCardMasterID {
		panic("same card master id")
	}

	oldPosition := 0
	// old position != 0 then new card is in current deck, we have to swap
	gjson.Parse(deckJson).ForEach(func(key, value gjson.Result) bool {
		if strings.Contains(key.String(), "card_master_id") {
			if int(value.Int()) == newCardMasterID {
				oldPosition = int(key.String()[len(key.String())-1] - '0')
				// don't change suit_master_id if we just swap card around
				newSuitMasterID = int(gjson.Get(deckJson, fmt.Sprintf("suit_master_id_%d", oldPosition)).Int())
				return false
			}
		}
		return true
	})

	parties := []model.UserLiveParty{}
	parties = append(parties, session.GetUserLivePartyWithDeckAndCardID(req.DeckID, oldCardMasterID))
	if oldPosition != 0 {
		oldParty := session.GetUserLivePartyWithDeckAndCardID(req.DeckID, newCardMasterID)
		if oldParty.PartyID != parties[0].PartyID {
			parties = append(parties, oldParty)
		}
	}

	// change card master id in deck, then change the suit master id to default
	deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("card_master_id_%d", position), newCardMasterID)
	deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", position), newSuitMasterID)
	if oldPosition != 0 {
		deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("card_master_id_%d", oldPosition), oldCardMasterID)
		deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", oldPosition), oldSuitMasterID)
	}
	err = json.Unmarshal([]byte(deckJson), &deck)
	CheckErr(err)
	session.UpdateUserLiveDeck(deck)

	for _, party := range parties {
		partyJsonByte, err := json.Marshal(party)
		CheckErr(err)
		partyJson := string(partyJsonByte)
		// change the live party and update the names
		gjson.Parse(partyJson).ForEach(func(key, value gjson.Result) bool {
			if strings.Contains(key.String(), "card_master_id") {
				if int(value.Int()) == oldCardMasterID {
					partyJson, _ = sjson.Set(partyJson, key.String(), newCardMasterID)
				} else if int(value.Int()) == newCardMasterID {
					partyJson, _ = sjson.Set(partyJson, key.String(), oldCardMasterID)
				}
			}
			return true
		})

		partyInfo := gjson.Parse(partyJson)
		roleIds := []int{}
		err = MainEng.Table("m_card").
			Where("id IN (?,?,?)", partyInfo.Get("card_master_id_1").Int(),
				partyInfo.Get("card_master_id_2").Int(),
				partyInfo.Get("card_master_id_3").Int()).
			Cols("role").Find(&roleIds)
		CheckErr(err)
		partyIcon, partyName := GetPartyInfoByRoleIds(roleIds)
		realPartyName := GetRealPartyName(partyName)
		partyJson, _ = sjson.Set(partyJson, "name.dot_under_text", realPartyName)
		partyJson, _ = sjson.Set(partyJson, "icon_master_id", partyIcon)
		err = json.Unmarshal([]byte(partyJson), &party)
		CheckErr(err)
		session.UpdateUserLiveParty(party)
	}

	signBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
