package main

import (
	"bufio"
	"github.com/araddon/dateparse"
	"io"
	"os"
	"regexp"
	"time"
)

const (
	chunkSize int64 = 128 * 1024
)

func searchFilePart(logFile string, timeParams timeParams, partsChannel chan filePart) {
	fileHandler, err := os.Open(logFile)
	if err != nil {
		partsChannel <- filePart{nil, 0, 0}
		return
	}
	from := searchOffset(fileHandler, timeParams, timeParams.borders.from)
	to := searchOffset(fileHandler, timeParams, timeParams.borders.to)
	partsChannel <- filePart{fileHandler, from, to}
}

func readFile(filePart filePart) {
	_, err := filePart.fileHandler.Seek(filePart.from, 0)
	if err != nil {
		return
	}
	chunk := make([]byte, chunkSize)
	remainBytes := filePart.to - filePart.from
	for remainBytes > 0 {
		portion, err := filePart.fileHandler.Read(chunk)
		if err != nil {
			break
		}
		readBytes := int64(portion)
		if readBytes > remainBytes {
			readBytes = remainBytes
		}
		os.Stdout.Write(chunk[:readBytes])
		remainBytes -= readBytes
	}
}

func parseTime(reader *bufio.Reader, timestampRegex *regexp.Regexp) (*time.Time, int64) {
	line, err := reader.ReadBytes('\n')
	lineLen := int64(len(line))
	if err != nil && err != io.EOF {
		return nil, lineLen
	}
	timeMatch := timestampRegex.Find(line)
	if timeMatch == nil {
		return nil, lineLen
	}
	parsedTime, err := dateparse.ParseLocal(string(timeMatch))
	if err != nil {
		return nil, lineLen
	}
	return &parsedTime, lineLen
}

func parseLine(fileHandler *os.File, timeParams timeParams, seek int64) (int64, *time.Time, int64) {
	curPos, err := fileHandler.Seek(seek, 0)
	if err != nil {
		return seek, nil, seek
	}
	var lineLen int64
	var curPosTime *time.Time
	reader := bufio.NewReader(fileHandler)
	for skipLines := int64(0); skipLines < timeParams.junkLines; skipLines++ {
		curPosTime, lineLen = parseTime(reader, timeParams.regex)
		if curPosTime != nil {
			break
		}
		curPos += lineLen
	}
	return curPos, curPosTime, curPos + lineLen
}

func searchOffset(fileHandler *os.File, timeParams timeParams, timeBorder time.Time) int64 {
	l := int64(0)
	fileStat, err := fileHandler.Stat()
	if err != nil {
		return 0
	}
	fileSize := fileStat.Size()
	r := fileSize
	for l < r && l >= 0 && r <= fileSize {
		curPos, curPosTime, lineEnd := parseLine(fileHandler, timeParams, (l+r)/2)
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
			if lineEnd <= fileSize {
				l = lineEnd
			} else {
				l = fileSize
			}
		}
	}
	return l
}
