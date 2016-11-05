# Recime Command Line Tool

The `Recime Command Line Tool` allows you to scaffold your bot from terminal. It creates your account based on the email you have provided and lets you deploy it to **Recime** cloud.


## Installation

You can install the `recime-cli` by typing the following command:

NOTE: You should already have [Go](https://golang.org/doc/install) installed and GOPATH, GOROOT set to correct values.

    go get -u github.com/Recime/recime-cli


if you are running OSX you can use [Homebrew](http://brew.sh/) to install go, you can do by typing the following command:

    brew install go


Once **Go** is installed, you need to add the following lines your $HOME/.bash_profile file to ensure that PATH values. In my case, I have created a go folder under /users/recime and added the following lines:

    export GOPATH=/users/recime/go
    export PATH=/users/recime/go/bin:$PATH
    

The following commands are currently available in the `Recime Command Line Tool`:

```bash
    build       Builds the bot module
    create      Scaffolds the bot from an interactive prompt
    deploy      Deploys the bot to Recime cloud
    init        Initializes your Recime account
    install     Installs the dependencies
```

## License

Copyright 2016 Recime Inc.