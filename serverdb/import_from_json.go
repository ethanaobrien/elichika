package serverdb

import (
	"elichika/config"
	"elichika/model"
	"elichika/utils"

	// "os"
	"encoding/json"
	"fmt"
	"strings"

	// "xorm.io/xorm"

	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

var (
	presetDataPath = "assets/preset/"
	userDataPath   = "assets/userdata/"
	UserID         = 588296696
	IsGlobal       = true
)

func GetJsonData(fileName string) string {
	userDataFile := userDataPath + fileName
	if utils.PathExists(userDataFile) {
		return utils.ReadAllText(userDataFile)
	}

	presetDataFile := presetDataPath + fileName
	if !utils.PathExists(presetDataFile) {
		panic("File not exists")
	}

	userData := utils.ReadAllText(presetDataFile)
	utils.WriteAllText(userDataFile, userData)

	return userData
}

func GetLiveDeckData() string {
	if IsGlobal {
		return GetJsonData("liveDeck_gl.json")
	}
	return GetJsonData("liveDeck.json")
}

func CreateNewUser() {
	fmt.Println("Insert new user: ", UserID)
	data := GetJsonData("userStatus.json")
	status := model.UserStatus{}
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		panic(err)
	}
	status.UserID = UserID
	// insert into the db
	_, err := Engine.Table("s_user_info").AllCols().Insert(&status)
	if err != nil {
		panic(err)
	}
}

func LoadMemberFromJson() []model.UserMemberInfo {
	members := []model.UserMemberInfo{}
	userMemberInfo := model.UserMemberInfo{}
	memberData := gjson.Parse(GetJsonData("memberSettings.json"))
	memberData.Get("user_member_by_member_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &userMemberInfo); err != nil {
				panic(err)
			}
			userMemberInfo.UserID = UserID
			members = append(members, userMemberInfo)
		}
		return true
	})
	return members
}

func LoadLiveDeckAndLivePartyFromJson() ([]model.UserLiveDeck, []model.UserLiveParty) {
	fmt.Println("importing live deck data to db")
	liveDecks := []model.UserLiveDeck{}
	liveDeckInfo := model.UserLiveDeck{}
	liveDeckData := gjson.Parse(GetLiveDeckData())
	liveDeckData.Get("user_live_deck_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &liveDeckInfo); err != nil {
				panic(err)
			}
			liveDeckInfo.UserID = UserID
			liveDecks = append(liveDecks, liveDeckInfo)
		}
		return true
	})

	liveParties := []model.UserLiveParty{}

	fmt.Println("importing live party data to db")
	livePartyInfo := model.UserLiveParty{}
	var livePartyData []json.RawMessage
	decoder := json.NewDecoder(strings.NewReader(liveDeckData.Get("user_live_party_by_id").String()))
	decoder.UseNumber()
	err := decoder.Decode(&livePartyData)
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(livePartyData); i += 2 {
		err := json.Unmarshal(livePartyData[i], &livePartyInfo)
		if err != nil {
			panic(err)
		}
		livePartyInfo.UserID = UserID
		liveParties = append(liveParties, livePartyInfo)
	}
	return liveDecks, liveParties
}

func LoadLessonDeckFromJson() []model.UserLessonDeck {
	lessonDecks := []model.UserLessonDeck{}
	fmt.Println("importing lesson deck data to db")
	lessonData := gjson.Parse(GetJsonData("lessonDeck.json"))
	userLessonDeck := model.UserLessonDeck{}
	lessonData.Get("user_lesson_deck_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &userLessonDeck); err != nil {
				panic(err)
			}
			userLessonDeck.UserID = UserID
			lessonDecks = append(lessonDecks, userLessonDeck)
		}
		return true
	})
	return lessonDecks
}

func LoadCardFromJson() []model.UserCard {
	cards := []model.UserCard{}
	fmt.Println("importing json card data to db")
	cardData := gjson.Parse(GetJsonData("userCard.json"))
	userCard := model.UserCard{}
	cardData.Get("user_card_by_card_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &userCard); err != nil {
				panic(err)
			}
			userCard.UserID = UserID
			cards = append(cards, userCard)
		}
		return true
	})
	return cards
}

func LoadSuitFromJson() []model.UserSuit {
	anys := []any{}
	suits := []model.UserSuit{}
	fmt.Println("importing json suit data to db")
	err := json.Unmarshal([]byte(gjson.Get(GetJsonData("login.json"), "user_model.user_suit_by_suit_id").String()), &anys)
	if err != nil {
		panic(err)
	}
	n := len(anys)
	for i := 0; i < n; i += 2 {
		suits = append(suits, model.UserSuit{
			UserID:       UserID,
			SuitMasterID: int(anys[i].(float64)),
			IsNew:        false})
	}
	// fmt.Println(suits)
	return suits
}

func (session *Session) InsertAccount(members []model.UserMemberInfo,
	liveDecks []model.UserLiveDeck,
	liveParties []model.UserLiveParty,
	lessonDecks []model.UserLessonDeck,
	cards []model.UserCard,
	suits []model.UserSuit,
	lovePanels []model.UserMemberLovePanel,
	trainingTreeCells []model.TrainingTreeCell) {

	session.InsertMembers(members)
	session.InsertLiveDecks(liveDecks)
	session.InsertLiveParties(liveParties)
	session.InsertLessonDecks(lessonDecks)
	session.InsertCards(cards)
	InsertUserSuits(suits)
	for _, lovePanel := range lovePanels {
		_, err := Engine.Table("s_user_member").AllCols().
			Where("user_id = ? AND member_master_id = ?", lovePanel.UserID, lovePanel.MemberID).
			Update(&lovePanel)
		if err != nil {
			panic(err)
		}
	}
	session.InsertTrainingCells(&trainingTreeCells)
}

