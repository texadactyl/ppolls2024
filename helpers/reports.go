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
	sqlText := fmt.Sprintf("SELECT state, end_date, pct_biden, pct_trump, pollster FROM history WHERE state = '%s' ORDER BY end_date DESC", state)

	// Get all the selected history table rows.
	counter := 0
	log.Printf("State report: %s\n", state)
	rows := sqlQuery(sqlText)

	// For each row, process it...
	fmt.Printf("%-8s    %-4s %-4s %-4s %-s\n", "EndPoll", "Biden", "Trump", "Other", "Pollster")
	for rows.Next() {
		counter += 1
		err := rows.Scan(&query.state, &query.endDate, &query.pctBiden, &query.pctTrump, &query.pollster)
		if err != nil {
			log.Fatalf("ReportSC: rows.Scan failed, row count: %d, reason: %s\n", counter, err.Error())
		}
		other := CalcOther(query.pctBiden, query.pctTrump)
		fmt.Printf("%-8s  %4.1f  %4.1f  %4.1f  %-s\n", query.endDate, query.pctBiden, query.pctTrump, other, query.pollster)
		if counter >= glob.PollHistoryLimit {
			break
		}
	}
}

func ReportEC() {
	glob := global.GetGlobalRef()
	var stateECV ECVote
	totalBidenECV := 0
	totalTrumpECV := 0
	totalTossupECV := 0
	counterBidenStates := 0
	counterTrumpStates := 0
	counterTossupStates := 0
	listBidenStates := ""
	listTrumpStates := ""
	listTossupStates := ""
	var arrayBidenPct []float64
	var arrayTrumpPct []float64
	var arrayOtherPct []float64
	prtDivider := "-----------------------------------------------------------"
	fmt.Println("\nSt  EV  Last Poll   Biden     Trump      Other      Leading")
	fmt.Println(prtDivider)
	for _, stateECV = range stateECVTable {
		// For the given state, query from the most recent to the least recent polling.
		sqlText := fmt.Sprintf("SELECT end_date, pct_biden, pct_trump FROM history WHERE state = '%s' ORDER BY end_date DESC",
			stateECV.state)
		rows := sqlQuery(sqlText)

		counter := 0
		var query dbparams
		aveBidenPct := 0.0
		aveTrumpPct := 0.0
		endDate := ""
		for rows.Next() {
			counter += 1
			err := rows.Scan(&query.endDate, &query.pctBiden, &query.pctTrump)
			if err != nil {
				log.Fatalf("ReportEC: rows.Scan failed, row count: %d, reason: %s\n", counter, err.Error())
			}
			if counter == 1 {
				endDate = query.endDate
			}
			aveBidenPct += query.pctBiden
			arrayBidenPct = append(arrayBidenPct, query.pctBiden)
			aveTrumpPct += query.pctTrump
			arrayTrumpPct = append(arrayTrumpPct, query.pctTrump)
			arrayOtherPct = append(arrayOtherPct, CalcOther(query.pctBiden, query.pctTrump))
			if counter >= glob.PollHistoryLimit {
				break
			}
		}

		// Averages for this state.
		aveBidenPct /= float64(counter)
		aveTrumpPct /= float64(counter)
		aveOtherPct := CalcOther(aveBidenPct, aveTrumpPct)
		leader := ""
		var increBiden, increTrump, increTossup int
		otherFactor := " "
		switch glob.ECVAlgorithm {
		case 1:
			leader, increBiden, increTrump, increTossup, otherFactor = ECVAward1(stateECV.votes, aveBidenPct, aveTrumpPct)
		case 2:
			leader, increBiden, increTrump, increTossup = ECVAward2(stateECV.votes, aveBidenPct, aveTrumpPct)
		default:
			log.Fatalf("ReportEC: global.ECVAlgoithm %d is not supported\n", glob.ECVAlgorithm)
		}

		totalBidenECV += increBiden
		totalTrumpECV += increTrump
		totalTossupECV += increTossup
		switch leader {
		case "Biden":
			counterBidenStates++
			listBidenStates += " " + stateECV.state
		case "Trump":
			counterTrumpStates++
			listTrumpStates += " " + stateECV.state
		default:
			counterTossupStates++
			listTossupStates += " " + stateECV.state
		}

		// Show results for current state.
		bidenTrend := CalcTrend(arrayBidenPct)
		trumpTrend := CalcTrend(arrayTrumpPct)
		otherTrend := CalcTrend(arrayOtherPct)
		fmt.Printf("%-2s  %2d  %-8s  %4.1f  %s  %4.1f  %s  %4.1f  %2s%2s  %-s\n",
			stateECV.state, stateECV.votes, endDate, aveBidenPct, bidenTrend,
			aveTrumpPct, trumpTrend, aveOtherPct, otherTrend, otherFactor, leader)

	}
	// Totals.
	fmt.Println(prtDivider)
	fmt.Printf("Biden  EV: %3d, states: (%2d)%s\n", totalBidenECV, counterBidenStates, listBidenStates)
	fmt.Printf("Trump  EV: %3d, states: (%2d)%s\n", totalTrumpECV, counterTrumpStates, listTrumpStates)
	fmt.Printf("Tossup EV: %3d, states: (%2d)%s\n", totalTossupECV, counterTossupStates, listTossupStates)
}
