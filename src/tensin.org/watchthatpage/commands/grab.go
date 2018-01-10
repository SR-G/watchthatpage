package commands

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"

	"tensin.org/watchthatpage/core"
)

func init() {
	RootCmd.AddCommand(grabCmd)
	grabCmd.PersistentFlags().BoolVarP(&cleanCachedContent, "clean", "c", false, "force clean cached content")
	grabCmd.PersistentFlags().BoolVarP(&activateDebug, "debug", "", false, "print additionnal debug")
	grabCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "", false, "removes colored output")
}

func renderNode(n *html.Node) []byte {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.Bytes() // buf.String()
}

func getTag(doc *html.Node, tag string) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	var found = false
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			b = n
			found = true
		}
		if !found {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Tag [" + tag + "] not found")
}

func removeTags(doc *html.Node, tags []string) (*html.Node, error) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		var nodesToRemove []*html.Node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			removed := false
			if c.Type == html.CommentNode {
				removed = true
			} else if c.Type == html.ElementNode {
				if c.Data == "" {
					removed = true
				} else {
					for _, tag := range tags {
						if c.Data == tag {
							removed = true
						}
					}
				}
			}
			if !removed {
				// keep this node and recurse
				f(c)
			} else {
				// we can't remove while iterating
				nodesToRemove = append(nodesToRemove, c)
			}
		}

		for _, c := range nodesToRemove {
			n.RemoveChild(c)
		}

	}
	f(doc)
	return doc, nil
}

// Makes a string safe to use in a file name (e.g. for saving file atttachments)
func sanitizeFileName(text string) string {
	// Start with lowercase string
	fileName := strings.ToLower(text)
	fileName = path.Clean(path.Base(fileName))
	fileName = strings.Trim(fileName, " ")

	// Replace certain joining characters with a dash
	seps, err := regexp.Compile(`[ &_=+:]`)
	if err == nil {
		fileName = seps.ReplaceAllString(fileName, "-")
	}

	// Remove all other unrecognised characters - NB we do allow any printable characters
	legal, err := regexp.Compile(`[^[:alnum:]-.]`)
	if err == nil {
		fileName = legal.ReplaceAllString(fileName, "")
	}

	// Remove any double dashes caused by existing - in name
	fileName = strings.Replace(fileName, "--", "-", -1)

	// NB this may be of length 0, caller must check
	return fileName
}

func buildFileName(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return sanitizeFileName(hex.EncodeToString(hasher.Sum(nil)))
}

func grabContentFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Can't parse [" + url + "]")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Can't parse [" + url + "]")
	}
	defer resp.Body.Close()
	return b, nil
}

func extractContent(content []byte, selector string, selectorsToSkip []string, sectionsToSkip []string) ([]byte, string) {
	doc, _ := html.Parse(strings.NewReader(string(content[:])))
	queryableFullDoc := goquery.NewDocumentFromNode(doc)
	title := queryableFullDoc.Find("title").Contents().Text()

	bn, err := getTag(doc, "body")
	if err != nil {
		fmt.Println("error", err)
	}

	bn2, err := removeTags(bn, sectionsToSkip)
	if err != nil {
		fmt.Println("error", err)
	}

	if selector != "" {
		queryDoc := goquery.NewDocumentFromNode(bn2)
		selected := queryDoc.Find(selector)
		if len(selected.Nodes) > 0 {
			bn2 = selected.Nodes[0]
		}
	}

	if len(selectorsToSkip) > 0 {
		for _, selector := range selectorsToSkip {
			queryDoc := goquery.NewDocumentFromNode(bn2)
			selected := queryDoc.Find(selector)
			for _, node := range selected.Nodes {
				node.Parent.RemoveChild(node)
			}
		}
	}

	// Removes non-breaksing space \u00a0 (TODO : maybe should be done with transformers)
	// Re-creates lines based on HTML tags (in order to be able to compute differencies)
	var text []byte
	text = renderNode(bn2)
	var s string
	s = strings.Replace(string(text[:]), "\u00a0", " ", -1)
	// s = strings.Replace(s, "</", "\n</", -1)
	s = strings.Replace(s, ">", ">\n", -1)

	return []byte(s), title
}

