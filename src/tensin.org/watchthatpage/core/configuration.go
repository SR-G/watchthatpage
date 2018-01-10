package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Configuration object contains entries read from the provided JSON file
type Configuration struct {
	Urls                 []string
	Selectors            map[string]URLConfiguration
	SectionsToSkip       []string
	LogLevel             string
	CacheDirectory       string
	ScreenshotCommand    string
	ScreenshotParameters string
	Gzip                 bool
	MinifyHTML           bool
	AutoBackup           bool
	GenerateScreenshots  bool
	NotificationByMail   NotificationMail `json:"NotificationMail"`
}

// SortedUrls is a sorted array of all the URLs
type SortedUrls []string

func (a SortedUrls) Len() int           { return len(a) }
func (a SortedUrls) Less(i, j int) bool { return a[i] < a[j] }
func (a SortedUrls) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// URLConfiguration additionnal configuration for one URL
type URLConfiguration struct {
	Selector        string
	SelectorsToSkip []string
}

// NotificationMail contains the mail configuration for notifications by mails
type NotificationMail struct {
	Template     string
	To           string
	From         string
	Subject      string
	SMTPHostname string `json:"smtp-hostname"`
	SMTPTLS      bool   `json:"smtp-tls"`
	SMTPPort     int    `json:"smtp-port"`
	SMTPLogin    string `json:"smtp-login"`
	SMTPPassword string `json:"smtp-password"`
}

func removeDuplicates(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

// InitDefaultConfiguration initializes the configuration with default values
func (conf *Configuration) InitDefaultConfiguration() {
	conf.Urls = []string{}
	conf.LogLevel = "INFO"
	conf.CacheDirectory = "cache/"
}

// LoadFromDefaultLocation loads the configuration from (if available) a JSON configuration file located alongside the program itself
func (conf *Configuration) LoadFromDefaultLocation(configurationFileName string) {
	if _, err := os.Stat(configurationFileName); os.IsNotExist(err) {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configurationFileName = dir + string(os.PathSeparator) + configurationFileName
	}
	conf.InitDefaultConfiguration()
	conf.Load(configurationFileName)
}

// Load loads the configuration from the provided JSON configuration file
func (conf *Configuration) Load(configurationFileName string) {
	if _, err := os.Stat(configurationFileName); err == nil {
		fmt.Println("Configuration file found under [" + configurationFileName + "], now loading content")
		file, _ := os.Open(configurationFileName)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&conf)
		if err != nil {
			fmt.Println("error while loading configuration :", err)
		}
		for url := range conf.Selectors {
			conf.Urls = append(conf.Urls, url)
		}
		conf.Urls = removeDuplicates(conf.Urls)
		sort.Sort(SortedUrls(conf.Urls))

		if !filepath.IsAbs(conf.CacheDirectory) {
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			conf.CacheDirectory = dir + string(os.PathSeparator) + conf.CacheDirectory
		}

		if conf.NotificationByMail.Template != "" && !filepath.IsAbs(conf.NotificationByMail.Template) {
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			conf.NotificationByMail.Template = dir + string(os.PathSeparator) + conf.NotificationByMail.Template
		}

		fmt.Printf("Configuration loaded with [%d] urls, gzip [%t], minify [%t], auto backup [%t], generate screenshots [%t], sections to skip %v\n", len(conf.Urls), conf.Gzip, conf.MinifyHTML, conf.AutoBackup, conf.GenerateScreenshots, conf.SectionsToSkip)
	} else {
		fmt.Println("No external configuration file found under [" + configurationFileName + "] will use default values")
	}
}
