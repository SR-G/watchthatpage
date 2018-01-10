package core

import (
	// "html"
	// "github.com/sergi/go-diff/diffmatchpatch"
	"tensin.org/watchthatpage/difflib"
)

func ComputeDifferences2(previousContent []byte, currentContent []byte) string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(previousContent[:])),
		B:        difflib.SplitLines(string(currentContent[:])),
		FromFile: "Old content",
		ToFile:   "New content",
		Context:  3,
	}
	s, _ := difflib.GetUnifiedDiffString(diff)
	return s
}

func ComputeDifferences(previousContent []byte, currentContent []byte) string {
	// dmp := diffmatchpatch.New()
	// diffs := dmp.DiffMain(string(previousContent[:]), string(currentContent[:]), true)
	// result.Differences = dmp.DiffPrettyText(diffs)

	/*
		array1 := strings.Split(html.EscapeString(string(previousContent[:])), "\n")
		array2 := strings.Split(html.EscapeString(string(currentContent[:])), "\n")
		result.Differences = difflib.HTMLDiff(array1, array2)
	*/

	/*
	 */
	return ""
}
