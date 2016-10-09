# Recime Command Line Tool

The `Recime Command Line Tool` allows you to scaffold your bot from terminal. It creates your account based on the email you have provided and lets you deploy it to **Recime** cloud.


## Installation

You can install the `recime-cli` by typing the following command:

NOTE: You should already have [Go](https://golang.org/doc/install) installed and GOPATH, GOROOT set to correct values.

    go get -u github.com/Recime/recime-cli


if you are running OSX you can use [Homebrew](http://brew.sh/) to install go, you can do by typing the following command:

    brew install go


Once **Go** is installed, you can need to add the following lines your $HOME/.bashrc file to ensure that PATH values.

    export GOPATH=$HOME/golang
    export GOROOT=/usr/local/opt/go/libexec
    export PATH=$PATH:$GOPATH/bin
    export PATH=$PATH:$GOROOT/bin


The following commands are currently available in the `Recime Command Line Tool`:

```bash
    create      Scaffolds the bot from an interactive prompt
    deploy      Deploys the bot to Recime cloud
    init        Initializes your Recime account

```

Recime template uses [TypeScript](https://www.typescriptlang.org/docs/tutorial.html), you can combine that with [atom-typescript](https://atom.io/packages/atom-typescript)  to generate javascript files for you and take advantage of auto-complete features.
