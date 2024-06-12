package helpers

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Fetch(filepath, url string) {
	// Hello.
	log.Printf("Fetch: Store %s from %s\n", filepath, url)

	// Get the data.
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("*** ERROR, http.Get(%s) failed, reason: %s\n", url, err.Error())
	}
	defer resp.Body.Close()

	// Create the file.
	outf, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("*** ERROR, os.Create(%s) failed, reason: %s\n", filepath, err.Error())
	}
	defer outf.Close()

	// Write the response body to file.
	_, err = io.Copy(outf, resp.Body)
	if err != nil {
		log.Fatalf("*** ERROR, io.Copy(%s) failed, reason: %s\n", url, err.Error())
	}

	// ByeBye.
	log.Println("Fetch: End")
}
