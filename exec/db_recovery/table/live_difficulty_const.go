package table

import (
	"elichika/exec/db_recovery/parser"
	"elichika/utils"

	"strconv"
)

type LiveDifficultyConst struct {
}

func (_ *LiveDifficultyConst) Table() string {
	return "m_live_difficulty_const"
}
func (_ *LiveDifficultyConst) ID(fields []parser.Field) int64 {
	if fields[0].Key != "id" {
		panic("wrong field order")
	}
	id, err := strconv.ParseInt(fields[0].Value, 10, 64)
	utils.CheckErr(err)
	return id
}
func (_ *LiveDifficultyConst) Value(field parser.Field) string {
	return field.Value
}
func (this *LiveDifficultyConst) Update(field parser.Field) string {
	return field.Key + "=" + this.Value(field)
}
func (this *LiveDifficultyConst) Condition(fields []parser.Field) string {
	return this.Update(fields[0])
}

func handleLiveDifficultyConstEvent(event parser.ModifierEvent[LiveDifficultyConst]) {
	if event.Type == parser.DELETE { // if deleted then we can add it back
		event.Type = parser.INSERT
	} else if event.Type == parser.INSERT {
		return
	}
	output += event.String() + "\n"
}

func handleLiveDifficultyConst() {
	var dummy LiveDifficultyConst
	events := parser.Parse[LiveDifficultyConst](readGitChange(dummy.Table()))
	for _, event := range events {
		handleLiveDifficultyConstEvent(event)
	}
}

func init() {
	addHandler(handleLiveDifficultyConst)
}
