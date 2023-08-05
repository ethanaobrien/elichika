package serverdb

import (
	"elichika/model"
)

func SaveLiveState(live model.LiveState) {
	// delete whatever is there
	affected, err := Engine.Table("s_user_live_state").Where("user_id = ?", live.UserID).Delete(&model.LiveState{})
	if err != nil {
		panic(err)
	}
	affected, err = Engine.Table("s_user_live_state").AllCols().Insert(live)
	if err != nil {
		panic(err)
	}
	if affected != 1 {
		panic("failed to insert")
	}
}

func LoadLiveState(userID int) (bool, model.LiveState) {
	live := model.LiveState{}
	exists, err := Engine.Table("s_user_live_state").Where("user_id = ?", userID).Get(&live)
	// fmt.Println(live)
	if err != nil {
		panic(err)
	}
	if exists {
		_, err = Engine.Table("s_user_live_state").Where("user_id = ?", userID).Delete(&model.LiveState{})
		if err != nil {
			panic(err)
		}
	}
	return exists, live
}
