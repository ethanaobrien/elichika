package userdata

import (
	"elichika/model"
	"elichika/protocol/response"
	"elichika/utils"

	"time"
)

func (session *Session) GetUserTowerCardUsed(towerID, cardMasterID int) model.UserTowerCardUsedCount {
	cardUsed := model.UserTowerCardUsedCount{}
	exist, err := session.Db.Table("u_tower_card_used").
		Where("user_id = ? AND tower_id = ? AND card_master_id = ?", session.UserStatus.UserID, towerID, cardMasterID).Get(&cardUsed)
	utils.CheckErr(err)
	if !exist {
		cardUsed = model.UserTowerCardUsedCount{
			UserID:         session.UserStatus.UserID,
			TowerID:        towerID,
			CardMasterID:   cardMasterID,
			UsedCount:      0,
			RecoveredCount: 0,
			LastUsedAt:     0,
		}
	}
	return cardUsed
}

func (session *Session) UpdateUserTowerCardUsed(card model.UserTowerCardUsedCount) {
	affected, err := session.Db.Table("u_tower_card_used").
		Where("user_id = ? AND tower_id = ? AND card_master_id = ?", session.UserStatus.UserID, card.TowerID, card.CardMasterID).
		AllCols().Update(card)
	utils.CheckErr(err)
	if affected == 0 {
		_, err := session.Db.Table("u_tower_card_used").Insert(card)
		utils.CheckErr(err)
	}
}

func (session *Session) GetUserTowerCardUsedList(towerID int) []model.UserTowerCardUsedCount {
	list := []model.UserTowerCardUsedCount{}
	err := session.Db.Table("u_tower_card_used").
		Where("user_id = ? AND tower_id = ?", session.UserStatus.UserID, towerID).Find(&list)
	utils.CheckErr(err)
	return list
}

func (session *Session) GetUserTower(towerID int) model.UserTower {
	pos, exist := session.UserTowerMapping.SetList(&session.UserModel.UserTowerByTowerID).Map[int64(towerID)]
	if exist {
		return session.UserModel.UserTowerByTowerID.Objects[pos]
	}
	tower := model.UserTower{}
	exist, err := session.Db.Table("u_tower").
		Where("user_id = ? AND tower_id = ?", session.UserStatus.UserID, towerID).Get(&tower)
	utils.CheckErr(err)
	if !exist {
		tower = model.UserTower{
			UserID:                      session.UserStatus.UserID,
			TowerID:                     towerID,
			ClearedFloor:                0,
			ReadFloor:                   0,
			Voltage:                     0,
			RecoveryPointFullAt:         int(time.Now().Unix() + 86400),
			RecoveryPointLastConsumedAt: int(time.Now().Unix()),
		}
	} else {

		// this is to make sure user can always mass recover LP
		// it's the only way to do this since these things are decided by the timestamps alone (+ plus the limit from the database)
		tower.RecoveryPointFullAt = int(time.Now().Unix() + 86400)
		tower.RecoveryPointLastConsumedAt = int(time.Now().Unix())
	}
	return tower
}

func (session *Session) UpdateUserTower(tower model.UserTower) {
	session.UserTowerMapping.SetList(&session.UserModel.UserTowerByTowerID).Update(tower)
}

func (session *Session) GetUserTowerVoltageRankingScores(towerID int) []model.UserTowerVoltageRankingScore {
	scores := []model.UserTowerVoltageRankingScore{}
	err := session.Db.Table("u_tower_voltage_ranking_score").
		Where("user_id = ? AND tower_id = ?", session.UserStatus.UserID, towerID).Find(&scores)
	utils.CheckErr(err)
	return scores
}

func (session *Session) GetUserTowerVoltageRankingScore(towerID, floorNo int) model.UserTowerVoltageRankingScore {
	score := model.UserTowerVoltageRankingScore{}
	exists, err := session.Db.Table("u_tower_voltage_ranking_score").
		Where("user_id = ? AND tower_id = ? AND floor_no = ?", session.UserStatus.UserID, towerID, floorNo).Get(&score)
	utils.CheckErr(err)
	if !exists {
		score = model.UserTowerVoltageRankingScore{
			UserID:  session.UserStatus.UserID,
			TowerID: towerID,
			FloorNo: floorNo,
			Voltage: 0,
		}
	}
	return score
}

func (session *Session) UpdateUserTowerVoltageRankingScore(score model.UserTowerVoltageRankingScore) {
	affected, err := session.Db.Table("u_tower_voltage_ranking_score").
		Where("user_id = ? AND tower_id = ? AND floor_no = ?", session.UserStatus.UserID, score.TowerID, score.FloorNo).AllCols().
		Update(score)
	utils.CheckErr(err)
	if affected == 0 {
		_, err := session.Db.Table("u_tower_voltage_ranking_score").Insert(score)
		utils.CheckErr(err)
	}
}

func (session *Session) GetTowerRankingCell(towerID int) response.TowerRankingCell {
	scores := session.GetUserTowerVoltageRankingScores(towerID)
	cell := response.TowerRankingCell{
		Order:            1,
		SumVoltage:       0,
		TowerRankingUser: session.GetRankingUser(),
	}
	for _, score := range scores {
		cell.SumVoltage += score.Voltage
	}
	return cell
}

func towerFinalizer(session *Session) {
	for _, userTower := range session.UserModel.UserTowerByTowerID.Objects {
		affected, err := session.Db.Table("u_tower").
			Where("user_id = ? AND tower_id = ?", session.UserStatus.UserID, userTower.TowerID).
			AllCols().Update(userTower)
		utils.CheckErr(err)
		if affected == 0 {
			_, err = session.Db.Table("u_tower").Insert(userTower)
			utils.CheckErr(err)
		}
	}
}

func init() {
	addFinalizer(towerFinalizer)
	addGenericTableFieldPopulator("u_tower", "UserTowerByTowerID")
}
