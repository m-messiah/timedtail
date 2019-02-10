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
	CHUNK_SIZE int64 = 128 * 1024
)

func searchFilePart(log_file string, timestampRegex *regexp.Regexp, timeBorders TimeBorders, junkLines int64, partsChannel chan FilePart) {
	fileHandler, err := os.Open(log_file)
	if err != nil {
		partsChannel <- FilePart{nil, 0, 0}
		return
	}
	from := searchOffset(fileHandler, timestampRegex, timeBorders.from, junkLines)
	to := searchOffset(fileHandler, timestampRegex, timeBorders.to, junkLines)
	partsChannel <- FilePart{fileHandler, from, to}
}

func readFile(filePart FilePart) {
	_, err := filePart.fileHandler.Seek(filePart.from, 0)
	if err != nil {
		return
	}
	chunk := make([]byte, CHUNK_SIZE)
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

func searchOffset(fileHandler *os.File, timestampRegex *regexp.Regexp, timeBorder time.Time, junkLines int64) int64 {
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