func ImportFromJson() {
	fmt.Println("Importing account data from json")
	CreateNewUser()
	session := GetSession(UserID)

	// member data
	members := LoadMemberFromJson()
	lovePanels := []model.UserMemberLovePanel{}
	for _, member := range members {
		lovePanel := model.UserMemberLovePanel{
			UserID:                 member.UserID,
			MemberID:               member.MemberMasterID,
			MemberLovePanelCellIDs: []int{}}

		for boardLevel := 0; boardLevel <= 31; boardLevel++ {
			for tile := 1; tile <= 5; tile++ {
				cellID := boardLevel*10000 + tile*1000 + member.MemberMasterID
				lovePanel.MemberLovePanelCellIDs = append(lovePanel.MemberLovePanelCellIDs, cellID)
			}
		}
		lovePanels = append(lovePanels, lovePanel)
	}

	liveDecks, liveParties := LoadLiveDeckAndLivePartyFromJson()
	lessonDecks := LoadLessonDeckFromJson()
	cards := LoadCardFromJson()
	suits := LoadSuitFromJson()
	trainingTreeCells := []model.TrainingTreeCell{}
	for _, card := range cards {
		for cellID := 1; cellID <= 89; cellID++ {
			// this is a bit lazy, could have used the correct number instead
			trainingTreeCells = append(trainingTreeCells,
				model.TrainingTreeCell{
					UserID:       card.UserID,
					CardMasterID: card.CardMasterID,
					CellID:       cellID,
					ActivatedAt:  int64(1688094000 + cellID)})
		}
	}
	session.InsertAccount(members, liveDecks, liveParties, lessonDecks, cards, suits, lovePanels, trainingTreeCells)
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

func ImportMinimalAccount() {
	// EXPERIMENTAL, for testing functionality only.
	// insert an account that can be upgraded to maximum strength from existing functionality
	// For now, it's something like this:
	// - all the same with the json account
	// - all card are level 1, grade 0 (no limit break), no tile unlocked, skill and ability level 1, no insight skill unlocked
	// - all members have max bond level such that when all of their cards are limit break, they get max bond level 500
	// - all members have bond level 1 (no bond)
	fmt.Println("Creating minimal data from json")
	CreateNewUser()
	session := GetSession(UserID)

	// member data
	members := LoadMemberFromJson()
	lovePanels := []model.UserMemberLovePanel{}
	memberIDToIndex := make(map[int]int)
	maxBondLevel := [30]int{}
	suitID := [30]int{}
	for i, _ := range members {
		memberIDToIndex[members[i].MemberMasterID] = i
		maxBondLevel[i] = 500
		members[i].LovePoint = 0
		members[i].LoveLevel = 1
	}

	liveDecks, liveParties := LoadLiveDeckAndLivePartyFromJson()
	lessonDecks := LoadLessonDeckFromJson()
	cards := LoadCardFromJson()
	suits := LoadSuitFromJson()
	// skip the suit awarded from cards
	suitCount := 0
	for ; suits[suitCount].SuitMasterID < 1000000; suitCount++ {
	}
	suits = suits[:suitCount]

	for i, _ := range cards {
		memberID := (cards[i].CardMasterID / 10000) % 1000
		rarity := (cards[i].CardMasterID / 100) % 100
		maxBondLevel[memberIDToIndex[memberID]] -= (rarity / 10) * 5 // 5 grade worths of limit break
		cards[i].Level = 1
		cards[i].IsAwakening = false
		cards[i].IsAwakeningImage = false
		cards[i].IsAllTrainingActivated = false
		cards[i].TrainingActivatedCellCount = 0
		exists, err := config.MainEng.Table("m_card").Where("id = ?", cards[i].CardMasterID).
			Cols("passive_skill_slot").Get(&cards[i].MaxFreePassiveSkill)
		if err != nil {
			panic(err)
		}
		if !exists {
			panic("not exist")
		}
		cards[i].Grade = 0
		cards[i].TrainingLife = 0
		cards[i].TrainingAttack = 0
		cards[i].TrainingDexterity = 0
		cards[i].ActiveSkillLevel = 1
		cards[i].PassiveSkillALevel = 1
		cards[i].AdditionalPassiveSkill1ID = 0
		cards[i].AdditionalPassiveSkill2ID = 0
		cards[i].AdditionalPassiveSkill3ID = 0
		cards[i].AdditionalPassiveSkill4ID = 0
	}

	for _, suit := range suits {
		memberID := (suit.SuitMasterID / 100) % 1000
		suitID[memberIDToIndex[memberID]] = suit.SuitMasterID
	}

	trainingTreeCells := []model.TrainingTreeCell{}
	for i, _ := range members {
		members[i].LovePointLimit = BondRequiredTotal(maxBondLevel[i])
		members[i].SuitMasterID = suitID[i]
	}
	session.InsertAccount(members, liveDecks, liveParties, lessonDecks, cards, suits, lovePanels, trainingTreeCells)
}
