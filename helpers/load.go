package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Load(dirCsv, fileName string) {
	var pollTable []string
	var pollFields dbparams

	// Form the full path of the CSV file.
	fullPath := filepath.Join(dirCsv, fileName)

	// Get the ECV table file contents as a slice of strings.
	fileBytes, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("Load: ReadFile(%s) failed, reason: %s\n", fullPath, err.Error())
	}

	// Create a table of strings.
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
			log.Fatalf("Load: File %s malformed at line %d, ncols: %d\n", fullPath, lineCounter, len(colArray))
		}

		// Collect all the column values.
		pollFields.state = strings.ToUpper(colArray[0])
		pollFields.pctDem, err = strconv.ParseFloat(colArray[1], 64)
		if err != nil {
			log.Fatalf("Load: Dem pct from %s is not a valid float at line %d\n", fullPath, lineCounter)
		}
		pollFields.pctGop, err = strconv.ParseFloat(colArray[2], 64)
		if err != nil {
			log.Fatalf("Load: Gop pct from %s is not a valid float at line %d\n", fullPath, lineCounter)
		}
		month, err := MonthToInt(colArray[4])
		if err != nil {
			log.Fatalf("Load: start month from %s at line %d: %s\n", fullPath, lineCounter, err.Error())
		}
		day, err := strconv.ParseInt(colArray[5], 10, 64)
		if err != nil {
			log.Fatalf("Load: start day from %s is not a valid integer at line %d\n", fullPath, lineCounter)
		}
		pollFields.startDate = fmt.Sprintf("2024-%02d-%02d", month, day)
		month, err = MonthToInt(colArray[6])
		if err != nil {
			log.Fatalf("Load: end month from %s at line %d: %s\n", fullPath, lineCounter, err.Error())
		}
		day, err = strconv.ParseInt(colArray[7], 10, 64)
		if err != nil {
			log.Fatalf("Load: end day from %s is not a valid integer at line %d\n", fullPath, lineCounter)
		}
		pollFields.endDate = fmt.Sprintf("2024-%02d-%02d", month, day)
		pollFields.pollster = strings.Join(colArray[8:], " ")

		// Insert this database row.
		DBStore(pollFields)
	}

	log.Printf("Loaded %d records into the database\n", lineCounter)
}
