package helpers

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"ppolls2024/global"
	"strconv"
	"strings"
)

type paramsStruct struct {
	// Embedded structs are not treated as embedded in YAML by default. To do that,
	// add the ",inline" annotation below
	Battleground     string `yaml:"Battleground"`
	DateThreshold    string `yaml:"DateThreshold"`
	ECVAlgorithm     string `yaml:"ECVAlgorithm"`
	PollHistoryLimit string `yaml:"PollHistoryLimit"`
	PlotHeight       string `yaml:"PlotHeight"`
	PlotWidth        string `yaml:"PlotWidth"`
	TossupThreshold  string `yaml:"TossupThreshold"`
}

func GetConfig() {

	var params paramsStruct
	glob := global.GetGlobalRef()
	bytes, err := os.ReadFile(glob.CfgFile)
	if err != nil {
		log.Fatalf("GetConfig: os.ReadFile(%s) failed, reason: %s\n", glob.CfgFile, err.Error())
	}
	err = yaml.Unmarshal(bytes, &params)
	if err != nil {
		log.Fatalf("yaml.Unmarshal from %s failed, reason: %s\n", glob.CfgFile, err.Error())
	}

	glob.DateThreshold, err = YYYY_MM_DDtoTime(params.DateThreshold)
	if err != nil {
		log.Fatalf("GetConfig: Cannot parse start date: %s, reason: %s\n", params.DateThreshold, err.Error())
	}
	log.Printf("GetConfig: DateThreshold: %s", params.DateThreshold)

	glob.ECVAlgorithm, err = strconv.Atoi(params.ECVAlgorithm)
	if err != nil {

		log.Fatalf("strconv.Atoi(ECVAlgorithm) from %s failed, reason: %s\n", glob.CfgFile, err.Error())
	}
	log.Printf("GetConfig: ECVAlgorithm: %d", glob.ECVAlgorithm)

	glob.Battleground = strings.Split(params.Battleground, ",")
	log.Printf("GetConfig: Battleground states: %s", params.Battleground)

	glob.PlotWidth, err = strconv.ParseFloat(params.PlotWidth, 64)
	if err != nil {
		log.Fatalf("GetConfig: strconv.ParseFloat(PlotWidth) from %s failed, reason: %s\n", glob.CfgFile, err.Error())
	}
	log.Printf("GetConfig: PlotWidth: %f", glob.PlotWidth)

	glob.PlotHeight, err = strconv.ParseFloat(params.PlotHeight, 64)
	if err != nil {
		log.Fatalf("GetConfig: strconv.ParseFloat(PlotHeight) from %s failed, reason: %s\n", glob.CfgFile, err.Error())
	}
	log.Printf("GetConfig: PlotHeight: %f", glob.PlotHeight)

	glob.PollHistoryLimit, err = strconv.Atoi(params.PollHistoryLimit)
	if err != nil {
		log.Fatalf("strconv.Atoi(PollHistoryLimit) from %s failed, reason: %s\n", glob.CfgFile, err.Error())
	}
	log.Printf("GetConfig: PollHistoryLimit: %d", glob.PollHistoryLimit)

	glob.TossupThreshold, err = strconv.ParseFloat(params.TossupThreshold, 64)
	if err != nil {
		log.Fatalf("GetConfig: strconv.ParseFloat(TossupThreshold) from %s failed, reason: %s\n", glob.CfgFile, err.Error())
	}
	log.Printf("GetConfig: TossupThreshold: %f", glob.TossupThreshold)

}
