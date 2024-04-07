#!/bin/bash

# This is a sample Bash script for WSL

# Function to greet the user
echo "Building docker image"

docker build -t lambda-build .

echo "Delete old files" 
rm ./infrastructure/bootstraps/bootstrap
rm -rf ./infrastructure/bootstraps/templates
rm ./infrastructure/bootstrap.zip

echo "Copy new files"
cp .env ./infrastructure/bootstraps/.env
cp -r templates ./infrastructure/bootstraps/templates/


echo "Create new compiled file"
docker run -it -v $(pwd):/project lambda-build go build -tags lambda.norpc -o ./infrastructure/bootstraps/bootstrap ./main.go