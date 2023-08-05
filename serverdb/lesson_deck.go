package serverdb

import (
	"elichika/model"

	"fmt"
)

func (session *Session) GetLessonDeck(userLessonDeckId int) model.UserLessonDeck {
	deck, exist := session.UserLessonDeckDiffs[userLessonDeckId]
	if exist {
		return deck
	}
	deck = model.UserLessonDeck{}
	exists, err := Engine.Table("s_user_lesson_deck").
		Where("user_id = ? AND user_lesson_deck_id = ?", session.UserStatus.UserID, userLessonDeckId).
		Get(&deck)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("deck not found")
	}
	return deck
}

func (session *Session) UpdateLessonDeck(userLessonDeck model.UserLessonDeck) {
	session.UserLessonDeckDiffs[userLessonDeck.UserLessonDeckID] = userLessonDeck
}

func (session *Session) FinalizeUserLessonDeckDiffs() []any {
	userLessonDeckByID := []any{}
	for userLessonDeckId, userLessonDeck := range session.UserLessonDeckDiffs {
		userLessonDeckByID = append(userLessonDeckByID, userLessonDeckId)
		userLessonDeckByID = append(userLessonDeckByID, userLessonDeck)
		affected, err := Engine.Table("s_user_lesson_deck").
			Where("user_id = ? AND user_lesson_deck_id = ?", session.UserStatus.UserID, userLessonDeckId).
			AllCols().Update(userLessonDeck)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userLessonDeckByID
}

func (session *Session) GetAllLessonDecks() []model.UserLessonDeck {
	decks := []model.UserLessonDeck{}
	err := Engine.Table("s_user_lesson_deck").Where("user_id = ?", session.UserStatus.UserID).Find(&decks)
	if err != nil {
		panic(err)
	}
	return decks
}

func (session *Session) InsertLessonDecks(decks []model.UserLessonDeck) {
	count, err := Engine.Table("s_user_lesson_deck").Insert(&decks)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", count, " lesson decks")
}
