package core

import(
	"testing"
	"strings"
	tests "tensin.org/watchthatpage/tests"
)

func TestScreenshotCommand(t *testing.T) {
	command := "/usr/bin/docker run --rm -v ${cache}:/images kevinsimper/wkhtmltoimage --quality 75 --format jpg ${url} /images/${filename}.jpg"
	c := PrepareCommand(command, "CACHE", "URL", "FILENAME")
	t.Log("Command is [" + strings.Join(c[:], " ") + "]")
	tests.AssertEquals(t, "/usr/bin/docker run --rm -v CACHE:/images kevinsimper/wkhtmltoimage --quality 75 --format jpg URL /images/FILENAME.jpg", strings.Join(c[:], " "))
}