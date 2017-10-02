#!/bin/bash

if [ `uname` == "Darwin" ]; then
    go get
    go build
    cp innervate /usr/local/bin/innervate
    cp __support__/innervated.plist /Library/LaunchDaemons/innervated.plist
    chown root:wheel /Library/LaunchDaemons/innervated.plist
    echo "Run 'sudo launchctl load -w /Library/LaunchDaemons/innervate.plist' or restart to launch Innervate"
fi
