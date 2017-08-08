//+build !nolog
//+build !js

package dlog

import (
	"bufio"
	"bytes"
	"fmt"

	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	byt    = bytes.NewBuffer(make([]byte, 0))
	writer *bufio.Writer
)

// dLog, the primary function of the package,
// prints out and writes to file a string
// containing the logged data separated by spaces,
// prepended with file and line information.
// It only includes logs which pass the current filters.
// Todo: use io.Multiwriter to simplify the writing to
// both logfiles and stdout
func dLog(console, override bool, in ...interface{}) {
	//(pc uintptr, file string, line int, ok bool)
	_, f, line, ok := runtime.Caller(2)
	if ok {
		f = truncateFileName(f)
		if !checkFilter(f, in) && !override {
			return
		}

		// Note on errors: these functions all return
		// errors, but they are always nil.
		byt.WriteRune('[')
		byt.WriteString(f)
		byt.WriteRune(':')
		byt.WriteString(strconv.Itoa(line))
		byt.WriteString("]  ")
		for _, elem := range in {
			byt.WriteString(fmt.Sprintf("%v ", elem))
		}
		byt.WriteRune('\n')

		if console {
			fmt.Print(byt.String())
		}

		if writer != nil {
			_, err := writer.WriteString(byt.String())
			if err != nil {
				// We can't log errors while we are in the error
				// logging function.
				fmt.Println("Logging error", err)
			}
			err = writer.Flush()
			if err != nil {
				fmt.Println("Logging error", err)
			}
		}

		byt.Reset()
	}
}

// FileWrite runs dLog, but JUST writes to file instead
// of also to stdout.
func FileWrite(in ...interface{}) {
	dLog(false, true, in...)
}

// CreateLogFile creates a file in the 'logs' directory
// of the starting point of this program to write logs to
func CreateLogFile() {
	file := "logs/dlog"
	file += time.Now().Format("_Jan_2_15-04-05_2006")
	file += ".txt"
	fHandle, err := os.Create(file)
	if err != nil {
		// We can't log an error that comes from
		// our error logging functions
		//panic(err)
		// But this is also not an error we want to panic on!
		fmt.Println("[oak]-------- No logs directory found. No logs will be written to file.")
		return
	}
	writer = bufio.NewWriter(fHandle)
}
