package gamedata

import (
	"elichika/dictionary"
	"elichika/enum"
	"elichika/model"
	"elichika/utils"
	"elichika/config"

	"encoding/json"
	"fmt"
	"sort"

	"xorm.io/xorm"
)

type LiveDifficulty struct {
	// from m_live_difficulty
	LiveDifficultyID int   `xorm:"pk 'live_difficulty_id'"`
	LiveID           *int  `xorm:"'live_id'"`
	Live             *Live `xorm:"-"`
	// Live3DAssetMasterID *int
	LiveDifficultyType int `xorm:"'live_difficulty_type'"`
	UnlockPattern      int `xorm:"'unlock_pattern'"`
	// DefaultAttribute int
	TargetVoltage int `xorm:"'target_voltage'"`
	NoteEmitMsec  int `xorm:"'note_emit_msec'"`
	ConsumedLP    int `xorm:"'consumed_lp'"`
	RewardUserExp int `xorm:"'reward_user_exp'"`
	// JudgeID int
	NoteDropGroupID *int `xorm:"'note_drop_group_id'"`

	// NoteDropGroup *NoteDropGroup `xorm:"-"`
	DropChooseCount    int  `xorm:"'drop_choose_count'"`
	RateDropRate       int  `xorm:"'rare_drop_rate'"`
	DropContentGroupID *int `xorm:"'drop_content_group_id'"`
	// DropContentGroup *DropContentGroup `xorm:"-"`
	RareDropContentGroupID *int `xorm:"'rare_drop_content_group_id'"`
	// RareDropContentGroup *RareDropContentGroup `xorm:"-"`
	AdditionalDropContentGroupID *int `xorm:"'additional_drop_content_group_id'"`
	// AdditionalDropContentGroup *AdditionalDropContentGroup `xorm:"-"`
	// ?????
	BottomTechnique              int `xorm:"'bottom_technique'"`
	AdditionalDropDecayTechnique int `xorm:"'additional_drop_decay_technique'"`

	RewardBaseLovePoint int `xorm:"'reward_base_love_point'"`
	EvaluationSScore    int `xorm:"'evaluation_s_score'"`
	EvaluationAScore    int `xorm:"'evaluation_a_score'"`
	EvaluationBScore    int `xorm:"'evaluation_b_score'"`
	EvaluationCScore    int `xorm:"'evaluation_c_score'"`
	// UpdatedAt int `xorm:"'updated_at'"`
	LoseAtDeath bool `xorm:"'lose_at_death'"`
	// AutoplayRequirementID *int `xorm:"'autoplay_requirement_id'"`
	SkipMasterID *int `xorm:"'skip_master_id'"`
	// StaminaVoltageGroupID int
	// ComboVoltageGroupID int
	// DifficultyConstMasterID int
	IsCountTarget bool `xorm:"'is_count_target'"`
	// InsufficentRate int

	// from m_live_difficulty_mission
	Missions []LiveDifficultyMission `xorm:"-"`

	// lazily constructed?
	LiveStage       *model.LiveStage `xorm:"-"`
	SimpleLiveStage *SimpleLiveStage `xorm:"-"`

	// from m_live_difficulty_gimmick
	LiveDifficultyGimmick *LiveDifficultyGimmick `xorm:"-"`

	// from m_live_difficulty_note_gimmick
	LiveDifficultyNoteGimmicks []LiveDifficultyNoteGimmick `xorm:"-"`
}

func (this *LiveDifficulty) populate(gamedata *Gamedata, masterdata_db, serverdata_db *xorm.Session, dictionary *dictionary.Dictionary) {
	this.Live = gamedata.Live[*this.LiveID]
	// 2-way links
	this.Live.LiveDifficulties = append(this.Live.LiveDifficulties, this)
	this.LiveID = &gamedata.Live[*this.LiveID].LiveID
	err := masterdata_db.Table("m_live_difficulty_mission").Where("live_difficulty_master_id = ?", this.LiveDifficultyID).
		OrderBy("position").Find(&gamedata.LiveDifficulty[this.LiveDifficultyID].Missions)
	utils.CheckErr(err)
	// if this.LiveDifficultyID == 9999 || this.LiveDifficultyID/10 == 6000000  {
	// 	return
	// }

	this.LiveDifficultyGimmick = new(LiveDifficultyGimmick)
	exists, err := masterdata_db.Table("m_live_difficulty_gimmick").Where("live_difficulty_master_id = ?", this.LiveDifficultyID).
		Get(this.LiveDifficultyGimmick)
	utils.CheckErr(err)

	if !exists {
		// doesn't exist for a small set of things that shouldn't matter
		// panic(fmt.Sprint("gimmick doesn't exist for: ", this.LiveDifficultyID))
		this.LiveDifficultyGimmick = nil
		// fmt.Println("gimmick doesn't exist for: ", this.LiveDifficultyID)
	}

	err = masterdata_db.Table("m_live_difficulty_note_gimmick").Where("live_difficulty_id = ?", this.LiveDifficultyID).
		Find(&this.LiveDifficultyNoteGimmicks)
	utils.CheckErr(err)
	for i := range this.LiveDifficultyNoteGimmicks {
		this.LiveDifficultyNoteGimmicks[i].populate()
	}
}

