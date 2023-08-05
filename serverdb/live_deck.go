package serverdb

import (
	"elichika/model"

	"fmt"
)

func (session *Session) GetUserLiveDeck(userLiveDeckID int) model.UserLiveDeck {
	liveDeck := model.UserLiveDeck{}
	exists, err := Engine.Table("s_user_live_deck").
		Where("user_id = ? AND user_live_deck_id = ?", session.UserStatus.UserID, userLiveDeckID).
		Get(&liveDeck)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("Deck doesn't exist")
	}
	return liveDeck
}

func (session *Session) UpdateUserLiveDeck(liveDeck model.UserLiveDeck) {
	session.UserLiveDeckDiffs[liveDeck.UserLiveDeckID] = liveDeck
}

func (session *Session) FinalizeUserLiveDeckDiffs() []any {
	userLiveDeckByID := []any{}
	for userLiveDeckId, userLiveDeck := range session.UserLiveDeckDiffs {
		userLiveDeckByID = append(userLiveDeckByID, userLiveDeckId)
		userLiveDeckByID = append(userLiveDeckByID, userLiveDeck)
		affected, err := Engine.Table("s_user_live_deck").
			Where("user_id = ? AND user_live_deck_id = ?", session.UserStatus.UserID, userLiveDeckId).
			AllCols().Update(userLiveDeck)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userLiveDeckByID
}

func (session *Session) GetAllLiveDecks() []model.UserLiveDeck {
	decks := []model.UserLiveDeck{}
	err := Engine.Table("s_user_live_deck").Where("user_id = ?", session.UserStatus.UserID).Find(&decks)
	if err != nil {
		panic(err)
	}
	return decks
}

func (session *Session) InsertLiveDecks(decks []model.UserLiveDeck) {
	count, err := Engine.Table("s_user_live_deck").Insert(&decks)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", count, " live decks")
}
