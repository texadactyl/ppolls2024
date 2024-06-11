package helpers

import (
	"gopkg.in/yaml.v3"
	"os"
	"ppolls2024/global"
	"strconv"
)

type paramsStruct struct {
	// Embedded structs are not treated as embedded in YAML by default. To do that,
	// add the ",inline" annotation below
	PollHistoryLimit string `yaml:"PollHistoryLimit"`
	TossupThreshold  string `yaml:"TossupThreshold"`
	ECVAlgorithm     string `yaml:"ECVAlgorithm"`
}

func GetConfig() {
	var params paramsStruct
	glob := global.GetGlobalRef()
	bytes, err := os.ReadFile(glob.CfgFile)
	if err != nil {
		Croak("os.ReadFile(%s) failed, reason: %s", glob.CfgFile, err.Error())
	}
	err = yaml.Unmarshal(bytes, &params)
	if err != nil {
		Croak("yaml.Unmarshal from %s failed, reason: %s", glob.CfgFile, err.Error())
	}
	glob.PollHistoryLimit, err = strconv.Atoi(params.PollHistoryLimit)
	if err != nil {
		Croak("strconv.Atoi(PollHistoryLimit) from %s failed, reason: %s", glob.CfgFile, err.Error())
	}
	glob.TossupThreshold, err = strconv.ParseFloat(params.TossupThreshold, 64)
	if err != nil {
		Croak("strconv.ParseFloat(TossupThreshold) from %s failed, reason: %s", glob.CfgFile, err.Error())
	}
	glob.ECVAlgorithm, err = strconv.Atoi(params.ECVAlgorithm)
	if err != nil {
		Croak("strconv.Atoi(ECVAlgorithm) from %s failed, reason: %s", glob.CfgFile, err.Error())
	}
}
