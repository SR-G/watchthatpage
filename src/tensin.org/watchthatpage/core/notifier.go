package core

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"gopkg.in/gomail.v2"
)

func renderTemplatedBody(templateFileName string, results *GlobalResults) (string, error) {
	templateName := filepath.Base(templateFileName)
	tmpl := template.New(templateName)
	tmpl.Funcs(template.FuncMap{
		"mod":       func(i, j int) bool { return i%j == 0 },
		"mod_start": func(i, j int) bool { return i%j == 0 },
		"mod_end":   func(i, j int) bool { return (i+1)%j == 0 },
		"add":       func(x, y int) int { return x + y }})
	tmpl, err := tmpl.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &results)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	var minified []byte
	minified, err = m.Bytes("text/html", buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(minified[:]), nil
}

func renderTemplatedSubject(subject string, results *GlobalResults) (string, error) {
	tmpl, err := template.New("subject").Parse(subject)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &results)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// Notify triggers the notification, embedding the provided results
func Notify(results *GlobalResults, configuration *Configuration, debug bool) {
	fmt.Println("Notifications to [" + configuration.NotificationByMail.To + "], from [" + configuration.NotificationByMail.From + "], server [" + configuration.NotificationByMail.SMTPHostname + ":" + strconv.Itoa(configuration.NotificationByMail.SMTPPort) + "], template [" + configuration.NotificationByMail.Template + "]")

	m := gomail.NewMessage()
	title, err := renderTemplatedSubject(configuration.NotificationByMail.Subject, results)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	m.SetHeader("Subject", title)
	m.SetHeader("From", configuration.NotificationByMail.From)
	m.SetHeader("To", configuration.NotificationByMail.To)
	m.SetHeader("MIME-Version", "1.0")
	m.SetHeader("Content-Transfer-Encoding", "base64")
	for _, result := range results.Results {
		if result.ScreenshotFullFileName != "" && result.HasDifferences {
			m.Embed(result.ScreenshotFullFileName) // attach image with full name
		}
	}
	body, err := renderTemplatedBody(configuration.NotificationByMail.Template, results)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	if debug {
		fmt.Println("Content : \n" + body)
	}

	m.SetBody("text/html; charset=\"utf-8\"", body)

	d := gomail.NewPlainDialer(configuration.NotificationByMail.SMTPHostname, configuration.NotificationByMail.SMTPPort, configuration.NotificationByMail.SMTPLogin, configuration.NotificationByMail.SMTPPassword)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error", err)
	}
}
