package main

import (
	"bufio"
	"fmt"
	"github.com/araddon/dateparse"
	"io"
	"os"
	"regexp"
	"sync"
	"time"
)

func readFile(wg *sync.WaitGroup, log_file string, timestampRegex *regexp.Regexp, timeLeftBorder, timeRightBorder time.Time, junkLines int64) {
	defer wg.Done()
	fileHandler, err := os.Open(log_file)
	if err != nil {
		return
	}
	fmt.Println("parse", log_file, "from", timeLeftBorder, "to", timeRightBorder)
	leftBorder := searchBorder(fileHandler, timestampRegex, timeLeftBorder, junkLines)
	rightBorder := searchBorder(fileHandler, timestampRegex, timeRightBorder, junkLines)
	fmt.Println("read", log_file, "from", leftBorder, "to", rightBorder)
	_, err = fileHandler.Seek(leftBorder, 0)
	if err != nil {
		return
	}
	var line []byte
	reader := bufio.NewReader(fileHandler)
	line, err = reader.ReadBytes('\n')
	fmt.Print(string(line))
}

func parseTime(reader *bufio.Reader, timestampRegex *regexp.Regexp) (*time.Time, int64) {
	line, err := reader.ReadBytes('\n')
	line_len := int64(len(line))
	if err != nil && err != io.EOF {
		return nil, line_len
	}
	timeMatch := timestampRegex.Find(line)
	if timeMatch == nil {
		return nil, line_len
	}
	parsedTime, err := dateparse.ParseLocal(string(timeMatch))
	if err != nil {
		return nil, line_len
	}
	return &parsedTime, line_len
}

func parseLine(fileHandler *os.File, timestampRegex *regexp.Regexp, junkLines int64, seek int64) (int64, *time.Time, int64) {
	curPos, err := fileHandler.Seek(seek, 0)
	if err != nil {
		return seek, nil, seek
	}
	var line_len int64
	var curPosTime *time.Time
	reader := bufio.NewReader(fileHandler)
	for skipLines := int64(0); skipLines < junkLines; skipLines++ {
		curPosTime, line_len = parseTime(reader, timestampRegex)
		if curPosTime != nil {
			break
		}
		curPos += line_len
	}
	return curPos, curPosTime, curPos + line_len
}

func searchBorder(fileHandler *os.File, timestampRegex *regexp.Regexp, timeBorder time.Time, junkLines int64) int64 {
	l := int64(0)
	fileStat, err := fileHandler.Stat()
	if err != nil {
		return 0
	}
	fileSize := fileStat.Size()
	r := fileSize
	for l < r && l >= 0 && r <= fileSize {
		curPos, curPosTime, line_end := parseLine(fileHandler, timestampRegex, junkLines, (l+r)/2)
		if r == curPos {
			if l == 0 {
				return l
			}
			return curPos
		}
		if curPosTime != nil && curPosTime.After(timeBorder) {
			if curPos <= fileSize {
				r = curPos
			} else {
				r = fileSize
			}
		} else {
			if line_end <= fileSize {
				l = line_end
			} else {
				l = fileSize
			}
		}
	}
	return l
}
