package config

import (
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	// Taken from:
	// https://github.com/RayFirefist/SukuStar_Datamine/blob/master/lib/sifas_api/sifas.py#L120
	// https://github.com/RayFirefist/SukuStar_Datamine/blob/master/lib/sifas_api/sifas.py#L400
	ServerEventReceiverKey = "31f1f9dc7ac4392d1de26acf99d970e425b63335b461e720c73d6914020d6014"
	JaKey                  = "78d53d9e645a0305602174e06b98d81f638eaf4a84db19c756866fddac360c96"

	SessionKey = "12345678123456781234567812345678"

	MainDb         = "assets/main.db"
	MasterdataDbGl = "assets/db/masterdata_gl.db"
	MasterdataDbJp = "assets/db/masterdata_jp.db"
	ServerdataDb   = "assets/db/serverdata.db"

	MainEng         *xorm.Engine
	MasterdataEngGl *xorm.Engine
	MasterdataEngJp *xorm.Engine

	Conf = &AppConfigs{}

	PresetDataPath = "assets/preset/"
	UserDataPath   = "assets/userdata/"
)

func init() {
	Conf = Load("./config.json")

	eng, err := xorm.NewEngine("sqlite", MainDb)
	if err != nil {
		panic(err)
	}
	err = eng.Ping()
	if err != nil {
		panic(err)
	}
	MainEng = eng
	MainEng.SetMaxOpenConns(50)
	MainEng.SetMaxIdleConns(10)

	eng1, err := xorm.NewEngine("sqlite", MasterdataDbGl)
	if err != nil {
		panic(err)
	}
	MasterdataEngGl = eng1
	MasterdataEngGl.SetMaxOpenConns(50)
	MasterdataEngGl.SetMaxIdleConns(10)

	eng2, err := xorm.NewEngine("sqlite", MasterdataDbJp)
	if err != nil {
		panic(err)
	}
	MasterdataEngJp = eng2
	MasterdataEngJp.SetMaxOpenConns(50)
	MasterdataEngJp.SetMaxIdleConns(10)
}