func (this *LiveDifficulty) loadSimpleLiveStage(gamedata *Gamedata) {
	if this.SimpleLiveStage != nil {
		return // already loaded
	}
	// fmt.Println("Loading for", this.LiveDifficultyID)
	liveNotes := utils.ReadAllText(fmt.Sprintf("assets/simple_stages/%d.json", this.LiveDifficultyID))
	if (liveNotes == "") || (this.UnlockPattern == enum.LiveUnlockPatternTowerOnly) {

		// song doesn't exist, use rule to find the original map
		if (this.UnlockPattern != enum.LiveUnlockPatternTowerOnly){
			// only accept event songs, SBL, or DLP
			return
		}
		originalLiveID := this.Live.LiveID%10000 + 10000
		for _, other := range gamedata.Live[originalLiveID].LiveDifficulties {
			if (other.NoteEmitMsec == this.NoteEmitMsec) &&  (other.LiveDifficultyType == this.LiveDifficultyType) {
				other.loadSimpleLiveStage(gamedata)
				if other.SimpleLiveStage != nil {
					this.SimpleLiveStage = other.SimpleLiveStage
					break
				}
			}
		}
		if this.SimpleLiveStage == nil {
			for _, other := range gamedata.Live[originalLiveID].LiveDifficulties {
				if other.NoteEmitMsec == this.NoteEmitMsec {
					other.loadSimpleLiveStage(gamedata)
					if other.SimpleLiveStage != nil {
						this.SimpleLiveStage = other.SimpleLiveStage
						break
					}
				}
			}
		}
	} else {
		err := json.Unmarshal([]byte(liveNotes), &this.SimpleLiveStage)
		utils.CheckErr(err)
	}
	if this.SimpleLiveStage == nil {
		panic(fmt.Sprint("Error finding live stage for: ", this.LiveDifficultyID))
	}
	if this.SimpleLiveStage.Original != nil {
		_, exists := gamedata.LiveDifficulty[*this.SimpleLiveStage.Original]
		if !exists {
			fmt.Println("Warning: original live referenced but do not exist in database: ",
				*this.SimpleLiveStage.Original, ". Attemping to just load the json.")
			gamedata.LiveDifficulty[*this.SimpleLiveStage.Original] = new(LiveDifficulty)
			gamedata.LiveDifficulty[*this.SimpleLiveStage.Original].LiveDifficultyID = *this.SimpleLiveStage.Original
			gamedata.LiveDifficulty[*this.SimpleLiveStage.Original].LiveDifficultyType = this.LiveDifficultyType
		}
		gamedata.LiveDifficulty[*this.SimpleLiveStage.Original].loadSimpleLiveStage(gamedata)
		this.SimpleLiveStage = gamedata.LiveDifficulty[*this.SimpleLiveStage.Original].SimpleLiveStage
	}
	if this.SimpleLiveStage == nil {
		panic(fmt.Sprint("Error finding original live stage for: ", this.LiveDifficultyID))
	}
}

