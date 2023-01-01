package gamedata

import (
	"elichika/config"
	"elichika/dictionary"
	"elichika/enum"
	"elichika/model"
	"elichika/utils"

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

func (ld *LiveDifficulty) populate(gamedata *Gamedata, masterdata_db, serverdata_db *xorm.Session, dictionary *dictionary.Dictionary) {
	ld.Live = gamedata.Live[*ld.LiveID]
	// 2-way links
	ld.Live.LiveDifficulties = append(ld.Live.LiveDifficulties, ld)
	ld.LiveID = &gamedata.Live[*ld.LiveID].LiveID
	err := masterdata_db.Table("m_live_difficulty_mission").Where("live_difficulty_master_id = ?", ld.LiveDifficultyID).
		OrderBy("position").Find(&gamedata.LiveDifficulty[ld.LiveDifficultyID].Missions)
	utils.CheckErr(err)
	// if ld.LiveDifficultyID == 9999 || ld.LiveDifficultyID/10 == 6000000  {
	// 	return
	// }

	ld.LiveDifficultyGimmick = new(LiveDifficultyGimmick)
	exists, err := masterdata_db.Table("m_live_difficulty_gimmick").Where("live_difficulty_master_id = ?", ld.LiveDifficultyID).
		Get(ld.LiveDifficultyGimmick)
	utils.CheckErr(err)

	if !exists {
		// doesn't exist for a small set of things that shouldn't matter
		// panic(fmt.Sprint("gimmick doesn't exist for: ", ld.LiveDifficultyID))
		ld.LiveDifficultyGimmick = nil
		// fmt.Println("gimmick doesn't exist for: ", ld.LiveDifficultyID)
	}

	err = masterdata_db.Table("m_live_difficulty_note_gimmick").Where("live_difficulty_id = ?", ld.LiveDifficultyID).
		Find(&ld.LiveDifficultyNoteGimmicks)
	utils.CheckErr(err)
	for i := range ld.LiveDifficultyNoteGimmicks {
		ld.LiveDifficultyNoteGimmicks[i].populate()
	}
}

func (ld *LiveDifficulty) loadSimpleLiveStage(gamedata *Gamedata) {
	if ld.SimpleLiveStage != nil {
		return // already loaded
	}
	// fmt.Println("Loading for", ld.LiveDifficultyID)
	liveNotes := utils.ReadAllText(fmt.Sprintf("assets/simple_stages/%d.json", ld.LiveDifficultyID))
	if (liveNotes == "") || (ld.UnlockPattern == enum.LiveUnlockPatternTowerOnly) {

		// song doesn't exist, use rule to find the original map
		if ld.UnlockPattern != enum.LiveUnlockPatternTowerOnly {
			// only accept event songs, SBL, or DLP
			return
		}
		originalLiveID := ld.Live.LiveID%10000 + 10000
		for _, other := range gamedata.Live[originalLiveID].LiveDifficulties {
			if (other.NoteEmitMsec == ld.NoteEmitMsec) && (other.LiveDifficultyType == ld.LiveDifficultyType) {
				other.loadSimpleLiveStage(gamedata)
				if other.SimpleLiveStage != nil {
					ld.SimpleLiveStage = other.SimpleLiveStage
					break
				}
			}
		}
		if ld.SimpleLiveStage == nil {
			for _, other := range gamedata.Live[originalLiveID].LiveDifficulties {
				if other.NoteEmitMsec == ld.NoteEmitMsec {
					other.loadSimpleLiveStage(gamedata)
					if other.SimpleLiveStage != nil {
						ld.SimpleLiveStage = other.SimpleLiveStage
						break
					}
				}
			}
		}
	} else {
		err := json.Unmarshal([]byte(liveNotes), &ld.SimpleLiveStage)
		utils.CheckErr(err)
	}
	if ld.SimpleLiveStage == nil {
		panic(fmt.Sprint("Error finding live stage for: ", ld.LiveDifficultyID))
	}
	if ld.SimpleLiveStage.Original != nil {
		_, exists := gamedata.LiveDifficulty[*ld.SimpleLiveStage.Original]
		if !exists {
			fmt.Println("Warning: original live referenced but do not exist in database: ",
				*ld.SimpleLiveStage.Original, ". Attemping to just load the json.")
			gamedata.LiveDifficulty[*ld.SimpleLiveStage.Original] = new(LiveDifficulty)
			gamedata.LiveDifficulty[*ld.SimpleLiveStage.Original].LiveDifficultyID = *ld.SimpleLiveStage.Original
			gamedata.LiveDifficulty[*ld.SimpleLiveStage.Original].LiveDifficultyType = ld.LiveDifficultyType
		}
		gamedata.LiveDifficulty[*ld.SimpleLiveStage.Original].loadSimpleLiveStage(gamedata)
		ld.SimpleLiveStage = gamedata.LiveDifficulty[*ld.SimpleLiveStage.Original].SimpleLiveStage
	}
	if ld.SimpleLiveStage == nil {
		panic(fmt.Sprint("Error finding original live stage for: ", ld.LiveDifficultyID))
	}
}

func (ld *LiveDifficulty) ConstructLiveStage(gamedata *Gamedata) {
	if ld.LiveStage != nil { // generated
		return
	}

	if !config.GenerateStageFromScratch { // load generated stage, it must exists
		text := utils.ReadAllText(fmt.Sprintf("assets/stages/%d.json", ld.LiveDifficultyID))
		if text == "" {
			panic(fmt.Sprintf("Stage %d doesn't exists in assets/stages", ld.LiveDifficultyID))
		}
		ld.LiveStage = new(model.LiveStage)
		err := json.Unmarshal([]byte(text), &ld.LiveStage)
		if err != nil {
			panic(fmt.Sprintf("Failed to load stage %d: wrong format", ld.LiveDifficultyID))
		}
		return
	}

	ld.loadSimpleLiveStage(gamedata)
	if ld.SimpleLiveStage == nil {
		if ld.UnlockPattern != enum.LiveUnlockPatternTowerOnly {
			return
		}
		panic(fmt.Sprint("Failed to load simple live stage for: ", ld.LiveDifficultyID))
	}

	// make the object and set relevant stuff
	ld.LiveStage = new(model.LiveStage)
	ld.LiveStage.LiveDifficultyID = ld.LiveDifficultyID
	ld.LiveStage.LiveNotes = []model.LiveNote{}
	ld.LiveStage.NoteGimmicks = []model.NoteGimmick{}
	ld.LiveStage.LiveWaveSettings = []model.LiveWaveSetting{}
	ld.LiveStage.StageGimmickDict = []any{}

	ld.LiveStage.LiveNotes = append(ld.LiveStage.LiveNotes, ld.SimpleLiveStage.LiveNotes...)
	for i := range ld.LiveStage.LiveNotes {
		ld.LiveStage.LiveNotes[i].ID = i + 1
		ld.LiveStage.LiveNotes[i].AutoJudgeType = enum.JudgeTypeGreat         // can be overwritten at runtime
		ld.LiveStage.LiveNotes[i].NoteRandomDropColor = enum.NoteDropColorNon // can be overwritten at runtime
	}
	ld.LiveStage.LiveWaveSettings = append(ld.LiveStage.LiveWaveSettings, ld.SimpleLiveStage.LiveWaveSettings...)

	// each note store its own gimmick, and the stage store unique note gimmicks in it
	noteGimmickDict := map[int]bool{}
	for _, noteGimmick := range ld.LiveDifficultyNoteGimmicks {
		ld.LiveStage.LiveNotes[noteGimmick.NoteID-1].GimmickID = noteGimmick.ID
		if !noteGimmickDict[noteGimmick.ID] {
			noteGimmickDict[noteGimmick.ID] = true
			ld.LiveStage.NoteGimmicks = append(ld.LiveStage.NoteGimmicks,
				model.NoteGimmick{
					ID:              noteGimmick.ID,
					NoteGimmickType: noteGimmick.NoteGimmickType,
					EffectMID:       noteGimmick.SkillMasterID,
					IconType:        noteGimmick.NoteGimmickIconType,
				})
		}
	}
	sort.Slice(ld.LiveStage.NoteGimmicks, func(i, j int) bool {
		return ld.LiveStage.NoteGimmicks[i].ID < ld.LiveStage.NoteGimmicks[j].ID
	})
	for i := range ld.LiveStage.NoteGimmicks {
		ld.LiveStage.NoteGimmicks[i].UniqID = 2001 + i
	}
	if ld.LiveDifficultyGimmick != nil {
		ld.LiveStage.StageGimmickDict = append(ld.LiveStage.StageGimmickDict, ld.LiveDifficultyGimmick.TriggerType)
		ld.LiveStage.StageGimmickDict = append(ld.LiveStage.StageGimmickDict, []model.StageGimmick{model.StageGimmick{
			GimmickMasterID:    ld.LiveDifficultyGimmick.ID,
			ConditionMasterID1: ld.LiveDifficultyGimmick.ConditionMasterID1,
			ConditionMasterID2: ld.LiveDifficultyGimmick.ConditionMasterID2,
			SkillMasterID:      ld.LiveDifficultyGimmick.SkillMasterID,
			UniqID:             1001,
		}})
	}

	// save the new map
	{
		output, err := json.Marshal(ld.LiveStage)
		utils.CheckErr(err)
		utils.WriteAllText(fmt.Sprintf("assets/stages/%d.json", ld.LiveDifficultyID), string(output))
	}

	// check against pregenerated map
	// skip checking for coop (SBL), because the database only has constant modifier while the actual
	// data will have some added bonus gimmick
	// not like we use those map right now anyway
	if ld.UnlockPattern == enum.LiveUnlockPatternCoopOnly {
		return
	}
	text := utils.ReadAllText(fmt.Sprintf("assets/full_stages/%d.json", ld.LiveDifficultyID))
	if text == "" {
		// fmt.Println("Newly generated map: ", ld.LiveDifficultyID)
		return
	}
	pregeneratedStage := model.LiveStage{}
	err := json.Unmarshal([]byte(text), &pregeneratedStage)
	utils.CheckErr(err)
	if !pregeneratedStage.IsSame(ld.LiveStage) {
		validDiff := map[int]bool{}
		if !validDiff[ld.LiveDifficultyID] {
			panic(fmt.Sprint("Difference detected for: ", ld.LiveDifficultyID, "\n", ld.LiveStage, "\n_______________\n", pregeneratedStage))
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
