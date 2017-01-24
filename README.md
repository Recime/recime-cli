# Recime Command Line Tool  [![Build Status](https://travis-ci.org/Recime/recime-cli.svg?branch=master)](https://travis-ci.org/Recime/recime-cli)

The `Recime Command Line Tool` allows you to scaffold your bot from terminal. It creates your account based on the email you have provided and lets you deploy it to **[Recime](https://recime.ai)** cloud.

## macOS

Install [Homebrew](http://brew.sh/)  by pasting the below command at a Terminal prompt:


    /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"


To install all the Go and Recime CLI, just run the following command, Homebrew takes care of the rest for you:

  brew tap recime/tools && brew install recime-cli


This will install all the dependencies, configure path and install the CLI accessible globally. Once installation is complete


## Windows

`Recime Command Line Tool` is available via [Chocolatey](https://chocolatey.org/). Choco is the package manager for Windows.


You can install chocolatey by typing the following command (Powershell V3+):

    iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex

For older systems and installing from cmd.exe. Please checkout the [installation guide](https://chocolatey.org/install):

Next, type the following to install recime-cli:

    choco install recime-cli

This will install all the dependencies, configure path and install the CLI accessible globally. Once installation is complete, type recime-cli to verify. Please do refreshenv or restart the shell for changes to take effect and dependencies to initialize correctly.


The following commands are currently available in the `Recime Command Line Tool`:

```bash
build       Builds the bot module
config      Initializes bot config vars
config:set  Sets a new or existing config var
create      Scaffolds the bot from an interactive prompt
deploy      Deploys the bot to Recime cloud
init        Initializes your Recime account
install     Installs the dependencies
run         Runs the bot locally
```

## License

Copyright © 2017 Recime Inc.
