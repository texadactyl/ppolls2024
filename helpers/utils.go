package helpers

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"ppolls2024/global"
	"strings"
	"time"
)

const ModeOutputFile = 0644
const FlagsOpen = os.O_CREATE | os.O_WRONLY

// GetUtcDate - Get UTC date string, YYYY-MM-DD
func GetUtcDate() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02")
}

// GetUtcTime - Get UTC time string in the format of hh:mm:ss.ddd.
func GetUtcTime() string {
	now := time.Now().UTC()
	return now.Format("15:04:05.000")
}

// Date (YYYY-MM-DD) to Time.
func YYYY_MM_DDtoTime(dateString string) (time.Time, error) {
	tm, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return global.DummyTime, err
	}
	return tm, err
}

// WriteOutputText - Write a text line to the given output file handle.
func WriteOutputText(outHandle *os.File, textLine string) {

	_, err := fmt.Fprintln(outHandle, textLine)
	if err != nil {
		outPath, _ := filepath.Abs(filepath.Dir(outHandle.Name()))
		log.Fatalf("WriteOutputText: fmt.Fprintln(%s) failed, reason: %s\n", outPath, err.Error())
	}

}

// MakeDir - If the specified directory does not yet exist, create it.
func MakeDir(pathDir string) {
	info, err := os.Stat(pathDir)
	if err == nil { // found it
		if !info.IsDir() { // expected a directory, not a simple file !!
			log.Fatal("MakeDir: Observed a simple file: %s (expected a directory)\n", pathDir)
		}
	} else { // not found or an error occurred
		if os.IsNotExist(err) {
			// Create directory
			err = os.Mkdir(pathDir, 0755)
			if err != nil {
				log.Fatalf("MakeDir: os.MkDir(%s) failed, reason: %s\n", pathDir, err)
			}
			log.Printf("MakeDir: %s was created\n", pathDir)
		} else { // some type of error
			log.Fatalf("MakeDir: os.Stat(%s) failed, reason: %s\n", pathDir, err.Error())
		}
	}
}

// StoreText - Store a text file in the specified directory.
func StoreText(targetDir string, argFile string, text string) {
	// Create the log file
	fullPath := filepath.Join(targetDir, argFile)
	outHandle, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("storeText: os.Create(%s) failed, reason: %s\n", fullPath, err.Error())
	}
	defer outHandle.Close()

	// Store the given text
	_, err = fmt.Fprintln(outHandle, text)
	if err != nil {
		log.Fatalf("storeText: fmt.Fprintln(%s) failed, reason: %s\n", fullPath, err.Error())
	}
}

// CleanerText - Replace nongraphics with '?'.
func CleanerText(argText string) string {
	rr := []rune(argText)
	for ii := 0; ii < len(rr); ii++ {
		if rr[ii] == '\n' || rr[ii] == '\r' || rr[ii] == '\t' {
			continue
		}
		if rr[ii] > 126 || rr[ii] < 32 {
			// Less than ASCII Space or greater than ASCII ~
			rr[ii] = 63 // ='?'
		}
	}
	return string(rr)
}

// Given a 3-character month code, return its month number.
var monthList = []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

func MonthToInt(monthStr string) (int, error) {
	arg := strings.ToUpper(monthStr)
	for ii := 0; ii < 12; ii++ {
		if arg == monthList[ii] {
			return ii + 1, nil
		}
	}
	errMsg := fmt.Sprintf("invalid month string (%s)", monthStr)
	return -1, errors.New(errMsg)
}

// Given a state, return the ECV for that state.
func StateToECV(state string) int {
	arg := strings.ToUpper(state)
	for ii := 0; ii < len(global.StateTable); ii++ {
		if arg == global.StateTable[ii].Stcode {
			return global.StateTable[ii].Votes
		}
	}
	log.Fatalf(fmt.Sprintf("StateToECV: invalid state code (%s)", state))
	return -1 // dummy return
}

