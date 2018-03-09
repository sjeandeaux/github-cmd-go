# Github Command Lines in Golang

[![Build Status](https://travis-ci.org/sjeandeaux/github-cmd-go.svg)](https://travis-ci.org/sjeandeaux/github-cmd-go) [![Coverage Status](https://coveralls.io/repos/sjeandeaux/github-cmd-go/badge.svg)](https://coveralls.io/r/sjeandeaux/github-cmd-go) [![Go Report Card](https://goreportcard.com/badge/github.com/sjeandeaux/github-cmd-go)](https://goreportcard.com/report/github.com/sjeandeaux/github-cmd-go)

## git-latest

The tool get the latest version in as tag if not found 0.0.0.

```
>git-latest
0.1.0
```

## incrementor

The tool increments the verison.

```
>incrementor -position minor -version 0.1.0
0.2.0
```

## associator

The tools associates the binary to a release in github.

```
>go build $(LDFLAGS) -o ./target/$(1)-$(2)-${APPL} ./cmd/${APPL}
#this command creates a release and attachs the file
>associator -create -name <name> \
                   -label <label> \
                   -content-type  <content-type>\
                   -owner <owner> \
                   -repo <repo> \
                   -tag  <tag>  \
                   -file <file>

#this command attachs the file
>associator -name <name> \
                   -label <label> \
                   -content-type  <content-type>\
                   -owner <owner> \
                   -repo <repo> \
                   -tag  <tag>  \
                   -file <file>
```
