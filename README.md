# This project is depricated You should use this [fork](https://github.com/Recime/nodejs-cli) instead.

# Recime Command Line Tool  

The `Recime Command Line Tool` allows you to scaffold your bot from terminal. It creates your account based on the email you have provided and lets you deploy it to **[Recime](https://recime.ai)** cloud.

[![Build Status](https://travis-ci.org/Recime/recime-cli.svg?branch=master)](https://travis-ci.org/Recime/recime-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/recime/recime-cli)](https://goreportcard.com/report/github.com/recime/recime-cli)

## macOS

Install [Homebrew](http://brew.sh/)  by pasting the below command at a Terminal prompt:


    /usr/bin/ru by -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"


To install all the Go and Recime CLI, just run the following command, Homebrew takes care of the rest for you:

  brew tap recime/tools && brew install recime-cli


This will install all the dependencies, configure path and install the CLI accessible globally. Once installation is complete

## Source

`Recime Command Line Tool` is built using [Go](https://golang.org/dl/). Therefore, if you have go tools installs then type the following command to install:

```
go get -v -u github.com/recime/recime-cli 

```

Plese make sure that `go` is installed and **GOROOT** and **GO_PATH** is configured correctly and you will need go 1.6+ to install `recime-cli`.


The following commands are currently available in the `Recime Command Line Tool`:

```bash
  config      Initializes bot config vars
  create      Scaffolds the bot from an interactive prompt
  deploy      Deploys the bot to Recime cloud
  login       Logs into your Recime account
  platform    Configures the Platform
  run         Runs the bot locally
```

## License

Copyright © 2017 Recime Inc. All rights reserved.
