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
	counterScan := 0
	counterPrint := 0

	log.Printf("State report: %s\n", state)
	rows := sqlQuery(sqlText)

	// For each row, process it...
	fmt.Printf("%-8s    %-4s  %-4s  %-4s %-s\n", "EndPoll", "Dem", "Gop", "Other", "Pollster")
	for rows.Next() {
		counterScan += 1
		err := rows.Scan(&query.state, &query.endDate, &query.pctDem, &query.pctGop, &query.pollster)
		if err != nil {
			log.Fatalf("ReportSC: rows.Scan failed, row count: %d, reason: %s\n", counterScan, err.Error())
		}
		tm, err := YYYY_MM_DDtoTime(query.endDate)
		if err != nil {
			log.Fatalf("ReportEC: Cannot parse start date: %s, reason: %s\n\n", query.endDate, err.Error())
		}
		if tm.Before(glob.DateThreshold) {
			continue
		}
		counterPrint++
		other := CalcOther(query.pctDem, query.pctGop)
		fmt.Printf("%-8s  %4.1f  %4.1f  %4.1f  %-s\n", query.endDate, query.pctDem, query.pctGop, other, query.pollster)
		if counterScan >= glob.PollHistoryLimit {
			break
		}
	}
	if counterPrint < 1 {
		fmt.Println("no data")
	}
}

func ReportEC() {
	glob := global.GetGlobalRef()
	var stateECV ECVote
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
	prtDivider := "-----------------------------------------------------------"
	fmt.Println("\nSt  EV  Last Poll   Dem       Gop       Other       Leading")
	fmt.Println(prtDivider)
	for _, stateECV = range stateECVTable {
		if glob.FlagBattleground {
			if !searchSlice(glob.Battleground, stateECV.state) {
				continue
			}
		}
		// For the given state, query from the most recent to the least recent polling.
		sqlText := fmt.Sprintf("SELECT end_date, pct_dem, pct_gop FROM history WHERE state = '%s' ORDER BY end_date DESC",
			stateECV.state)
		rows := sqlQuery(sqlText)

		counterScan := 0
		counterPrint := 0
		var query dbparams
		aveDemPct := 0.0
		aveGopPct := 0.0
		endDate := ""
		for rows.Next() {
			counterScan += 1
			err := rows.Scan(&query.endDate, &query.pctDem, &query.pctGop)
			if err != nil {
				log.Fatalf("ReportEC: rows.Scan failed, row count: %d, reason: %s\n", counterScan, err.Error())
			}
			tm, err := YYYY_MM_DDtoTime(query.endDate)
			if err != nil {
				log.Fatalf("ReportEC: Cannot parse start date: %s, reason: %s\n\n", query.endDate, err.Error())
			}
			if tm.Before(glob.DateThreshold) {
				continue
			}
			counterPrint += 1
			if counterScan == 1 {
				endDate = query.endDate
			}
			aveDemPct += query.pctDem
			arrayDemPct = append(arrayDemPct, query.pctDem)
			aveGopPct += query.pctGop
			arrayGopPct = append(arrayGopPct, query.pctGop)
			arrayOtherPct = append(arrayOtherPct, CalcOther(query.pctDem, query.pctGop))
			if counterScan >= glob.PollHistoryLimit {
				break
			}
		}

		if counterPrint < 1 {
			fmt.Printf("%-2s  no data\n", stateECV.state)
			continue
		}

		// Averages for this state.
		aveDemPct /= float64(counterScan)
		aveGopPct /= float64(counterScan)
		aveOtherPct := CalcOther(aveDemPct, aveGopPct)
		leader := ""
		var increDem, increGop, increTossup int
		otherFactor := ""
		switch glob.ECVAlgorithm {
		case 1:
			leader, increDem, increGop, increTossup = ECVAward1(stateECV.votes, aveDemPct, aveGopPct)
		case 2:
			leader, increDem, increGop, increTossup, otherFactor = ECVAward2(stateECV.votes, aveDemPct, aveGopPct)
		case 3:
			leader, increDem, increGop, increTossup, otherFactor = ECVAward3(stateECV.votes, aveDemPct, aveGopPct)
		default:
			log.Fatalf("ReportEC: global.ECVAlgoithm %d is not supported\n", glob.ECVAlgorithm)
		}

		totalDemECV += increDem
		totalGopECV += increGop
		totalTossupECV += increTossup
		switch leader {
		case "Dem":
			counterDemStates++
			listDemStates += " " + stateECV.state
		case "Gop":
			counterGopStates++
			listGopStates += " " + stateECV.state
		default:
			counterTossupStates++
			listTossupStates += " " + stateECV.state
		}

		// Show results for current state.
		demTrend := CalcTrend(arrayDemPct)
		gopTrend := CalcTrend(arrayGopPct)
		otherTrend := CalcTrend(arrayOtherPct)
		fmt.Printf("%-2s  %2d  %-8s  %4.1f  %s  %4.1f  %s  %4.1f  %2s%2s  %-s\n",
			stateECV.state, stateECV.votes, endDate, aveDemPct, demTrend,
			aveGopPct, gopTrend, aveOtherPct, otherTrend, otherFactor, leader)

	}
	// Totals.
	fmt.Println(prtDivider)
	if glob.ECVAlgorithm != 1 {
		fmt.Println("** The Other percentage exceeds the difference between Dem and Gop.")
	}
	fmt.Printf("Dem    EV: %3d, states: (%2d)%s\n", totalDemECV, counterDemStates, listDemStates)
	fmt.Printf("Gop    EV: %3d, states: (%2d)%s\n", totalGopECV, counterGopStates, listGopStates)
	fmt.Printf("Tossup EV: %3d, states: (%2d)%s\n", totalTossupECV, counterTossupStates, listTossupStates)
}
