package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"
	"encoding/json"

	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	// "xorm.io/xorm"
)

func FetchTrainingTree(ctx *gin.Context) {
	// signBody := GetData("fetchTrainingTree.json")
	// resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	// ctx.Header("Content-Type", "application/json")
	// ctx.String(http.StatusOK, resp)

	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	type FetchTrainingTreeReq struct {
		CardMasterID int `json:"card_master_id"`
	}
	req := FetchTrainingTreeReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}
	session := serverdb.GetSession(UserID)
	signBody := `"{}"`
	signBody, _ = sjson.Set(signBody, "user_card_training_tree_cell_list", session.GetTrainingTree(req.CardMasterID))
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LevelUpCard(ctx *gin.Context) {
	session := serverdb.GetSession(UserID)

	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]

	type LevelUpCardReq struct {
		CardMasterId    int `json:"card_master_id"`
		AdditionalLevel int `json:"additional_level"`
	}

	req := LevelUpCardReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	userCard := session.GetUserCard(req.CardMasterId)
	userCard.Level += req.AdditionalLevel
	session.UpdateUserCard(userCard)
	signBody := session.Finalize(GetData("userModelDiff.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)

	// TODO: Handle things like exp and gold cost

	// SendCardInfoDiff(ctx, &userCard)

	// SendCardInfoDiff(ctx, &userCard)
}

func BondRequired(l int) int {
	res := 30 * l
	if l > 2 {
		res += 10 * (l - 2)
	}
	if l > 6 {
		res += 10 * (l - 6)
	}
	if l > 20 {
		res += 10 * (l - 20)
	}
	if l > 59 {
		res += 10 * (l - 59)
	}
	return res
}

func BondRequiredTotal(l int) int {
	res := 0
	for i := 2; i <= l; i++ {
		res += BondRequired(i)
	}
	return res
}

func GetBondLevel(maxBond int) int {
	res := 0
	for i := 2; ; i++ {
		res += BondRequired(i)
		if res > maxBond {
			return i - 1
		}
	}
}

func GradeUpCard(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	type GradeUpCardReq struct {
		CardMasterID int `json:"card_master_id"`
	}
	req := GradeUpCardReq{}

	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)

	userCard := session.GetUserCard(req.CardMasterID)
	memberInfo := session.GetMember(GetMemberMasterIdByCardMasterId(req.CardMasterID))
	userCard.Grade += 1
	currentBondLevel := GetBondLevel(memberInfo.LovePointLimit)
	currentBondLevel += 3
	memberInfo.LovePointLimit = BondRequiredTotal(currentBondLevel)
	session.UpdateUserCard(userCard)
	session.UpdateMember(memberInfo)

	// we need to set user_info_trigger_card_grade_up_by_trigger_id
	// for the pop up after limit breaking

	type Trigger struct {
		TriggerId            int64 `json:"trigger_id"`
		CardMasterId         int   `json:"card_master_id"`
		BeforeLoveLevelLimit int   `json:"before_love_level_limit"`
		AfterLoveLevelLimit  int   `json:"after_love_level_limit"`
	}

	trigger := Trigger{}
	// this trigger show the pop up after limit break

	// TODO: we load the card again, the animation will be played again
	// this has something to do with the state of the game, as restarting fix this
	// the first 10 digit is certainly the time stamp in unix second
	// after that there's 9 digit, but it's unclear what they actually mean.
	// could be that it's just a time stamp is unix nanosecond, and something else control how the pop-up behave
	trigger.TriggerId = ClientTimeStamp * 1000000
	trigger.CardMasterId = userCard.CardMasterID
	trigger.BeforeLoveLevelLimit = currentBondLevel - 3
	trigger.AfterLoveLevelLimit = currentBondLevel

	session.AddCardGradeUpTrigger(trigger.TriggerId, trigger)

	resp := session.Finalize(GetData("userModelDiff.json"), "user_model_diff")
	resp = SignResp(ctx.GetString("ep"), resp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}

func ActivateTrainingTreeCell(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	type ActivateTrainingTreeCellReq struct {
		CardMasterID  int   `json:"card_master_id"`
		CellMasterIDs []int `json:"cell_master_ids"`
		PayType       int   `json:"pay_type"`
	}
	req := ActivateTrainingTreeCellReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)

	db := GetMasterdataDb()

	type TrainingTreeMapping struct {
		// ID int  // card master id
		TrainingTreeMappingMID int `xorm:"'training_tree_mapping_m_id'"` // the training tree
		// TrainingTreeCardParamMID int  // same as card master id
		TrainingTreeCardPassiveSkillIncreaseMID int `xorm:"'training_tree_card_passive_skill_increase_m_id'"` // 1 to 2
	}

	treeMapping := TrainingTreeMapping{}
	exists, err := db.Table("m_training_tree").Where("id = ?", req.CardMasterID).
		Cols("training_tree_mapping_m_id", "training_tree_card_passive_skill_increase_m_id").Get(&treeMapping)
	if (err != nil) || (!exists) {
		panic(err)
	}

	type TrainingTreeCellContent struct {
		// Id int // tree id
		CellID               int `xorm:"'cell_id'"` // cell id
		TrainingTreeCellType int // type of the cell
		TrainingContentNo    int // the content of the cell
		// RequiredGrade int // the limit break required, no need to read this here
		TrainingTreeCellItemSetMID int `xorm:"'training_tree_cell_item_set_m_id'"` // the set of items used to unlock this cell
		// SnsCoin int // always 1, maybe we could earn gem by unlocking cell?
	}

	cellContents := []TrainingTreeCellContent{}
	err = db.Table("m_training_tree_cell_content").Where("id = ?", treeMapping.TrainingTreeMappingMID).
		Cols("cell_id", "training_tree_cell_type", "training_content_no", "training_tree_cell_item_set_m_id").
		OrderBy("cell_id").Find(&cellContents)

	if err != nil {
		panic(err)
	}

	// stats reward is in "m_training_tree_card_param"
	type TrainingTreeCardParam struct {
		TrainingContentType int
		Value               int
	}
	cellParams := []TrainingTreeCardParam{}
	db.Table("m_training_tree_card_param").
		Where("id = ?", req.CardMasterID).
		Cols("training_content_type", "value").OrderBy("training_content_no").Find(&cellParams)

	increasedStats := [5]int{}
	card := session.GetUserCard(req.CardMasterID)
	for _, cellID := range req.CellMasterIDs {
		switch cellContents[cellID].TrainingTreeCellType {
		case 2:
			// param cells, reference the
			paramCell := &cellParams[cellContents[cellID].TrainingContentNo-1]
			increasedStats[paramCell.TrainingContentType] += paramCell.Value
		case 3:
			// voice
			// TODO: Unlock voice
		case 4:
			// story
			// TODO: Unlock story
		case 5:
			// idolize
			card.IsAwakening = true
			card.IsAwakeningImage = true
			// idolize stats is calculated by the flag?
		case 7:
			// skill
			card.ActiveSkillLevel += 1
		case 8:
			// insight
			card.MaxFreePassiveSkill += 1
		case 9:
			// ability
			card.PassiveSkillALevel += 1
		}
	}

	card.TrainingLife += increasedStats[2]
	card.TrainingAttack += increasedStats[3]
	card.TrainingDexterity += increasedStats[4]
	card.TrainingActivatedCellCount += len(req.CellMasterIDs)

	if card.TrainingActivatedCellCount+1 == len(cellContents) {
		card.IsAllTrainingActivated = true
	}

	session.UpdateUserCard(card)

	// set "user_card_training_tree_cell_list" to the cell unlocked and insert the cell to db
	unlockedCells := []model.TrainingTreeCell{}
	for _, cellID := range req.CellMasterIDs {
		cell := model.TrainingTreeCell{}
		cell.UserID = UserID
		cell.CardMasterID = req.CardMasterID
		cell.CellID = cellID
		cell.ActivatedAt = ClientTimeStamp
		unlockedCells = append(unlockedCells, cell)
	}

	session.InsertTrainingCells(&unlockedCells)
	jsonResp := session.Finalize(GetUserData("userModelDiff.json"), "user_model_diff")
	jsonResp, _ = sjson.Set(jsonResp, "user_card_training_tree_cell_list", session.GetTrainingTree(req.CardMasterID))
	resp := SignResp(ctx.GetString("ep"), jsonResp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
