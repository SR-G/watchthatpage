package core

// Results is a sorted collection of Result objects
type Results []Result

func (results Results) Len() int {
	return len(results)
}

func (results Results) Less(i, j int) bool {
	return results[i].URL < results[j].URL
}

func (results Results) Swap(i, j int) {
	results[i], results[j] = results[j], results[i]
}

// GlobalResults contains both the information about the program execution results and the list of Result entries (one for each analyzed URL)
type GlobalResults struct {
	NbUrls        int
	NbDiff        int
	NbErrors      int
	ExecutionTime string
	Date          string
	Version       string
	Results       Results
}
