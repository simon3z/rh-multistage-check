# Installing and Upgrading

    $ go get -u github.com/simon3z/rh-multistage-check

You may consider to set your GOBIN path to your local bin directory, e.g.:

    $ GOBIN=$HOME/.local/bin go get -u github.com/simon3z/rh-multistage-check

# Usage

    $ rh-multistage-check --branch <branch> --repository <repo1> --repository <repo2> ...

# Requirements

Go 1.11 is required to build (modules support). You can check your Go version with:

    $ go version
