package core

import (
	"fmt"
	"testing"

	tests "tensin.org/watchthatpage/tests"
)

func TestVersion(t *testing.T) {
	Version.Major = 2
	Version.Minor = 3
	Version.Patch = 4
	Version.Label = "TEST"
	Version.Name = "TestVersion"
	Build = "TODAY"
	fmt.Println("Version [" + Version.String() + "]")
	tests.AssertEquals(t, "WatchThatPage version 2.3.4-TEST \"TestVersion\"\nGit commit hash: TODAY", Version.String())
}
