package serverdb

import (
	"elichika/model"
	// "elichika/utils"

	// "encoding/json"
	"fmt"
	// "github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// A session is a complete transation between server and client
// so 1 session per request
// A session fetch the data needs to be modified.
type Session struct {
	UserStatus          model.UserStatus
	CardDiffs           map[int]model.UserCard
	UserMemberDiffs     map[int]model.UserMemberInfo
	UserLessonDeckDiffs map[int]model.UserLessonDeck
	UserLiveDeckDiffs   map[int]model.UserLiveDeck
	UserLivePartyDiffs  map[int]model.UserLiveParty
	UserSuitDiffs       []model.UserSuit
	CardGradeUpTriggers []any
}

// Push update into the db and create the diff
// The actual response depend on the API, but they often contain the diff somewhere
// The mainKey is the key to the diff
func (session *Session) Finalize(jsonBody string, mainKey string) string {
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_status", session.FinalizeUserInfo())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_member_by_member_id", session.FinalizeUserMemberDiffs())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_card_by_card_id", session.FinalizeCardDiffs())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_info_trigger_card_grade_up_by_trigger_id", session.FinalizeCardGradeUpTrigger())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_lesson_deck_by_id", session.FinalizeUserLessonDeckDiffs())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_live_deck_by_id", session.FinalizeUserLiveDeckDiffs())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_live_party_by_id", session.FinalizeUserLivePartyDiffs())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_suit_by_suit_id", session.FinalizeUserSuitDiffs())
	return jsonBody
}

// fetch the user, this is always sent back to client
func (session *Session) InitUser(userID int) {
	session.UserStatus.UserID = userID
	exists, err := Engine.Table("s_user_info").Where("user_id = ?", userID).Get(&session.UserStatus)
	if err != nil {
		panic(err)
	}
	if !exists { // create one
		panic(fmt.Sprintf("user doesn't exist %d", userID))
	}
}

func (session *Session) FinalizeUserInfo() model.UserStatus {
	_, err := Engine.Table("s_user_info").Where("user_id = ?", session.UserStatus.UserID).Update(&session.UserStatus)
	if err != nil {
		panic(err)
	}
	return session.UserStatus
}

func GetSession(userId int) Session {
	s := Session{}
	s.CardDiffs = make(map[int]model.UserCard)
	s.UserMemberDiffs = make(map[int]model.UserMemberInfo)
	s.UserLessonDeckDiffs = make(map[int]model.UserLessonDeck)
	s.UserLiveDeckDiffs = make(map[int]model.UserLiveDeck)
	s.UserLivePartyDiffs = make(map[int]model.UserLiveParty)
	s.CardGradeUpTriggers = make([]any, 0)
	s.UserSuitDiffs = make([]model.UserSuit, 0)
	s.InitUser(userId)
	return s
}
