package serverdb

import (
	"elichika/model"

	"fmt"
)

func (session *Session) GetUserLiveParty(partyID int) model.UserLiveParty {
	liveParty := model.UserLiveParty{}
	exists, err := Engine.Table("s_user_live_party").
		Where("user_id = ? AND party_id = ?", session.UserStatus.UserID, partyID).
		Get(&liveParty)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("Party doesn't exist")
	}
	return liveParty
}

func (session *Session) GetUserLivePartyWithDeckAndCardID(deckID int, cardID int) model.UserLiveParty {
	liveParty := model.UserLiveParty{}
	exists, err := Engine.Table("s_user_live_party").
		Where("user_id = ? AND user_live_deck_id = ? AND (card_master_id_1 = ? OR card_master_id_2 = ? OR card_master_id_3 = ?)",
			session.UserStatus.UserID, deckID, cardID, cardID, cardID).
		Get(&liveParty)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("Party doesn't exist")
	}
	return liveParty
}

func (session *Session) UpdateUserLiveParty(liveParty model.UserLiveParty) {
	session.UserLivePartyDiffs[liveParty.PartyID] = liveParty
}

func (session *Session) FinalizeUserLivePartyDiffs() []any {
	userLivePartyByID := []any{}
	for userLivePartyId, userLiveParty := range session.UserLivePartyDiffs {
		userLivePartyByID = append(userLivePartyByID, userLivePartyId)
		userLivePartyByID = append(userLivePartyByID, userLiveParty)
		affected, err := Engine.Table("s_user_live_party").
			Where("user_id = ? AND party_id = ?", session.UserStatus.UserID, userLivePartyId).
			AllCols().Update(userLiveParty)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userLivePartyByID
}

func (session *Session) GetAllLiveParties() []model.UserLiveParty {
	parties := []model.UserLiveParty{}
	err := Engine.Table("s_user_live_party").Where("user_id = ?", session.UserStatus.UserID).Find(&parties)
	if err != nil {
		panic(err)
	}
	return parties
}

func (session *Session) InsertLiveParties(parties []model.UserLiveParty) {
	count, err := Engine.Table("s_user_live_party").Insert(&parties)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", count, " live parties")
}
