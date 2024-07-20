package global

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Constants.
const PATH_VERSION = "./VERSION.txt"
const CSV_FILE_NAME = "president_poll.csv"
const INTERNET_FILE = "https://www.electoral-vote.com/evp2024/Pres/pres_polls.txt"

var DummyTime = time.Date(1776, time.July, 4, 23, 59, 59, 0, time.UTC)

// Definition of the singleton global.
type GlobalsStruct struct {
	Battleground     []string  // List of battleground states
	CfgFile          string    // Configuration file
	DateThreshold    time.Time // No polls used in reports nor plots before this date
	DbDriver         string    // Database driver name
	DbFile           string    // Database file name + extension
	DirCsv           string    // CSV input directory (before database load)
	DirDatabase      string    // Database directory path
	DirPlots         string    // Plots directory path
	DirTemp          string    // Temporary holding area directory path
	ECVAlgorithm     int       // Cfg: ECV distribution algorithm
	FlagBattleground bool      // Only report on battleground states (-r ec)? true/false
	FlagFetch        bool      // Fetch new data from the internet? true/false
	FlagLoad         bool      // Load new data into the database? true/false
	FlagPlot         bool      // Plots requested? true/false
	FlagReport       bool      // Report requested? true/false
	InternetCsvFile  string    // INTERNET_PREFIX + CSV_FILE_NAME + ".txt"
	LocalCsvFile     string    // CSV file name + extension
	PlotHeight       float64   // Height of plot canvase in dots
	PlotWidth        float64   // Width of plot canvase in dots
	PollHistoryLimit int       // Cfg: Limit of how many polls are entertained
	TossupThreshold  float64   // Cfg: Threshold of difference below which a tossup can be inferred
	Version          string    // Software version string
}

// Here's the singleton.
var global GlobalsStruct

// Initialise the singleton global and return a reference to it.
func InitGlobals() *GlobalsStruct {

	versionBytes, err := os.ReadFile(PATH_VERSION)
	if err != nil {
		fmt.Printf("\n*** InitGlobals: ReadFile(%s) failed, error: %s\n\n", PATH_VERSION, err.Error())
		os.Exit(1)
	}
	versionString := string(versionBytes[:])
	versionString = strings.TrimSpace(versionString)

	global = GlobalsStruct{
		CfgFile:          "config.yaml",
		DateThreshold:    DummyTime,
		DbFile:           "ppolls2024.db",
		DbDriver:         "sqlite",
		DirCsv:           "./csv/",
		DirDatabase:      "./database/",
		DirPlots:         "./plots/",
		DirTemp:          "./temp/",
		FlagFetch:        false,
		FlagLoad:         false,
		FlagReport:       false,
		FlagPlot:         false,
		FlagBattleground: false,
		InternetCsvFile:  INTERNET_FILE,
		LocalCsvFile:     CSV_FILE_NAME,
		Version:          versionString,
	}

	return &global
}

// GetGlobalRef returns a pointer to the singleton instance of GlobalsStruct
func GetGlobalRef() *GlobalsStruct {
	return &global
}
