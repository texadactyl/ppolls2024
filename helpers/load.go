package helpers

import (
	"fmt"
	"log"
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
		log.Fatalf("loadOneCsv: ReadFile(%s) failed, reason: %s\n", absPath, err.Error())
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
			log.Fatalf("loadOneCsv: File %s malformed at line %d, ncols: %d\n", absPath, lineCounter, len(colArray))
		}

		// Collect all the column values.
		pollFields.state = strings.ToUpper(colArray[0])
		pollFields.pctBiden, err = strconv.ParseFloat(colArray[1], 64)
		if err != nil {
			log.Fatalf("loadOneCsv: Biden pct from %s is not a valid float at line %d\n", absPath, lineCounter)
		}
		pollFields.pctTrump, err = strconv.ParseFloat(colArray[2], 64)
		if err != nil {
			log.Fatalf("loadOneCsv: Trump pct from %s is not a valid float at line %d\n", absPath, lineCounter)
		}
		month, err := MonthToInt(colArray[4])
		if err != nil {
			log.Fatalf("loadOneCsv: start month from %s at line %d: %s\n", absPath, lineCounter, err.Error())
		}
		day, err := strconv.ParseInt(colArray[5], 10, 64)
		if err != nil {
			log.Fatalf("loadOneCsv: start day from %s is not a valid integer at line %d\n", absPath, lineCounter)
		}
		pollFields.startDate = fmt.Sprintf("2024-%02d-%02d", month, day)
		month, err = MonthToInt(colArray[6])
		if err != nil {
			log.Fatalf("loadOneCsv: end month from %s at line %d: %s\n", absPath, lineCounter, err.Error())
		}
		day, err = strconv.ParseInt(colArray[7], 10, 64)
		if err != nil {
			log.Fatalf("loadOneCsv: end day from %s is not a valid integer at line %d\n", absPath, lineCounter)
		}
		pollFields.endDate = fmt.Sprintf("2024-%02d-%02d", month, day)
		pollFields.pollster = strings.Join(colArray[8:], " ")

		// Insert this database row.
		DBStore(pollFields)
	}

	log.Printf("Loaded %d records into the database\n", lineCounter)

}

func Load(dirCsvIn string) {
	entries, err := os.ReadDir(dirCsvIn)
	if err != nil {
		log.Fatalf("Load: os.ReadDir(%s) failed, reason: %s\n", dirCsvIn, err.Error())
	}
	for _, dirEntry := range entries {
		fileName := dirEntry.Name()

		// Get full path of file name
		fullPathFile := filepath.Join(dirCsvIn, fileName)

		// If a directory, we will skip it
		fileInfo, err := os.Stat(fullPathFile)
		if err != nil {
			log.Fatalf("Load: os.Stat(%s) failed, reason: %s\n", fullPathFile, err.Error())
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