func crawl(url string, urlConfiguration core.URLConfiguration, ch chan core.Result) {
	startingTime := time.Now().UTC()
	var result = core.Result{}
	result.URL = url
	result.HasDifferences = false
	result.HasError = false

	cacheBaseName := buildFileName(url)
	result.CacheFileName = configuration.CacheDirectory + string(os.PathSeparator) + cacheBaseName
	var buf bytes.Buffer
	buf.WriteString("Now parsing URL [" + url + "]")

	if urlConfiguration.Selector != "" {
		buf.WriteString(", selector [" + urlConfiguration.Selector + "]")
	}
	if len(urlConfiguration.SelectorsToSkip) > 0 {
		buf.WriteString(", excludes [")
		sep := ""
		for _, item := range urlConfiguration.SelectorsToSkip {
			buf.WriteString(sep)
			buf.WriteString(item)
			sep = ", "
		}
		buf.WriteString("]")
	}
	fmt.Println(buf.String())

	b, err := grabContentFromURL(url)
	if err != nil {
		result.HasError = true
	} else {
		currentContent, title := extractContent(b, urlConfiguration.Selector, urlConfiguration.SelectorsToSkip, configuration.SectionsToSkip)
		result.Title = title

		var dest = result.CacheFileName
		cacheExists, err := core.Exists(dest)
		if err != nil {
			fmt.Println("Error : ", err)
		}
		if cleanCachedContent || !cacheExists {
			result.FirstExecution = true
			core.CacheWrite(dest, currentContent, configuration.Gzip, configuration.MinifyHTML, configuration.AutoBackup)
		} else {
			result.FirstExecution = false
			previousContent, err := core.CacheRead(dest, configuration.Gzip)
			if err != nil {
				fmt.Print("Error : ", err)
			}
			if configuration.MinifyHTML {
				currentContent, err = core.Minify(currentContent)
				if err != nil {
					fmt.Print("Error : ", err)
				}
			}
			if bytes.Equal(previousContent, currentContent) {
				// same content
			} else {
				// different content
				core.CacheWrite(dest, currentContent, configuration.Gzip, configuration.MinifyHTML, configuration.AutoBackup)
				result.HasDifferences = true

				if configuration.GenerateScreenshots {
					core.GenerateScreenshot(configuration.ScreenshotCommand, configuration.CacheDirectory, url, cacheBaseName)
					result.ScreenshotFullFileName = configuration.CacheDirectory + string(os.PathSeparator) + cacheBaseName + ".jpg"
					result.ScreenshotFileName = cacheBaseName + ".jpg"
				}

				result.Differences = core.ComputeDifferences(previousContent, currentContent)
			}
		}
	}

	endingTime := time.Now().UTC()
	result.AnalysisExecutionTime = endingTime.Sub(startingTime)

	ch <- result
}

var configuration = core.Configuration{}
var cleanCachedContent bool
var activateDebug bool
var noColor bool
var grabCmd = &cobra.Command{
	Use:   "grab",
	Short: "Grab pages",
	Long:  `Grab the content of configured web pages and trigger notifications`,
	Run: func(cmd *cobra.Command, args []string) {
		startingTime := time.Now().UTC()

		configuration.LoadFromDefaultLocation(ConfigurationFileName)
		core.CreateDirectoryIfNeeded(configuration.CacheDirectory)

		var nb = len(configuration.Urls)
		var results = make(core.Results, nb)

		chResults := make(chan core.Result)
		for _, url := range configuration.Urls {
			go crawl(url, configuration.Selectors[url], chResults)
		}

		for c := 0; c < nb; {
			select {
			case result := <-chResults:
				results[c] = result
				c++
			}
		}

		fmt.Println("Results : ")
		sort.Sort(results)
		var globalResults = core.GlobalResults{}
		globalResults.Results = results
		for _, result := range results {
			globalResults.NbUrls++
			if result.HasError {
				globalResults.NbErrors++
			} else if result.HasDifferences {
				globalResults.NbDiff++
			}
			if noColor {
				fmt.Println("  - " + result.String())
			} else {
				fmt.Println("  - " + result.ColoredString())
			}
		}

		endingTime := time.Now().UTC()
		globalResults.Date = time.Now().String()
		globalResults.ExecutionTime = endingTime.Sub(startingTime).String()
		globalResults.Version = core.Version.String()
		fmt.Printf("Total execution time [%s], analyzed urls [%d], errors [%d], diffs [%d]\n", globalResults.ExecutionTime, globalResults.NbUrls, globalResults.NbErrors, globalResults.NbDiff)

		if globalResults.NbDiff > 0 {
			core.Notify(&globalResults, &configuration, activateDebug)
		} else {
			fmt.Println("No differences - no notification")
		}
	},
}
