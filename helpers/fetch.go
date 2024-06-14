package helpers

import (
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func moveTempCSVToCurrentCSV(pathTemp, pathCurrent string) {
	fileBytes, err := os.ReadFile(pathTemp)
	if err != nil {
		log.Fatalf("moveTempCSVToCurrentCSV: ReadFile(%s) failed, reason: %s\n", pathTemp, err.Error())
	}

	err = os.WriteFile(pathCurrent, fileBytes, 0644)
	if err != nil {
		log.Fatalf("moveTempCSVToCurrentCSV: WriteFile(%s) failed, reason: %s\n", pathCurrent, err.Error())
	}

	err = os.Remove(pathTemp)
	if err != nil {
		log.Fatalf("moveTempCSVToCurrentCSV: os.Remove(%s) failed, reason: %s\n", pathTemp, err.Error())
	}
}

func Fetch(dirCsv, fileName, url, dirTemp string) {

	// Form the full path of the CSV file in the final directory and the temporary directory.
	pathCurrentCSV := filepath.Join(dirCsv, fileName)
	pathTempCSV := filepath.Join(dirTemp, fileName)
	log.Printf("Fetch: Will retrieve %s from %s\n", pathCurrentCSV, url)

	// Get the data.
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Fetch: http.Get(%s) failed, reason: %s\n", url, err.Error())
	}
	defer resp.Body.Close()

	// Create the temp CSV file.
	outf, err := os.Create(pathTempCSV)
	if err != nil {
		log.Fatalf("Fetch: os.Create(%s) failed, reason: %s\n", pathTempCSV, err.Error())
	}

	// Write the response body to the temp CSV file.
	_, err = io.Copy(outf, resp.Body)
	if err != nil {
		log.Fatalf("Fetch: io.Copy(%s) failed, reason: %s\n", pathTempCSV, err.Error())
	}

	// Close the temp CSV file.
	err = outf.Close()
	if err != nil {
		log.Fatalf("Fetch: os.Close(%s) failed, reason: %s\n", pathTempCSV, err.Error())
	}

	// Compute checksum for current CSV file.
	fileBytes, err := os.ReadFile(pathCurrentCSV)
	if err != nil {
		// An error occurred.
		// Assume that there is no current CSV file.
		// Copy the temp CSV file to the current CSV file.
		log.Println("Fetch: No previous poll data.")
		moveTempCSVToCurrentCSV(pathTempCSV, pathCurrentCSV)
		log.Println("Fetch: End")
		return
	}
	cksumCurrent := crc32.ChecksumIEEE(fileBytes)

	// Compute checksum for temp file.
	fileBytes, err = os.ReadFile(pathTempCSV)
	if err != nil {
		log.Fatalf("Fetch: ReadFile(%s) failed, reason: %s\n", pathTempCSV, err.Error())
	}
	cksumTemp := crc32.ChecksumIEEE(fileBytes)

	// Any changes from last time?
	if cksumCurrent != cksumTemp {
		log.Println("Fetch: Internet poll data has changed.")
		moveTempCSVToCurrentCSV(pathTempCSV, pathCurrentCSV)
		log.Println("Fetch: End")
		return
	}

	// ByeBye.
	log.Println("Fetch: Internet poll data has not changed. Nothing to do.")
	log.Println("Fetch: End")
}
