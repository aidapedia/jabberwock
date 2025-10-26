#!/bin/bash

# enter project name
read -p "Enter Git Project ex: github.com/aidapedia/jabberwock: " newProjectName
# check if newName is empty
if [ -z "$newProjectName" ]; then
    echo "newProjectName is empty"
    exit 1
fi

# enter name of application
read -p "Enter Name of Application ex: jabberwock: " newAppName
# check if newAppName is empty
if [ -z "$newAppName" ]; then
    echo "newAppName is empty"
    exit 1
fi

echo "Init Project: $newProjectName with Name: $newAppName ..."

# replace content file that using github.com/aidapedia/jabberwock with newProjectName, especially this file. expect not replace on vendor directory
find . -type f -name "*.go" -not -path "./vendor/*" -exec grep -l "github.com/aidapedia/jabberwock" {} + | xargs -n 1 sed -i "" "s|github.com/aidapedia/jabberwock|$newProjectName|g"
# find and replace on go mod. expect not replace on vendor directory
find . -type f -name "go.mod" -not -path "./vendor/*" -exec grep -l "github.com/aidapedia/jabberwock" {} + | xargs -n 1 sed -i "" "s|github.com/aidapedia/jabberwock|$newProjectName|g"

# find config file main.json in path files/config. expect not replace on vendor directory
find ./files/config -type f -name "*.json" -not -path "./vendor/*" -exec grep -l "jabberwock" {} + | xargs -n 1 sed -i "" "s|jabberwock|$newAppName|g"

# run go mod tidy and vendor
go mod tidy
go mod vendor

# print successfully replace
echo "Replace all jabberwock with $newAppName ..."

# add make generate key
echo "Create private key ..."
make generate-key
echo "Private key created successfully."

# print successfully replace
echo "Successfully init Project: $newProjectName with Name: $newAppName"
echo "Please check and commit the changes"
