package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

type FilePart struct {
	fileHandler *os.File
	from        int64
	to          int64
}

type TimeBorders struct {
	from time.Time
	to   time.Time
}

type TimeParams struct {
	regex     *regexp.Regexp
	borders   TimeBorders
	junkLines int64
}

func AssertFlagIsPositiveInt(value int64) {
	if value < 0 {
		flag.Usage()
	}
}

func getTimeStampRegex(timestampCustomRegex *string, timestampType *string) *regexp.Regexp {
	var timestampRegexString string
	if *timestampCustomRegex != "" {
		timestampRegexString = *timestampCustomRegex
	} else {
		lookupRegexString, ok := TIMESTAMP_TYPES[*timestampType]
		if ok {
			timestampRegexString = lookupRegexString
		} else {
			fmt.Fprintf(os.Stderr, "Unkown type %s, possible values are: ", *timestampType)
			for k := range TIMESTAMP_TYPES {
				fmt.Fprint(os.Stderr, k, " ")
			}
			fmt.Fprintln(os.Stderr, "")
			os.Exit(1)
		}
	}
	return regexp.MustCompile(timestampRegexString)
}

func getTimeBorders(fromTime, deltaSeconds *int64) TimeBorders {
	var timeRightBorder time.Time
	if *fromTime != 0 {
		timeRightBorder = time.Unix(*fromTime, 0)
	} else {
		timeRightBorder = time.Now().Round(time.Second)
	}
	timeLeftBorder := timeRightBorder.Add(-time.Duration(*deltaSeconds) * time.Second)
	return TimeBorders{timeLeftBorder, timeRightBorder}
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
	flag.Parse()
	log_files := flag.Args()
	if len(log_files) == 0 {
		flag.Usage()
	}
	AssertFlagIsPositiveInt(*deltaSeconds)
	AssertFlagIsPositiveInt(*fromTime)
	AssertFlagIsPositiveInt(*junkLines)
	if *timestampType != "common" && *timestampCustomRegex != "" {
		flag.Usage()
	}
	timestampRegex := getTimeStampRegex(timestampCustomRegex, timestampType)
	timeBorders := getTimeBorders(fromTime, deltaSeconds)
	partsChannel := make(chan FilePart)
	for _, log_file := range log_files {
		go searchFilePart(log_file, TimeParams{timestampRegex, timeBorders, *junkLines}, partsChannel)
	}

	for i := 0; i < len(log_files); i++ {
		readFile(<-partsChannel)
	}
}
