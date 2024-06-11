package global

import (
	"fmt"
	"os"
	"strings"
)

// Constants.
const PATH_VERSION = "./VERSION.txt"
const CSV_FILE_NAME = "president_poll.csv"
const INTERNET_FILE = "https://www.electoral-vote.com/evp2024/Pres/pres_polls.txt"

// Definition of the singleton global.
type GlobalsStruct struct {
	Version          string  // Software version string
	InternetCsvFile  string  // INTERNET_PREFIX + CSV_FILE_NAME + ".txt"
	LocalCsvFile     string  // CSV file name + extension
	DirCsvIn         string  // CSV input directory (before database load)
	DbDriver         string  // Database driver name
	DbDirectory      string  // Database directory path
	DbFile           string  // Database file name + extension
	CfgFile          string  // Configuration file
	PollHistoryLimit int     // Cfg: Limit of how many polls are entertained
	TossupThreshold  float64 // Cfg: Threshold of difference below which a tossup can be inferred
	ECVAlgorithm     int     // Cfg: ECV distribution algorithm
	FlagFetch        bool    // Fetch new data from the internet? true/false
	FlagLoad         bool    // Load new data into the database? true/false
	FlagReport       bool    // Report requested? true/false
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
		Version:          versionString,
		InternetCsvFile:  INTERNET_FILE,
		LocalCsvFile:     CSV_FILE_NAME,
		DirCsvIn:         "./csv/",
		DbDriver:         "sqlite",
		DbDirectory:      "./database/",
		DbFile:           "ppolls2024.db",
		CfgFile:          "config.yaml",
		PollHistoryLimit: -42,  // dummy value
		TossupThreshold:  42.0, // dummy value
		ECVAlgorithm:     42,   // dummy value
		FlagFetch:        false,
		FlagLoad:         false,
		FlagReport:       false,
	}

	return &global
}

// GetGlobalRef returns a pointer to the singleton instance of GlobalsStruct
func GetGlobalRef() *GlobalsStruct {
	return &global
}
