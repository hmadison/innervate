#!/bin/bash

if [ `uname` == "Darwin" ]; then
    go get
    go build
    cp innervate /usr/local/bin/innervate
    cp __support__/launchd.plist /Library/LaunchDaemons/launchd.plist
    chmod chown root:wheel /Library/LaunchDaemons/launchd.plist
    echo "Run `sudo launchctl -w /Library/LaunchDaemons/launchd.plist` or restart to launch Innervate"
fi
