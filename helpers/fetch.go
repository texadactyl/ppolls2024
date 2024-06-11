package helpers

import (
	"io"
	"net/http"
	"os"
)

func Fetch(filepath, url string) {
	// Hello.
	Logger("Fetch: Store %s from %s", filepath, url)

	// Get the data.
	resp, err := http.Get(url)
	if err != nil {
		Croak("*** ERROR, http.Get(%s) failed, reason: %s", url, err.Error())
	}
	defer resp.Body.Close()

	// Create the file.
	outf, err := os.Create(filepath)
	if err != nil {
		Croak("*** ERROR, os.Create(%s) failed, reason: %s", filepath, err.Error())
	}
	defer outf.Close()

	// Write the response body to file.
	_, err = io.Copy(outf, resp.Body)
	if err != nil {
		Croak("*** ERROR, io.Copy(%s) failed, reason: %s", url, err.Error())
	}

	// ByeBye.
	Logger("Fetch: End")
}
