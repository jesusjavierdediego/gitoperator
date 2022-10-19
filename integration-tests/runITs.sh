#!/bin/bash

# https://dzone.com/articles/viewing-junit-xml-files-locally
# https://www.npmjs.com/package/junit-viewer

if [[ -f "gitea" ]]
then
    echo "The test IT Git repo DOES NOT EXIST in the filesystem. The Integration Test is not possible"
else
    echo "GIT OPERATOR integration tests start"
    docker-compose up -d
    echo "Docker containers for ITs ready."
    sleep 20
    sh prepareTests.sh
    echo "Preparation to tests OK"
    echo "Integration tests start"
    cd ../mediator
    PROFILE=dev go test xqledger/gitoperator/apilogger -v  2>&1 | go-junit-report > ../testreports/apilogger.xml
    PROFILE=dev go test xqledger/gitoperator/configuration -v 2>&1 | go-junit-report > ../testreports/configuration.xml
    PROFILE=dev go test xqledger/gitoperator/utils -v 2>&1 | go-junit-report > ../testreports/utils.xml
    PROFILE=dev go test xqledger/gitoperator/giteaapiclient -v  2>&1 | go-junit-report > ../testreports/giteaapiclient.xml
    echo "Integration tests complete"
    echo "Cleaning up..."
    cd ../integration-tests
    docker-compose down
    echo "Clean up complete. Bye!"
fi