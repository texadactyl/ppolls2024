package main

import (
	"fmt"
	"os"
	"path/filepath"
	"ppolls2024/global"
	"ppolls2024/helpers"
	"strings"
)

// Show help and then exit to the O/S
func showHelp() {
	suffix := filepath.Base(os.Args[0])
	fmt.Printf("\nUsage:  %s  {-f  -l  -v   -r ID}\n\nwhere\n\n", suffix)
	fmt.Printf("\t-f:\tFetch latest poll data from Internet --> directory csv\n")
	fmt.Printf("\t-l:\tLoad poll data from directory csv\n")
	fmt.Printf("\t-r ID:\tReport by identifier (ID):\n")
	fmt.Printf("\t\tSC\tSC = state code (E.g. AL).\n")
	fmt.Printf("\t\tEC\tElectoral College outcome.\n")
	fmt.Printf("\t-p:\tGenerate plots\n\n")
	fmt.Printf("Exit codes:\n")
	fmt.Printf("\t0\tNormal completion.\n")
	fmt.Printf("\t1\tSomething went wrong during execution.\n\n")
	os.Exit(1)
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
			rpt = strings.ToUpper(params[ii])
			glob.FlagReport = true
		default:
			fmt.Printf("*** The specified parameter (%s) is not supported!\n", params[ii])
			showHelp()
		}
	}

	// Create subdirectories.
	helpers.MakeDir(glob.DirCsv)
	helpers.MakeDir(glob.DirDatabase)
	helpers.MakeDir(glob.DirPlots)
	helpers.MakeDir(glob.DirTemp)

	// Fetch new data?
	if glob.FlagFetch {
		helpers.Fetch(glob.DirCsv, glob.LocalCsvFile, glob.InternetCsvFile, glob.DirTemp)
	}

	// Load newly-fetched data into the database?
	if glob.FlagLoad {
		helpers.DBOpen(glob.DbDriver, glob.DirDatabase, glob.DbFile)
		helpers.Load(glob.DirCsv, glob.LocalCsvFile)
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
	}

	// Generate plots?
	if glob.FlagPlot {
		helpers.DBOpen(glob.DbDriver, glob.DirDatabase, glob.DbFile)
		helpers.Plodder()
	}

}
