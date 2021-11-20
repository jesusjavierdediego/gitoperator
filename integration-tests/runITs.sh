#!/bin/bash

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
    PROFILE=dev go test xqledger/gitoperator/mediator -v  2>&1 | go-junit-report > ../testreports/mediator.xml
    PROFILE=dev go test xqledger/gitoperator/gitactors -v -run TestGetLocalRepoPath  2>&1 | go-junit-report > ../testreports/TestGetLocalRepoPath.xml
    PROFILE=dev go test xqledger/gitoperator/gitactors -v -run TestCloneRepo  2>&1 | go-junit-report > ../testreports/TestCloneRepo.xml
    PROFILE=dev go test xqledger/gitoperator/gitactors -v -run TestGitProcessNewFile  2>&1 | go-junit-report > ../testreports/TestGitProcessNewFile.xml
    PROFILE=dev go test xqledger/gitoperator/gitactors -v -run TestGitProcessUpdatedFile  2>&1 | go-junit-report > ../testreports/TestGitProcessUpdatedFile.xml
    PROFILE=dev go test xqledger/gitoperator/gitactors -v -run TestGitProcessDeleteFile  2>&1 | go-junit-report > ../testreports/TestGitProcessDeleteFile.xml
    echo "Integration tests complete"
    echo "Cleaning up..."
    cd ../integration-tests
    docker-compose down
    echo "Clean up complete. Bye!"
fi