package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"ppolls2024/global"
	"ppolls2024/helpers"
	"strings"
)

// Show help and then exit to the O/S
func showHelp() {
	suffix := filepath.Base(os.Args[0])
	fmt.Printf("\nUsage:  %s  {-f  -l  -p  -r ID}\n\nwhere\n\n", suffix)
	fmt.Printf("\t-f:\tFetch latest poll data from Internet --> directory csv\n")
	fmt.Printf("\t-l:\tLoad poll data from directory csv\n")
	fmt.Printf("\t-p:\tGenerate plots\n")
	fmt.Printf("\t-r ID:\tReport by identifier (ID):\n")
	fmt.Printf("\t\tSC\tSC = state code (E.g. AL).\n")
	fmt.Printf("\t\tEC\tElectoral College tallies for all states.\n")
	fmt.Printf("\t-b:\tProcess only battleground states in -r ec\n")
	fmt.Printf("\nExit codes:\n")
	fmt.Printf("\t0\tNormal completion or help shown due to command line error.\n")
	fmt.Printf("\t1\tSomething went wrong during execution.\n\n")
	os.Exit(0)
}

// Walk the base JMOD file and count the bytes for all classes/*.class entries
func main() {

	var params []string
	rpt := ""
	glob := global.InitGlobals()
	helpers.GetConfig()

	// Parse command line arguments.
	for _, singleVar := range os.Args[1:] {
		params = append(params, singleVar)
	}
	if len(params) < 1 {
		fmt.Printf("*** Missing arguments!\n")
		showHelp()
	}

	for ii := 0; ii < len(params); ii++ {
		switch params[ii] {
		case "-h":
			showHelp()
		case "-f":
			glob.FlagFetch = true
		case "-l":
			glob.FlagLoad = true
		case "-p":
			glob.FlagPlot = true
		case "-r":
			ii++
			if ii >= len(params) {
				fmt.Println("*** The -r parameter lacks a value!")
				showHelp()
			}
			rpt = strings.ToUpper(params[ii])
			glob.FlagReport = true
		case "-b":
			glob.FlagBattleground = true
		default:
			fmt.Printf("*** The specified parameter (%s) is not supported!\n", params[ii])
			showHelp()
		}
	}

	// If plotting, delete old plots.
	if glob.FlagPlot {
		err := os.RemoveAll(glob.DirPlots)
		if err != nil {
			log.Fatalf(fmt.Sprintf("os.Remove(%s) failed, reason: %s", glob.DirPlots, err.Error()))
		}
	}

	// Create subdirectories.
	helpers.MakeDir(glob.DirCsv)
	helpers.MakeDir(glob.DirDatabase)
	helpers.MakeDir(glob.DirPlots)
	helpers.MakeDir(glob.DirTemp)

	// Validate the use of -b.
	if glob.FlagBattleground && !glob.FlagReport {
		log.Println("Warning: No reports requested. The battleground flag (-b) is ignored")
	}

	// Fetch new data?
	if glob.FlagFetch {
		if !helpers.Fetch(glob.DirCsv, glob.LocalCsvFile, glob.InternetCsvFile, glob.DirTemp) {
			os.Exit(1)
		}
	}

	// Load newly-fetched data into the database?
	if glob.FlagLoad {
		helpers.DBOpen(glob.DbDriver, glob.DirDatabase, glob.DbFile)
		helpers.Load(glob.DirCsv, glob.LocalCsvFile)
		helpers.DBClose()
	}

	// Generate plots?
	if glob.FlagPlot {
		helpers.DBOpen(glob.DbDriver, glob.DirDatabase, glob.DbFile)
		helpers.Plodder()
		helpers.DBClose()
	}

	// Run a report?
	if glob.FlagReport {
		helpers.DBOpen(glob.DbDriver, glob.DirDatabase, glob.DbFile)
		if rpt == "EC" {
			helpers.ReportEC()
		} else {
			helpers.ReportSC(rpt)
		}
		helpers.DBClose()
	}

}
