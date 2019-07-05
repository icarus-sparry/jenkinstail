package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type LogData struct {
	Data     []byte
	TextSize int64 // Current output position. Is it worth holding as an int?
	MoreData bool  // jenkins says there is potentially more data
}

// GetLogFromURL returns part of the log from the given start point
// updates the LogData to have returned data and updated start point
func GetLogFromURL(buildUrl string, logData *LogData) error {
	path := fmt.Sprintf("%s/logText/progressiveText?start=%d",
		buildUrl, logData.TextSize)
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//	fmt.Fprintln(os.Stderr, path)
	// fmt.Fprintln(os.Stderr, resp.Header)
	more := resp.Header["X-More-Data"]
	sizeHeaders := resp.Header["X-Text-Size"]

	// Does Jenkins say there is more data?
	logData.MoreData = more != nil && len(more) > 0 && more[0] == "true"

	// Where does Jenkins say that this data ends?
	// Badly named header in Jenkins
	// Is it worth converting to an int64, or just hold as string?
	if sizeHeaders != nil && len(sizeHeaders) > 0 {
		sh := sizeHeaders[0]
		if len(sh) > 0 {
			i, err := strconv.ParseInt(sh, 10, 64)
			if err == nil && i >= 0 {
				logData.TextSize = i
			}
		}
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid response %d", resp.StatusCode)
	}
	logData.Data, err = ioutil.ReadAll(resp.Body)
	return err
}

func main() {
	l := LogData{}
	for {
		err := GetLogFromURL(os.Args[1], &l)
		if err != nil {
			log.Fatal(err)
		}
		// Change CR/LF to LF.
		// TODO: Maybe make this faster by an inplace change
		// TODO: Makeby change CR/CR/LF to LF as well?
		s := string(bytes.Replace(l.Data, []byte{13, 10}, []byte{10}, -1))
		// Just send it to stdout.
		fmt.Print(s)
		if !l.MoreData {
			return
		}
		// If we didn't get any data last time, then throttle ourselves
		if len(s) == 0 {
			// fmt.Fprintln(os.Stderr, "Sleeping")
			time.Sleep(2000 * time.Millisecond)
		}
	}
}
