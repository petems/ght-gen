package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

// Version is what is returned by the `-v` flag
const Version = "0.1.0"

// gitCommit is the gitcommit its built from
var gitCommit = "development"

// List of all Github Scopes
// Taken from https://github.com/google/go-github/blob/5a0c02dff48b10470db5615d1a46a7784cb3cb22/github/authorizations.go#L19
// Date: March 9th 2020
var allGithubScopes = []github.Scope{
	github.ScopeAdminGPGKey,
	github.ScopeAdminOrg,
	github.ScopeAdminOrgHook,
	github.ScopeAdminPublicKey,
	github.ScopeAdminRepoHook,
	github.ScopeDeleteRepo,
	github.ScopeGist,
	github.ScopeNotifications,
	github.ScopePublicRepo,
	github.ScopeReadGPGKey,
	github.ScopeReadOrg,
	github.ScopeReadPublicKey,
	github.ScopeReadRepoHook,
	github.ScopeRepo,
	github.ScopeRepoDeployment,
	github.ScopeRepoStatus,
	github.ScopeUser,
	github.ScopeUserEmail,
	github.ScopeUserFollow,
	github.ScopeWriteGPGKey,
	github.ScopeWriteOrg,
	github.ScopeWritePublicKey,
	github.ScopeWriteRepoHook,
}

func main() {
	app := &cli.App{
		Name:    "ght-gen",
		Usage:   "A simple cli to fetch a github api OAuth token",
		Version: gitCommit + "-" + Version,
		Action: func(c *cli.Context) error {
			err := cmdToken()
			return err
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "scopes",
				Usage: fmt.Sprintf("The List of scopes to use.[TODO: Not implemented yet, will do all scopes by default]"),
			},
		},
	}
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func readPassword(prompt string) (pw []byte, err error) {
	fd := int(os.Stdin.Fd())
	if terminal.IsTerminal(fd) {
		fmt.Fprint(os.Stderr, prompt)
		pw, err = terminal.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		return
	}

	fmt.Fprint(os.Stderr, prompt)

	var b [1]byte
	for {
		n, err := os.Stdin.Read(b[:])
		// terminal.ReadPassword discards any '\r', so we do the same
		if n > 0 && b[0] != '\r' {
			if b[0] == '\n' {
				return pw, nil
			}
			pw = append(pw, b[0])
			// limit size, so that a wrong input won't fill up the memory
			if len(pw) > 1024 {
				err = errors.New("password too long")
			}
		}
		if err != nil {
			// terminal.ReadPassword accepts EOF-terminated passwords
			// if non-empty, so we do the same
			if err == io.EOF && len(pw) > 0 {
				err = nil
			}
			return pw, err
		}
	}
}

func cmdToken() (err error) {
	tokenNote := "ght-gen generated at" + time.Now().String()

	fmt.Print("Github Username: ")
	var username string
	fmt.Scanln(&username)
	bytePassword, err := readPassword("Github Password: ")
	if err != nil {
		return err
	}
	password := string(bytePassword)

	transport := github.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	var auth *github.Authorization

	for {
		client := github.NewClient(transport.Client())

		auth, _, err = client.Authorizations.Create(context.Background(),
			&github.AuthorizationRequest{
				Scopes: allGithubScopes,
				Note:   &tokenNote,
			})

		if _, ok := err.(*github.TwoFactorAuthError); ok {
			fmt.Println("2FA Code: ")
			var otp string
			fmt.Scanln(&otp)
			transport.OTP = otp
		} else if err != nil {
			fmt.Println()
			return err
		} else {
			break
		}
	}

	fmt.Printf("export GITHUB_TOKEN=\"%s\"\n", *auth.Token)

	return nil
}
