package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

type filePart struct {
	fileHandler *os.File
	from        int64
	to          int64
}

type timeBorders struct {
	from time.Time
	to   time.Time
}

type timeParams struct {
	regex     *regexp.Regexp
	borders   timeBorders
	junkLines int64
}

func assertFlagIsPositiveInt(value int64) {
	if value < 0 {
		flag.Usage()
	}
}

func getTimeStampRegex(timestampCustomRegex *string, timestampType *string) *regexp.Regexp {
	var timestampRegexString string
	if *timestampCustomRegex != "" {
		timestampRegexString = *timestampCustomRegex
	} else {
		lookupRegexString, ok := TimestampTypes[*timestampType]
		if ok {
			timestampRegexString = lookupRegexString
		} else {
			fmt.Fprintf(os.Stderr, "Unkown type %s, possible values are: ", *timestampType)
			for k := range TimestampTypes {
				fmt.Fprint(os.Stderr, k, " ")
			}
			fmt.Fprintln(os.Stderr, "")
			os.Exit(1)
		}
	}
	return regexp.MustCompile(timestampRegexString)
}

func getTimeBorders(fromTime, deltaSeconds *int64) timeBorders {
	var timeRightBorder time.Time
	if *fromTime != 0 {
		timeRightBorder = time.Unix(*fromTime, 0)
	} else {
		timeRightBorder = time.Now().Round(time.Second)
	}
	timeLeftBorder := timeRightBorder.Add(-time.Duration(*deltaSeconds) * time.Second)
	return timeBorders{timeLeftBorder, timeRightBorder}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <log files>...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	deltaSeconds := flag.Int64("n", 300, "Delta time count from now() to look at from the end of the <log file>")
	fromTime := flag.Int64("b", 0, "From which datetime start to operate seconds of the <log file>")
	timestampType := flag.String("t", "common", "Timestamp type")
	timestampCustomRegex := flag.String("r", "", "Regexp to pick timestamp from string ($1 must select timestamp)")
	junkLines := flag.Int64("j", 500, "Max number of junk lines to read")
	utcLogTime := flag.Bool("utc", false, "Enable UTC timezone for time parse")
	flag.Parse()
	logFiles := flag.Args()
	if len(logFiles) == 0 {
		flag.Usage()
	}
	assertFlagIsPositiveInt(*deltaSeconds)
	assertFlagIsPositiveInt(*fromTime)
	assertFlagIsPositiveInt(*junkLines)
	if *timestampType != "common" && *timestampCustomRegex != "" {
		flag.Usage()
	}
	timestampRegex := getTimeStampRegex(timestampCustomRegex, timestampType)
	timeBordersVal := getTimeBorders(fromTime, deltaSeconds)
	if *utcLogTime {
		utcLoc, _ := time.LoadLocation("UTC")
		time.Local = utcLoc
	}
	partsChannel := make(chan filePart)
	for _, logFile := range logFiles {
		go searchFilePart(logFile, timeParams{timestampRegex, timeBordersVal, *junkLines}, partsChannel)
	}

	for i := 0; i < len(logFiles); i++ {
		readFile(<-partsChannel)
	}
}
