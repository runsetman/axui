#!/bin/bash

USERNAME=$1
PASSWORD=$2
TOKEN=$3

wget https://github.com/runsetman/axui/releases/download/v0.3.204/x-ui-linux-amd64.tar.gz

tar xvfz x-ui-linux-amd64.tar.gz -C ./

x-ui stop

rm -rf /usr/local/x-ui/x-ui

mv x-ui/x-ui /usr/local/x-ui/x-ui

/usr/local/x-ui/x-ui setting -username $USERNAME -password $PASSWORD -token $TOKEN

x-ui start