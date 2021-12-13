package gitactors

import (
	"errors"
	"fmt"
	utils "xqledger/gitoperator/utils"
)

func GetLocalRepoPath(event *utils.RecordEvent) (string, error) {
	repoPath := config.Gitserver.Localreposlocation
	if !(len(repoPath) > 0) || !(len(event.DBName)>0){
		return "", errors.New(fmt.Sprintf("The path for the local git repo cannot be composed - event.Unit: %s - Root path in config: %s" + event.DBName, repoPath))
	}
	return repoPath + "/" + event.DBName, nil
}
