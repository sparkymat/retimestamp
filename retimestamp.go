package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

func timestampRegex(timestampFormat string) *regexp.Regexp {
	timestampRegexFormat := strings.Replace(timestampFormat, "yyyy", "[1,2][0-9]{3}", -1)
	timestampRegexFormat = strings.Replace(timestampRegexFormat, "mm", "[0,1][0-9]", -1)
	timestampRegexFormat = strings.Replace(timestampRegexFormat, "dd", "[0-3][0-9]", -1)
	timestampRegexFormat = strings.Replace(timestampRegexFormat, "HH", "[0-2][0-9]", -1)
	timestampRegexFormat = strings.Replace(timestampRegexFormat, "MM", "[0-5][0-9]", -1)
	timestampRegexFormat = strings.Replace(timestampRegexFormat, "SS", "[0-5][0-9]", -1)
	return regexp.MustCompile(timestampRegexFormat)
}

func RetimestampString(inputString string, timestampFormat string) (string, string, string) {
	timestampDisplayFormat := strings.Replace(timestampFormat, "yyyy", "2006", -1)
	timestampDisplayFormat = strings.Replace(timestampDisplayFormat, "mm", "01", -1)
	timestampDisplayFormat = strings.Replace(timestampDisplayFormat, "dd", "02", -1)
	timestampDisplayFormat = strings.Replace(timestampDisplayFormat, "HH", "15", -1)
	timestampDisplayFormat = strings.Replace(timestampDisplayFormat, "MM", "04", -1)
	timestampDisplayFormat = strings.Replace(timestampDisplayFormat, "SS", "05", -1)
	newTimestamp := time.Now().UTC().Format(timestampDisplayFormat)

	matchedTimestamp := timestampRegex(timestampFormat).FindStringSubmatch(inputString)[0]
	updatedFilenameBytes := timestampRegex(timestampFormat).ReplaceAll([]byte(inputString), []byte(newTimestamp))
	return string(updatedFilenameBytes), matchedTimestamp, newTimestamp
}

func main() {
	var timestampFormat string

	flag.StringVar(&timestampFormat, "format", "yyyymmddHHMMSS", "Format of the timestamp")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("[Error] No file provided")
		os.Exit(1)
	}

	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		fmt.Println("[Error] File not found")
		os.Exit(1)
	}

	if !timestampRegex(timestampFormat).Match([]byte(args[0])) {
		fmt.Println("[Error] No timestamp found in file. Did you specify the correct format?")
		os.Exit(1)
	}

	updatedFilename, oldTimestamp, newTimestamp := RetimestampString(args[0], timestampFormat)

	highlight := color.New(color.Bold, color.FgHiWhite).Add(color.Underline).SprintfFunc()
	fmt.Printf("Are you sure you want to rename\n\t%v\nto\n\t%v\n? (yes/no) :", strings.Replace(args[0], oldTimestamp, highlight(oldTimestamp), -1), strings.Replace(updatedFilename, newTimestamp, highlight(newTimestamp), -1))

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (strings.TrimSpace(strings.ToLower(input)) != "yes" && strings.TrimSpace(strings.ToLower(input)) != "y") {
		fmt.Printf("[Error] User abort")
		os.Exit(1)
	}

	os.Rename(args[0], updatedFilename)

	fmt.Printf("Successfully renamed")
}
