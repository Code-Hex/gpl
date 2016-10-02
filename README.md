gpl - Update multiple local repositories with parallel
=======

[![Build Status](https://travis-ci.org/Code-Hex/gpl.svg?branch=master)](https://travis-ci.org/Code-Hex/gpl)
[![Go Report Card](https://goreportcard.com/badge/github.com/Code-Hex/gpl)](https://goreportcard.com/report/github.com/Code-Hex/gpl)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Description
gpl will update your selected multiple local repositories with parallel.  
It was developed on the assumption that mainly to rely on other commands.  
Support git, git-svn, svn, mercurial, darcs.

## Synopsis
    gpl -h # help
    gpl -v # version
    gpl /user/local/project/repo1 /user/local/project/repo2 ...

If you have created such as this file

    cat list.txt
    /user/local/project/repo1
    /user/local/project/repo2

You can do this

    cat list.txt | gpl

## Installation
    go get -u github.com/Code-Hex/gpl/cmd/gpl

## Author
[codehex](https://twitter.com/CodeHex)

