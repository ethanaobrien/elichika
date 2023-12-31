package enum

const (
	NoteTypeSingleTap         = 0x00000001
	NoteTypeLongStart         = 0x00000002
	NoteTypeLongEnd           = 0x00000003
	NoteTypeWaveStart         = 0x00000004
	NoteTypeWaveEnd           = 0x00000005
	NoteTypeTutorialSingleTap = 0x00000006
)

const (
	JudgeTypeNone    = 0x00000001
	JudgeTypeMiss    = 0x0000000a
	JudgeTypeBad     = 0x0000000c
	JudgeTypeGood    = 0x0000000e
	JudgeTypeGreat   = 0x00000014
	JudgeTypePerfect = 0x0000001e
)

const (
	NoteDropColorGold   = 0x00000001
	NoteDropColorSilver = 0x00000002
	NoteDropColorBronze = 0x00000003
	NoteDropColorNon    = 0x00000063
)
