gpl - Update multiple local repositories with parallel
=======

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Description
gpl will update your selected multiple local repositories with parallel.  
It was developed on the assumption that mainly to rely on other commands.  
Support git, git-svn, svn, mercurial, darcs.

## Synopsis
    gpl -h # help
    gpl -v # version
    gpl /user/local/project/repo1 /user/local/project/repo2

If you have created such as this file

    cat list.txt
    /user/local/project/repo1
    /user/local/project/repo2

You can do this

    cat list.txt | gpl

## Installation
    go get https://github.com/Code-Hex/gpl

