package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <log files>...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	timeNow := time.Now()
	deltaTime := flag.Duration("n", time.Duration(300)*time.Second, "Delta time count from now() to look at from the end of the <log file>")
	fromTime := flag.Int("b", 0, "From which datetime start to operate seconds of the <log file>")
	timestampType := flag.String("t", "common", "Timestamp type")
	timestampRegex := flag.String("r", "", "Regexp to pick timestamp from string ($1 must select timestamp)")
	junkLines := flag.Int("j", 500, "Max number of junk lines to read")
	flag.Parse()
	log_files := flag.Args()
	if len(log_files) == 0 {
		flag.Usage()
	}
	for _, log_file := range log_files {
		fmt.Println(log_file)
	}
	fmt.Println(timeNow)
	fmt.Println(*deltaTime)
	fmt.Println(*fromTime)
	fmt.Println(*timestampType)
	fmt.Println(*timestampRegex)
	fmt.Println(*junkLines)
}
