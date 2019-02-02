package main

import (
	"bufio"
	"fmt"
	"github.com/araddon/dateparse"
	"os"
	"regexp"
	"sync"
	"time"
)

const (
	SEEK_STEP int64 = -16384
)

func ReadFile(wg *sync.WaitGroup, log_file string, timestampRegex *regexp.Regexp, timeLeftBorder, timeRightBorder time.Time, junkLines int64) {
	defer wg.Done()
	fileHandler, err := os.Open(log_file)
	if err != nil {
		return
	}
	fmt.Println("parse", log_file, "from", timeLeftBorder, "to", timeRightBorder)
	leftBorder, rightBorder := SearchTimeBorders(fileHandler, timestampRegex, timeLeftBorder, timeRightBorder, junkLines)
	fmt.Println("read", log_file, "from", leftBorder, "to", rightBorder)
	_, err = fileHandler.Seek(leftBorder, 0)
	if err != nil {
		return
	}
	var line []byte
	reader := bufio.NewReader(fileHandler)
	line, _, err = reader.ReadLine()
	fmt.Println(string(line))
}

func parseTime(reader *bufio.Reader, timestampRegex *regexp.Regexp) ([]byte, *time.Time) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return line, nil
	}
	timeMatch := timestampRegex.Find(line)
	if timeMatch == nil {
		return line, nil
	}
	parsedTime, err := dateparse.ParseLocal(string(timeMatch))
	if err != nil {
		return line, nil
	}
	return line, &parsedTime
}

func SearchTimeBorders(fileHandler *os.File, timestampRegex *regexp.Regexp, timeLeftBorder, timeRightBorder time.Time, junkLines int64) (int64, int64) {
	// fileStat, err := fileHandler.Stat()
	// if err != nil {
	// 	return 0, 0
	// }

	curPos, err := fileHandler.Seek(SEEK_STEP, 2)
	if err != nil {
		curPos = 0
	}
	var line []byte
	var nextLine []byte
	var curPosTime *time.Time
	var nextLineTime *time.Time
	// Find minimum left border
	for curPos > 0 {
		reader := bufio.NewReader(fileHandler)
		for skipLines := int64(0); skipLines < junkLines; skipLines++ {
			line, curPosTime = parseTime(reader, timestampRegex)
			if curPosTime != nil {
				break
			}
			curPos += int64(len(line) + 1)
		}
		fmt.Println("cur", curPosTime)
		if curPosTime.Before(timeLeftBorder) {
			nextLine, nextLineTime = parseTime(reader, timestampRegex)
			fmt.Println("next", nextLineTime)
			if !nextLineTime.Before(timeLeftBorder) {
				curPos += int64(len(nextLine) + 1)
				break
			}

		}
		curPos, err = fileHandler.Seek(SEEK_STEP, 1)
		if err != nil {
			curPos = 0
		}
	}
	// Skip lines, while not equal or after
	// _, err = fileHandler.Seek(curPos, 0)
	// for curPos < fileStat.Size() {
	// 	reader := bufio.NewReader(fileHandler)
	// 	for skipLines := int64(0); skipLines < junkLines; skipLines++ {
	// 		line, _, err = reader.ReadLine()
	// 		if err == nil {
	// 			timeMatch := timestampRegex.Find(line)
	// 			if timeMatch != nil {
	// 				curPosTime, err = dateparse.ParseLocal(string(timeMatch))
	// 				if err == nil {
	// 					break
	// 				}
	// 			}
	// 		}
	// 		curPos += int64(len(line) + 1)
	// 	}
	// 	if !curPosTime.Before(timeLeftBorder) {
	// 		break
	// 	}
	// 	curPos += int64(len(line) + 1)
	// 	_, err = fileHandler.Seek(curPos, 1)
	// }
	return curPos, 0
}
