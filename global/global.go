package global

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	CfgFile          string    // Configuration file path
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
	PollHistoryLimit int       // Limit of how many polls are entertained
	StateTableFile   string    // State table file path
	StronglyDem      []string  // List of strongly Democratic states
	StronglyGop      []string  // List of strongly GOP states
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
		StateTableFile:   "state_table.txt",
		Version:          versionString,
	}

	bytes, err := os.ReadFile(global.StateTableFile)
	if err != nil {
		log.Fatalf("InitGlobals: os.ReadFile(%s) failed, reason: %s\n", global.StateTableFile, err.Error())
	}
	allTheLines := string(bytes)
	lineSplice := strings.Split(allTheLines, "\n")
	lineCount := 0
	for _, line := range lineSplice {
		lineCount++
		line = strings.TrimSpace(line)
		if len(line) < 1 || strings.HasPrefix(line, "#") {
			continue
		}
		triplet := strings.Fields(line)
		if len(triplet) != 3 {
			log.Fatalf("InitGlobals: State table line %d is not a triplet\n", lineCount)
		}
		votesValue, err := strconv.Atoi(triplet[1])
		if err != nil {
			log.Fatalf("InitGlobals: strconv.Atoi(Votes) failed on line %d, reason: %s\n", lineCount, err.Error())
		}
		StateTable = append(StateTable, StateTableEntry_t{Stcode: triplet[0], Votes: votesValue, Category: triplet[2]})
		switch triplet[2] {
		case "B":
			global.Battleground = append(global.Battleground, triplet[0])
		case "D":
			global.StronglyDem = append(global.StronglyDem, triplet[0])
		case "G":
			global.StronglyGop = append(global.StronglyGop, triplet[0])
		default:
			log.Fatalf("InitGlobals: State Category on line %d is not supported: %s\n", lineCount, triplet[2])
		}
	}
	log.Printf("InitGlobals: Battleground states: %v (%d)\n", global.Battleground, len(global.Battleground))
	log.Printf("InitGlobals: Strongly Democrat states: %v (%d)\n", global.StronglyDem, len(global.StronglyDem))
	log.Printf("InitGlobals: Strongly GOP states: %v (%d)\n", global.StronglyGop, len(global.StronglyGop))

	return &global
}

// GetGlobalRef returns a pointer to the singleton instance of GlobalsStruct
func GetGlobalRef() *GlobalsStruct {
	return &global
}

// State table entry definition
type StateTableEntry_t struct {
	Stcode   string // 2-character state code
	Votes    int    // number of Electoral College Votes
	Category string // "B" (battleground), "D" (strongly democrat), or "G" (strongly GOP)
}

// State table
var StateTable = []StateTableEntry_t{}
