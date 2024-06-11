package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func loadOneCsv(absPath string) {
	var pollTable []string
	var pollFields dbparams

	// Get the EV table file contents as a slice of strings.
	fileBytes, err := os.ReadFile(absPath)
	if err != nil {
		Croak("loadOneCsv: ReadFile(%s) failed, reason: %s", absPath, err.Error())
	}
	giantString := string(fileBytes)
	pollTable = strings.Split(string(giantString), "\n")

	// For each line, parse the state string and the EV count.
	lineCounter := 0
	for _, oneLine := range pollTable {
		lineCounter += 1
		oneLine = strings.TrimSpace(oneLine)
		if len(oneLine) < 1 {
			continue
		}
		colArray := strings.Split(oneLine, " ")
		if len(colArray) < 1 {
			continue
		}
		if strings.HasPrefix(colArray[0], "#") {
			continue
		}
		if len(colArray) < 9 {
			Croak("loadOneCsv: File %s malformed at line %d, ncols: %d", absPath, lineCounter, len(colArray))
		}

		// Collect all the column values.
		pollFields.state = strings.ToUpper(colArray[0])
		pollFields.pctBiden, err = strconv.ParseFloat(colArray[1], 64)
		if err != nil {
			Croak("loadOneCsv: Biden pct from %s is not a valid float at line %d", absPath, lineCounter)
		}
		pollFields.pctTrump, err = strconv.ParseFloat(colArray[2], 64)
		if err != nil {
			Croak("loadOneCsv: Trump pct from %s is not a valid float at line %d", absPath, lineCounter)
		}
		month, err := MonthToInt(colArray[4])
		if err != nil {
			Croak("loadOneCsv: start month from %s at line %d: %s", absPath, lineCounter, err.Error())
		}
		day, err := strconv.ParseInt(colArray[5], 10, 64)
		if err != nil {
			Croak("loadOneCsv: start day from %s is not a valid integer at line %d", absPath, lineCounter)
		}
		pollFields.startDate = fmt.Sprintf("2024-%02d-%02d", month, day)
		month, err = MonthToInt(colArray[6])
		if err != nil {
			Croak("loadOneCsv: end month from %s at line %d: %s", absPath, lineCounter, err.Error())
		}
		day, err = strconv.ParseInt(colArray[7], 10, 64)
		if err != nil {
			Croak("loadOneCsv: end day from %s is not a valid integer at line %d", absPath, lineCounter)
		}
		pollFields.endDate = fmt.Sprintf("2024-%02d-%02d", month, day)
		pollFields.pollster = strings.Join(colArray[8:], " ")

		// Insert this database row.
		DBStore(pollFields)
	}

	Logger("Loaded %d records into the database", lineCounter)

}

func Load(dirCsvIn string) {
	entries, err := os.ReadDir(dirCsvIn)
	if err != nil {
		Croak("Load: os.ReadDir(%s) failed, reason: %s", dirCsvIn, err.Error())
	}
	for _, dirEntry := range entries {
		fileName := dirEntry.Name()

		// Get full path of file name
		fullPathFile := filepath.Join(dirCsvIn, fileName)

		// If a directory, we will skip it
		fileInfo, err := os.Stat(fullPathFile)
		if err != nil {
			Croak("Load: os.Stat(%s) failed, reason: %s", fullPathFile, err.Error())
		}
		if fileInfo.IsDir() {
			continue
		}

		// Not a directory.
		// If not a .csv file, skip it
		if filepath.Ext(fileName) != ".csv" {
			continue
		}

		loadOneCsv(fullPathFile)
	}

}
