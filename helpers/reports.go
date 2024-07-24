package helpers

import (
	"fmt"
	"log"
	"ppolls2024/global"
)

func ReportSC(state string) {
	glob := global.GetGlobalRef()
	// Query record
	var query dbparams

	// For the given state, query from the most recent to the least recent polling.
	sqlText := fmt.Sprintf("SELECT state, end_date, pct_dem, pct_gop, pollster FROM history WHERE state = '%s' ORDER BY end_date DESC", state)

	// Get all the selected history table rows.
	counterRows := 0
	log.Printf("State report: %s\n", state)
	rows := sqlQuery(sqlText)
	fmt.Printf("%-8s    %-4s  %-4s  %-4s %-s\n", "EndPoll", "Dem", "Gop", "Other", "Pollster")

	// For each row, process it...
	for rows.Next() {
		counterRows += 1
		err := rows.Scan(&query.state, &query.endDate, &query.pctDem, &query.pctGop, &query.pollster)
		if err != nil {
			log.Fatalf("ReportSC: rows.Scan failed, row count: %d, reason: %s\n", counterRows, err.Error())
		}
		tm, err := YYYY_MM_DDtoTime(query.endDate)
		if err != nil {
			log.Fatalf("ReportEC: Cannot parse start date: %s, reason: %s\n\n", query.endDate, err.Error())
		}
		if tm.Before(glob.DateThreshold) {
			continue
		}
		other := CalcOther(query.pctDem, query.pctGop)
		fmt.Printf("%-8s  %4.1f  %4.1f  %4.1f  %-s\n", query.endDate, query.pctDem, query.pctGop, other, query.pollster)
		if counterRows >= glob.PollHistoryLimit {
			break
		}
	}
	if counterRows < 1 {
		fmt.Println("no data")
	}
}

func ReportEC() {
	glob := global.GetGlobalRef()
	var stateTableEntry global.StateTableEntry_t
	totalDemECV := 0
	totalGopECV := 0
	totalTossupECV := 0
	counterDemStates := 0
	counterGopStates := 0
	counterTossupStates := 0
	listDemStates := ""
	listGopStates := ""
	listTossupStates := ""
	var arrayDemPct []float64
	var arrayGopPct []float64
	var arrayOtherPct []float64
	prtDivider := "------------------------------------------------------------"
	fmt.Println("\nSt   EV  Last Poll   Dem       Gop       Other       Leading")
	fmt.Println(prtDivider)
	for _, stateTableEntry = range global.StateTable {
		if glob.FlagBattleground {
			if !searchSlice(glob.Battleground, stateTableEntry.Stcode) {
				continue
			}
		}
		// For the given state, query from the most recent to the least recent polling.
		sqlText := fmt.Sprintf("SELECT end_date, pct_dem, pct_gop FROM history WHERE state = '%s' ORDER BY end_date DESC",
			stateTableEntry.Stcode)
		rows := sqlQuery(sqlText)

		counterRows := 0
		var query dbparams
		aveDemPct := 0.0
		aveGopPct := 0.0
		aveOtherPct := 0.0
		endDate := ""
		for rows.Next() {
			err := rows.Scan(&query.endDate, &query.pctDem, &query.pctGop)
			if err != nil {
				log.Fatalf("ReportEC: rows.Scan failed, row count: %d, reason: %s\n", counterRows, err.Error())
			}
			tm, err := YYYY_MM_DDtoTime(query.endDate)
			if err != nil {
				log.Fatalf("ReportEC: Cannot parse start date: %s, reason: %s\n\n", query.endDate, err.Error())
			}
			if tm.Before(glob.DateThreshold) {
				continue
			}
			counterRows += 1

			// If first row, that is the end date.
			if counterRows == 1 {
				endDate = query.endDate
			}

			aveDemPct += query.pctDem
			arrayDemPct = append(arrayDemPct, query.pctDem)
			aveGopPct += query.pctGop
			arrayGopPct = append(arrayGopPct, query.pctGop)
			arrayOtherPct = append(arrayOtherPct, CalcOther(query.pctDem, query.pctGop))

			// Don't go over the poll history threshold.
			if counterRows >= glob.PollHistoryLimit {
				break
			}
		}

		// Got any data for this state?
		if counterRows < 1 { // NO DATA
			endDate = "no data   "
			var strongly bool
			strongly = searchSlice(glob.StronglyDem, stateTableEntry.Stcode)
			if strongly { // Strongly Democrat
				aveDemPct = 99.9
				aveGopPct = 0.0
				aveOtherPct = 0.0
			} else {
				strongly = searchSlice(glob.StronglyGop, stateTableEntry.Stcode)
				if strongly { // Strongly GOP
					aveDemPct = 0.0
					aveGopPct = 99.9
					aveOtherPct = 0.0
				} else { // Battleground
					aveDemPct = 0.0
					aveGopPct = 0.0
					aveOtherPct = 99.9
				}
			}
		} else { // We have data for this state.
			// Averages for this state.
			aveDemPct /= float64(counterRows)
			aveGopPct /= float64(counterRows)
			aveOtherPct = CalcOther(aveDemPct, aveGopPct)
		}

		// Compute leader and the increments.
		var leader string
		var increDem, increGop, increTossup int
		otherFactor := ""
		switch glob.ECVAlgorithm {
		case 1:
			leader, increDem, increGop, increTossup = ECVAward1(stateTableEntry.Votes, aveDemPct, aveGopPct)
		case 2:
			leader, increDem, increGop, increTossup, otherFactor = ECVAward2(stateTableEntry.Votes, aveDemPct, aveGopPct)
		case 3:
			leader, increDem, increGop, increTossup, otherFactor = ECVAward3(stateTableEntry.Votes, aveDemPct, aveGopPct)
		default:
			log.Fatalf("ReportEC: global.ECVAlgoithm %d is not supported\n", glob.ECVAlgorithm)
		}

		totalDemECV += increDem
		totalGopECV += increGop
		totalTossupECV += increTossup
		switch leader {
		case "Dem":
			counterDemStates++
			listDemStates += " " + stateTableEntry.Stcode
		case "Gop":
			counterGopStates++
			listGopStates += " " + stateTableEntry.Stcode
		default:
			counterTossupStates++
			listTossupStates += " " + stateTableEntry.Stcode
		}

		// Show results for current state.
		demTrend := CalcTrend(arrayDemPct)
		gopTrend := CalcTrend(arrayGopPct)
		otherTrend := CalcTrend(arrayOtherPct)
		fmt.Printf("%-2s  %3d  %-8s  %4.1f  %s  %4.1f  %s  %4.1f  %2s%2s  %-s\n",
			stateTableEntry.Stcode, stateTableEntry.Votes, endDate, aveDemPct, demTrend,
			aveGopPct, gopTrend, aveOtherPct, otherTrend, otherFactor, leader)

	} // for _, stateTableEntry = range global.StateTable

	// Totals.
	fmt.Println(prtDivider)
	if glob.ECVAlgorithm != 1 {
		fmt.Println("** The Other percentage exceeds the difference between Dem and Gop.")
	}
	fmt.Printf("Dem    EV: %3d, states: (%2d)%s\n", totalDemECV, counterDemStates, listDemStates)
	fmt.Printf("Gop    EV: %3d, states: (%2d)%s\n", totalGopECV, counterGopStates, listGopStates)
	fmt.Printf("Tossup EV: %3d, states: (%2d)%s\n", totalTossupECV, counterTossupStates, listTossupStates)
}
