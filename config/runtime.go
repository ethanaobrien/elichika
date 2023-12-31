package config

import (
	"elichika/enum"
	"elichika/utils"

	"encoding/json"
)

type RuntimeConfig struct {
	CdnServer     *string `json:"cdn_server"`
	ServerAddress *string `json:"server_address"`
	TapBondGain   *int    `json:"tap_bond_gain"`
	AutoJudgeType *int    `json:"auto_judge_type"`
}

func defaultConfigs() *RuntimeConfig {
	configs := RuntimeConfig{
		CdnServer:     new(string), // self-hosted
		ServerAddress: new(string),
		TapBondGain:   new(int),
		AutoJudgeType: new(int),
	}
	*configs.CdnServer = "https://llsifas.catfolk.party/static/"
	*configs.ServerAddress = "0.0.0.0:8080"
	*configs.TapBondGain = 20
	*configs.AutoJudgeType = enum.JudgeTypeGreat
	return &configs
}

func Load(p string) *RuntimeConfig {
	if !utils.PathExists(p) {
		c := defaultConfigs()
		c.Save(p)
		return c
	}

	c := RuntimeConfig{}
	err := json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		panic("config file is wrong, change/delete it and try again")
	}
	d := defaultConfigs()

	if c.CdnServer == nil {
		c.CdnServer = d.CdnServer
	}
	if c.ServerAddress == nil {
		c.ServerAddress = d.ServerAddress
	}
	if c.TapBondGain == nil {
		c.TapBondGain = d.TapBondGain
	}
	if c.AutoJudgeType == nil {
		c.AutoJudgeType = d.AutoJudgeType
	}

	return &c
}

func (c *RuntimeConfig) Save(p string) error {
	data, err := json.Marshal(c)
	utils.CheckErr(err)
	utils.WriteAllText(p, string(data))
	return nil
}
