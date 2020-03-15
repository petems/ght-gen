package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
)

func getENV(value string) string {
	envValue := os.Getenv(value)
	if len(value) == 0 || envValue == "" {
		panic(fmt.Sprintf("No ENV value for %s", value))
	}
	return envValue
}

func TestIntegration(t *testing.T) {
	makeBuild := exec.Command("make", "int_build")
	err := makeBuild.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", "ght-gen-int-testing", err)
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	c, _, err := vt10x.NewVT10XConsole(expect.WithStdout(buf), expect.WithDefaultTimeout(2*time.Second))
	assert.NoError(t, err, "Failed to create a console")
	defer c.Close()

	cmd := exec.Command("./bin/ght-gen-int-testing")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	githubUsername := getENV("GITHUB_USERNAME")
	githubPassword := getENV("GITHUB_PASSWORD")

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		c.ExpectString("Github Username: ")
		c.SendLine(githubUsername)
		c.ExpectString("Github Password: ")
		c.SendLine(githubPassword)
		c.ExpectEOF()
	}()

	err = cmd.Run()

	// Dump the terminal's screen.
	// t.Logf("\n%s", expect.StripTrailingEmptyLines(buf.String()))

	pattern := `export GITHUB_TOKEN="\w+"`

	assert.Contains(t, buf.String(), fmt.Sprintf("Github Username: %s", githubUsername))

	assert.Regexp(t, regexp.MustCompile(pattern), buf.String())

	if err != nil {
		t.Errorf("Expected no error but got '%s'", err)
	}
}
