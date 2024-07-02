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

// Given a state code, return the EC vote value.
type ECVote struct {
	state string
	votes int
}

var stateECVTable = []ECVote{
	{"AL", 9},
	{"AK", 3},
	{"AZ", 11},
	{"AR", 6},
	{"CA", 54},
	{"CO", 10},
	{"CT", 7},
	{"DE", 3},
	{"DC", 3},
	{"FL", 30},
	{"GA", 16},
	{"HI", 4},
	{"ID", 4},
	{"IL", 19},
	{"IN", 11},
	{"IA", 6},
	{"KS", 6},
	{"KY", 8},
	{"LA", 8},
	{"ME", 4},
	{"MD", 10},
	{"MA", 11},
	{"MI", 15},
	{"MN", 10},
	{"MS", 6},
	{"MO", 10},
	{"MT", 4},
	{"NE", 5},
	{"NV", 6},
	{"NH", 4},
	{"NJ", 14},
	{"NM", 5},
	{"NY", 28},
	{"NC", 16},
	{"ND", 3},
	{"OH", 17},
	{"OK", 7},
	{"OR", 8},
	{"PA", 19},
	{"RI", 4},
	{"SC", 9},
	{"SD", 3},
	{"TN", 11},
	{"TX", 40},
	{"UT", 6},
	{"VT", 3},
	{"VA", 13},
	{"WA", 12},
	{"WV", 4},
	{"WI", 10},
	{"WY", 3},
}

// Given a state, return the ECV for that state.
func StateToECV(state string) int {
	arg := strings.ToUpper(state)
	for ii := 0; ii < len(stateECVTable); ii++ {
		if arg == stateECVTable[ii].state {
			return stateECVTable[ii].votes
		}
	}
	log.Fatalf(fmt.Sprintf("StateToECV: invalid state code (%s)", state))
	return -1 // dummy return
}

// Calculate not Biden nor Trump.
func CalcOther(biden, trump float64) float64 {
	return 100.0 - (biden + trump)
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

// Translate true/false into "-*"/"  ".
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
func ECVAward1(stateVotes int, pctBiden, pctTrump float64) (string, int, int, int) {
	glob := global.GetGlobalRef()
	pctOther := CalcOther(pctBiden, pctTrump)

	pctBiden += pctOther * pctBiden / 100.0
	pctTrump += pctOther * pctTrump / 100.0
	diff := math.Abs(pctBiden - pctTrump)

	if diff < glob.TossupThreshold {
		return "TOSSUP", 0, 0, stateVotes
	} else {
		if pctBiden > pctTrump {
			return "Biden", stateVotes, 0, 0
		} else {
			return "Trump", 0, stateVotes, 0
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
func ECVAward2(stateVotes int, pctBiden, pctTrump float64) (string, int, int, int, string) {
	glob := global.GetGlobalRef()
	diff := math.Abs(pctBiden - pctTrump)
	pctOther := CalcOther(pctBiden, pctTrump)
	if diff < glob.TossupThreshold {
		return "TOSSUP", 0, 0, stateVotes, getFactorString(pctOther > diff)
	}
	if pctBiden > pctTrump {
		return "Biden", stateVotes, 0, 0, getFactorString(pctOther > diff)
	}
	return "Trump", 0, stateVotes, 0, getFactorString(pctOther > diff)
}

/*
ECVAward3 - ECV Award Algorithm 3.

	Calculate the Other percentage = 100 - the sum of the candidate percentages.
	Calculate the difference = absolute value of the difference between the candidates.

	If the "Other" percentage exceeds the difference, then this state is a tossup.
	If the difference is below the tossup threshold, then this state is a tossup.
*/
func ECVAward3(stateVotes int, pctBiden, pctTrump float64) (string, int, int, int, string) {
	glob := global.GetGlobalRef()
	diff := math.Abs(pctBiden - pctTrump)
	pctOther := CalcOther(pctBiden, pctTrump)
	if pctOther > diff {
		return "TOSSUP", 0, 0, stateVotes, getFactorString(true)
	}
	if diff < glob.TossupThreshold {
		return "TOSSUP", 0, 0, stateVotes, " "
	}
	if pctBiden > pctTrump {
		return "Biden", stateVotes, 0, 0, " "
	}
	return "Trump", 0, stateVotes, 0, " "
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
