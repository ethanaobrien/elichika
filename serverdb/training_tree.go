package serverdb

import (
	"elichika/model"

	"fmt"
)

// return the training tree for a card
func (session *Session) GetTrainingTree(cardMasterID int) []model.TrainingTreeCell {
	cells := []model.TrainingTreeCell{}
	err := Engine.Table("s_user_training_tree_cell").
		Where("user_id = ? AND card_master_id = ?", session.UserStatus.UserID, cardMasterID).Find(&cells)
	if err != nil {
		panic(err)
	}
	return cells
}

// insert a training cell sets
func (session *Session) InsertTrainingCells(cells *[]model.TrainingTreeCell) {
	n := len(*cells)
	step := 5000
	for begin := 0; begin < n; begin += step {
		end := begin + step
		if end > n {
			end = n
		}
		affected, err := Engine.Table("s_user_training_tree_cell").AllCols().Insert((*cells)[begin:end])
		if err != nil {
			panic(err)
		}
		fmt.Println("Inserted ", affected, " training cells")
	}
}
