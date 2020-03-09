# ght-gen

A basic CLI app to generate an OAuth token for Github. 

### Create a token:

```bash
$ ght-gen
Github Username: janedoe
Github Password: 
export GITHUB_TOKEN="b5b3eany92oxsx1jsvwtkn6fd550dny92oxsx"
```

When 2FA is Enabled:

```
Github Username: janedoe
Github Password: 
2FA Code: 199673
export GITHUB_TOKEN="b5b3eany92oxsx1jsvwtkn6fd550dny92oxsx"
```


## Installation

```bash
go get -v github.com/petems/ght-gen
```

Eventually I'll configure Travis to build binaries and setup a `brew tap` for OSX and Linux.

## Background

As part of my new laptop onboarding, I generate a new SSH key and upload it to Github with Terraform (https://github.com/petems/github-ssh-key-terraform).

Generating a token in the GUI via the website is a bit of a pain, this makes things much easier!