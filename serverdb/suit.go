package serverdb

import (
	"elichika/model"

	"fmt"
)

func (session *Session) GetAllSuits() []model.UserSuit {
	suits := []model.UserSuit{}
	err := Engine.Table("s_user_suit").Where("user_id = ?", session.UserStatus.UserID).Find(&suits)
	if err != nil {
		panic(err)
	}
	return suits
}

func InsertUserSuits(suits []model.UserSuit) {
	if len(suits) == 0 {
		return
	}
	count, err := Engine.Table("s_user_suit").AllCols().Insert(suits)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", count, "suits")
}

func (session *Session) InsertUserSuit(suit model.UserSuit) {
	session.UserSuitDiffs = append(session.UserSuitDiffs, suit)
}

func (session *Session) FinalizeUserSuitDiffs() []any {
	InsertUserSuits(session.UserSuitDiffs)
	diffs := []any{}
	for _, suit := range session.UserSuitDiffs {
		diffs = append(diffs, suit.SuitMasterID)
		diffs = append(diffs, suit)
	}
	return diffs
}
