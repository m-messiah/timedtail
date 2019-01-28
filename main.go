package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

func assertFlagIsPositiveInt(value int64) {
	if value < 0 {
		flag.Usage()
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <log files>...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	timeNow := time.Now()
	deltaTime := flag.Duration("n", time.Duration(300)*time.Second, "Delta time count from now() to look at from the end of the <log file>")
	fromTime := flag.Int64("b", 0, "From which datetime start to operate seconds of the <log file>")
	timestampType := flag.String("t", "common", "Timestamp type")
	timestampCustomRegex := flag.String("r", "", "Regexp to pick timestamp from string ($1 must select timestamp)")
	junkLines := flag.Int64("j", 500, "Max number of junk lines to read")
	flag.Parse()
	log_files := flag.Args()
	if len(log_files) == 0 {
		flag.Usage()
	}
	assertFlagIsPositiveInt(deltaTime.Nanoseconds())
	assertFlagIsPositiveInt(*fromTime)
	assertFlagIsPositiveInt(*junkLines)
	if *timestampType != "common" && *timestampCustomRegex != "" {
		flag.Usage()
	}
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
	timestampRegex := regexp.MustCompile(timestampRegexString)

	for _, log_file := range log_files {
		fmt.Print(log_file, ",")
	}
	fmt.Println("")
	fmt.Println(timeNow)
	fmt.Println(*deltaTime)
	fmt.Println(*fromTime)
	fmt.Println(*timestampType)
	fmt.Println(*timestampCustomRegex)
	fmt.Println(*junkLines)
	fmt.Println(timestampRegex)

}
