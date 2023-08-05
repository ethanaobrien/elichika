package serverdb

import (
	"elichika/config"
	"elichika/model"

	"fmt"

	"os"
	"xorm.io/xorm"
)

var (
	Engine *xorm.Engine
	IsNew  bool = false
)

func InitTable(tableName string, structure interface{}) bool {
	exist, err := Engine.Table(tableName).IsTableExist(tableName)
	if err != nil {
		panic(err)
	}

	if !exist {
		fmt.Println("Creating new table: ", tableName)
		err = Engine.Table(tableName).CreateTable(structure)
		if err != nil {
			panic(err)
		}
		return true
	} else {
		return false
	}
}

func InitTables() bool {
	type DbUser struct {
		model.UserStatus             `xorm:"extends"`
		model.DBUserProfileLiveStats `xorm:"extends"`
	}
	isNew := false

	isNew = InitTable("s_user_info", DbUser{})
	type DbCard struct {
		model.UserCard       `xorm:"extends"`
		model.DBCardPlayInfo `xorm:"extends"`
	}
	isNew = InitTable("s_user_card", DbCard{}) || isNew
	isNew = InitTable("s_user_suit", model.UserSuit{}) || isNew
	isNew = InitTable("s_user_training_tree_cell", model.TrainingTreeCell{}) || isNew

	type DbMembers struct {
		model.UserMemberInfo   `xorm:"extends"`
		MemberLovePanelCellIDs []int `xorm:"'member_love_panel_cell_ids'"`
	}
	isNew = InitTable("s_user_member", DbMembers{}) || isNew

	isNew = InitTable("s_user_lesson_deck", model.UserLessonDeck{}) || isNew

	isNew = InitTable("s_user_live_deck", model.UserLiveDeck{}) || isNew
	isNew = InitTable("s_user_live_party", model.UserLiveParty{}) || isNew
	isNew = InitTable("s_user_live_state", model.LiveState{}) || isNew
	return isNew
}

func InitDb(isGlobal bool) {
	if IsNew { // init the db depend on argv
		IsGlobal = isGlobal
		if len(os.Args) == 1 { // import from existing jsons
			ImportFromJson()
		} else {
			ImportMinimalAccount() // make a minimal account
		}
		IsNew = false
	}
}

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite", config.ServerdataDb)
	if err != nil {
		panic(err)
	}
	Engine.SetMaxOpenConns(50)
	Engine.SetMaxIdleConns(10)
	// Engine.ShowSQL(true)
	IsNew = InitTables()
}