func (this *LiveDifficulty) ConstructLiveStage(gamedata *Gamedata) {
	if this.LiveStage != nil { // generated
		return
	}

	if !config.GenerateStageFromScratch { // load generated stage, it must exists 
		text := utils.ReadAllText(fmt.Sprintf("assets/stages/%d.json", this.LiveDifficultyID))
		if text == "" {
			panic(fmt.Sprintf("Stage %d doesn't exists in assets/stages"))
		}
		this.LiveStage = new(model.LiveStage)
		err := json.Unmarshal([]byte(text), &this.LiveStage)
		if err != nil {
			panic(fmt.Sprintf("Failed to load stage %d: wrong format"))
		}
		return 
	}

	this.loadSimpleLiveStage(gamedata)
	if this.SimpleLiveStage == nil {
		if this.UnlockPattern != enum.LiveUnlockPatternTowerOnly {
			return
		}
		panic(fmt.Sprint("Failed to load simple live stage for: ", this.LiveDifficultyID))
	}

	// make the object and set relevant stuff
	this.LiveStage = new(model.LiveStage)
	this.LiveStage.LiveDifficultyID = this.LiveDifficultyID
	this.LiveStage.LiveNotes = []model.LiveNote{}
	this.LiveStage.NoteGimmicks = []model.NoteGimmick{}
	this.LiveStage.LiveWaveSettings = []model.LiveWaveSetting{}
	this.LiveStage.StageGimmickDict = []any{}

	this.LiveStage.LiveNotes = append(this.LiveStage.LiveNotes, this.SimpleLiveStage.LiveNotes...)
	for i := range this.LiveStage.LiveNotes {
		this.LiveStage.LiveNotes[i].ID = i + 1
		this.LiveStage.LiveNotes[i].AutoJudgeType = enum.JudgeTypeGreat // can be overwritten at runtime
		this.LiveStage.LiveNotes[i].NoteRandomDropColor = enum.NoteDropColorNon // can be overwritten at runtime
	}
	this.LiveStage.LiveWaveSettings = append(this.LiveStage.LiveWaveSettings, this.SimpleLiveStage.LiveWaveSettings...)

	// each note store its own gimmick, and the stage store unique note gimmicks in it
	noteGimmickDict := map[int]bool{}
	for _, noteGimmick := range this.LiveDifficultyNoteGimmicks {
		this.LiveStage.LiveNotes[noteGimmick.NoteID-1].GimmickID = noteGimmick.ID
		if !noteGimmickDict[noteGimmick.ID] {
			noteGimmickDict[noteGimmick.ID] = true
			this.LiveStage.NoteGimmicks = append(this.LiveStage.NoteGimmicks,
				model.NoteGimmick{
					ID:              noteGimmick.ID,
					NoteGimmickType: noteGimmick.NoteGimmickType,
					EffectMID:       noteGimmick.SkillMasterID,
					IconType:        noteGimmick.NoteGimmickIconType,
				})
		}
	}
	sort.Slice(this.LiveStage.NoteGimmicks, func(i, j int) bool {
		return this.LiveStage.NoteGimmicks[i].ID < this.LiveStage.NoteGimmicks[j].ID
	})
	for i := range this.LiveStage.NoteGimmicks {
		this.LiveStage.NoteGimmicks[i].UniqID = 2001 + i
	}
	if this.LiveDifficultyGimmick != nil {
		this.LiveStage.StageGimmickDict = append(this.LiveStage.StageGimmickDict, this.LiveDifficultyGimmick.TriggerType)
		this.LiveStage.StageGimmickDict = append(this.LiveStage.StageGimmickDict, []model.StageGimmick{model.StageGimmick{
			GimmickMasterID:    this.LiveDifficultyGimmick.ID,
			ConditionMasterID1: this.LiveDifficultyGimmick.ConditionMasterID1,
			ConditionMasterID2: this.LiveDifficultyGimmick.ConditionMasterID2,
			SkillMasterID:      this.LiveDifficultyGimmick.SkillMasterID,
			UniqID:             1001,
		}})
	}

	// save the new map
	{
		output, err := json.Marshal(this.LiveStage)
		utils.CheckErr(err)
		utils.WriteAllText(fmt.Sprintf("assets/stages/%d.json", this.LiveDifficultyID), string(output))
	}

	// check against pregenerated map
	// skip checking for coop (SBL), because the database only has constant modifier while the actual 
	// data will have some added bonus gimmick
	// not like we use those map right now anyway
	if this.UnlockPattern == enum.LiveUnlockPatternCoopOnly {
		return
	}
	text := utils.ReadAllText(fmt.Sprintf("assets/full_stages/%d.json", this.LiveDifficultyID))
	if text == "" {
		// fmt.Println("Newly generated map: ", this.LiveDifficultyID)
		return
	}
	pregeneratedStage := model.LiveStage{}
	err := json.Unmarshal([]byte(text), &pregeneratedStage)
	utils.CheckErr(err)
	if !pregeneratedStage.IsSame(this.LiveStage) {
		validDiff := map[int]bool{}
		if !validDiff[this.LiveDifficultyID] {
			panic(fmt.Sprint("Difference detected for: ", this.LiveDifficultyID, "\n", this.LiveStage, "\n_______________\n", pregeneratedStage))
		}
	}

}

func loadLiveDifficulty(gamedata *Gamedata, masterdata_db, serverdata_db *xorm.Session, dictionary *dictionary.Dictionary) {
	fmt.Println("Loading LiveDifficulty")
	gamedata.LiveDifficulty = make(map[int]*LiveDifficulty)
	err := masterdata_db.Table("m_live_difficulty").Find(&gamedata.LiveDifficulty)
	utils.CheckErr(err)
	// ordered iteration is important here
	ids := []int{}
	for id := range gamedata.LiveDifficulty {
		ids = append(ids, id)
	}
	// order by unlock pattern then id
	sort.Slice(ids, func(i, j int) bool {
		if gamedata.LiveDifficulty[ids[i]].UnlockPattern != gamedata.LiveDifficulty[ids[j]].UnlockPattern {
			return gamedata.LiveDifficulty[ids[i]].UnlockPattern < gamedata.LiveDifficulty[ids[j]].UnlockPattern
		}
		return gamedata.LiveDifficulty[ids[i]].LiveDifficultyID < gamedata.LiveDifficulty[ids[j]].LiveDifficultyID
	})
	for _, id := range ids {
		liveDifficulty := gamedata.LiveDifficulty[id]
		liveDifficulty.populate(gamedata, masterdata_db, serverdata_db, dictionary)
	}

	if config.GenerateStageFromScratch {
		for _, liveDifficulty := range gamedata.LiveDifficulty {
			liveDifficulty.ConstructLiveStage(gamedata)
		}
	}
}

func init() {
	addLoadFunc(loadLiveDifficulty)
	addPrequisite(loadLiveDifficulty, loadLive)
}