// Calculate not Dem nor Gop.
func CalcOther(dem, gop float64) float64 {
	return 100.0 - (dem + gop)
}

// Trend calculation.
// If less than 3 points, return "--".
// If rising, return "up".
// If falling, return "dn".
// Otherwise, return "--".
func CalcTrend(array []float64) string {
	num := len(array)
	if num < 3 {
		return "--"
	}
	if array[num-1] > array[num-2] {
		if array[num-2] > array[num-3] {
			return "u2"
		}
		return "u1"
	}
	if array[num-1] < array[num-2] {
		if array[num-2] < array[num-3] {
			return "d2"
		}
		return "d1"
	}
	return "--"
}

// Translate true | false into "**" | "  ".
func getFactorString(arg bool) string {
	if arg {
		return "**"
	}
	return "  "
}

/*
ECVAward1 - ECV Award Algorithm 1.

	Calculate the Other percentage = 100 - the sum of the candidate percentages.
	Split the "Other" percentage proportionally amongst the candidates.
	Calculate the difference = absolute value of the difference between the candidates.

	If the difference is below the tossup threshold, then this state is a tossup.
*/
func ECVAward1(stateVotes int, pctDem, pctGop float64) (string, int, int, int) {
	glob := global.GetGlobalRef()
	pctOther := CalcOther(pctDem, pctGop)

	pctDem += pctOther * pctDem / 100.0
	pctGop += pctOther * pctGop / 100.0
	diff := math.Abs(pctDem - pctGop)

	if diff < glob.TossupThreshold {
		return "TOSSUP", 0, 0, stateVotes
	} else {
		if pctDem > pctGop {
			return "Dem", stateVotes, 0, 0
		} else {
			return "Gop", 0, stateVotes, 0
		}
	}
}

/*
ECVAward2 - ECV Award Algorithm 2.

	Calculate the Other percentage = 100 - the sum of the candidate percentages.
	Calculate the difference = absolute value of the difference between the candidates.

	If the "Other" percentage exceeds the difference, then flag this state on return.
	If the difference is below the tossup threshold, then this state is a tossup.
*/
func ECVAward2(stateVotes int, pctDem, pctGop float64) (string, int, int, int, string) {
	glob := global.GetGlobalRef()
	diff := math.Abs(pctDem - pctGop)
	pctOther := CalcOther(pctDem, pctGop)
	if diff < glob.TossupThreshold {
		return "TOSSUP", 0, 0, stateVotes, getFactorString(pctOther > diff)
	}
	if pctDem > pctGop {
		return "Dem", stateVotes, 0, 0, getFactorString(pctOther > diff)
	}
	return "Gop", 0, stateVotes, 0, getFactorString(pctOther > diff)
}

/*
ECVAward3 - ECV Award Algorithm 3.

	Calculate the Other percentage = 100 - the sum of the candidate percentages.
	Calculate the difference = absolute value of the difference between the candidates.

	If the "Other" percentage exceeds the difference, then this state is a tossup.
	If the difference is below the tossup threshold, then this state is a tossup.
*/
func ECVAward3(stateVotes int, pctDem, pctGop float64) (string, int, int, int, string) {
	glob := global.GetGlobalRef()
	diff := math.Abs(pctDem - pctGop)
	pctOther := CalcOther(pctDem, pctGop)
	if pctOther > diff {
		return "TOSSUP", 0, 0, stateVotes, getFactorString(true)
	}
	if diff < glob.TossupThreshold {
		return "TOSSUP", 0, 0, stateVotes, " "
	}
	if pctDem > pctGop {
		return "Dem", stateVotes, 0, 0, " "
	}
	return "Gop", 0, stateVotes, 0, " "
}

// searchSlice looks for a target string in an array of strings.
// It returns true if the target string is found,
// Otherwise, it returns false.
func searchSlice(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
