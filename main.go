package main

import (
	"context"
	"fmt"
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

func cmdToken() (err error) {
	tokenNote := "ght-gen generated at" + time.Now().String()

	fmt.Print("Github Username: ")
	var username string
	fmt.Scanln(&username)
	fmt.Print("Github Password: ")
	bytePassword, err := terminal.ReadPassword(0)
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
			return err
		} else {
			break
		}
	}

	fmt.Printf("export GITHUB_TOKEN=\"%s\"\n", *auth.Token)

	return nil
}
