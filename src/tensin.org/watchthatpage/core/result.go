package core

import (
	"bytes"
	"time"

	"github.com/fatih/color"
)

// Result objects, containing various fields related to an analysis execution
type Result struct {
	URL                    string
	Title                  string
	Differences            string
	FirstExecution         bool
	HasDifferences         bool
	HasError               bool
	Error                  string
	CacheFileName          string
	ScreenshotFileName     string
	ScreenshotFullFileName string
	AnalysisExecutionTime  time.Duration
}

func (result *Result) dump(colored bool) string {
	var buf bytes.Buffer

	if result.FirstExecution {
		buf.WriteString("[")
		if colored {
			yellow := color.New(color.FgWhite, color.BgYellow).Add(color.Bold).SprintFunc()
			buf.WriteString(yellow("FIRST"))
		} else {
			buf.WriteString("FIRST")
		}
		buf.WriteString("] ")
	} else if result.HasError {
		buf.WriteString("[")
		if colored {
			red := color.New(color.FgWhite, color.BgRed).Add(color.Bold).SprintFunc()
			buf.WriteString(red("ERROR"))
		} else {
			buf.WriteString("ERROR")
		}
		buf.WriteString("] ")
	} else if result.HasDifferences {
		buf.WriteString("[")
		if colored {
			blue := color.New(color.FgWhite, color.BgBlue).Add(color.Bold).SprintFunc()
			buf.WriteString(blue("DIFF"))
		} else {
			buf.WriteString("DIFF")
		}
		buf.WriteString("]  ")
	} else {
		buf.WriteString("[")
		if colored {
			green := color.New(color.FgWhite, color.BgGreen).Add(color.Bold).SprintFunc()
			buf.WriteString(green("SAME"))
		} else {
			buf.WriteString("SAME")
		}
		buf.WriteString("]  ")
	}

	buf.WriteString("URL [")
	buf.WriteString(result.URL)
	buf.WriteString("], analysis took [")
	buf.WriteString(result.AnalysisExecutionTime.String())
	buf.WriteString("], cached content [")
	buf.WriteString(result.CacheFileName)
	buf.WriteString("]")
	return buf.String()
}

func (result *Result) String() string {
	return result.dump(false)
}

// ColoredString dumps the result as a string with a few colored entries
func (result *Result) ColoredString() string {
	return result.dump(true)
}
