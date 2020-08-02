package gitactors

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	configuration "me/gitoperator/configuration"
	utils "me/gitoperator/utils"
	"os"
	"path/filepath"
	"time"
    //"fmt"
	//"github.com/go-git/go-git/v5"
	//"github.com/go-git/go-git/v5/plumbing"
	//"github.com/go-git/go-git"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const componentConsumerMessage = "Git New File Processor"

var config = configuration.GlobalConfiguration

// This example receives a new file to be added to git
// - Add
// - commit
// - push
func GitProcessNewFile(event *utils.RecordEvent) error {
	var methodMsg = "ProcessNewFile"
	var repoPath = ""
	var repoName = ""
	var fileName = event.Id + ".json"

	//utils.PrintLogInfo(componentConsumerMessage, methodMsg, "GIT NEW FILE1")
	
	for _, dbowner := range config.Dbowners {
		if dbowner.Repo == event.DBName {
			repoName = dbowner.Repo
			repoPath = config.Localgitbasicpath + dbowner.Repo
		}
	}
	if !(len(repoPath) > 0) {
		utils.PrintLogError(nil, componentConsumerMessage, methodMsg, "Not found match with Unit in event in configuration - event.Unit: "+event.DBName)
		return errors.New("Not found match with Unit in event in configuration - event.Unit: " + event.DBName)
	}

	var completeFileName = event.Group + "/" + fileName

	var r *git.Repository
	var openErr error
	r, openErr = git.PlainOpen(repoPath)
	if openErr != nil {
		utils.PrintLogError(openErr, componentConsumerMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		/*
		Error opening the local repo -> Try to clone the remote repo
		*/
		remote_repo_url := config.Gitserver.Url + "/" + config.Gitserver.Username + "/" + repoName

		utils.PrintLogInfo(componentConsumerMessage, methodMsg, "We are going to clone the remote repo if it exists - URL: " + remote_repo_url)
		cloneErr := Clone(remote_repo_url, repoPath)
		if cloneErr != nil {
			return cloneErr
		}
		r, openErr = git.PlainOpen(repoPath)
		if openErr != nil {
			return openErr
		}
	}

	w, err := r.Worktree()
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error getting Worktree in local Git repository: "+repoPath)
		return err
	}

	// We need a file to commit so let's create a new file inside of the
	// worktree of the project using the go standard library.
	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "File to process: "+completeFileName)
	var prettyJSON bytes.Buffer
	jsonErr := json.Indent(&prettyJSON, []byte(event.RecordContent), "", "\t")
	if jsonErr != nil {
		utils.PrintLogError(jsonErr, componentConsumerMessage, methodMsg, "Error in JSON pretty printing")
		return jsonErr
	}

	filePathAndName := filepath.Join(repoPath, completeFileName)
	err = ioutil.WriteFile(filePathAndName, prettyJSON.Bytes(), 0644)
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error writing to local file: "+filePathAndName)
		if len(event.Group) > 0 {
			utils.PrintLogInfo(componentConsumerMessage, methodMsg, "We are going to make the tree if it does not exist")
			makedirErr := os.Mkdir(filepath.Join(repoPath, event.Group), 0755)
			if makedirErr != nil {
				utils.PrintLogError(makedirErr, componentConsumerMessage, methodMsg, "Error making new dir: "+event.Group)
				return makedirErr
			}
			writeFileErr := ioutil.WriteFile(filePathAndName, prettyJSON.Bytes(), 0644)
			if writeFileErr != nil {
				utils.PrintLogError(writeFileErr, componentConsumerMessage, methodMsg, "Error writing file in new tree: "+filePathAndName)
				return writeFileErr
			}
			return nil
		}else{
			return err
		}
	}

	//PULL FIRST
	w.Pull(&git.PullOptions{RemoteName: "origin"})
	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error getting HEAD reference")
		return err
	}
	commitPull, err := r.CommitObject(ref.Hash())
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error in commit - Ref Hash: "+ref.Hash().String())
		return err
	}

	utils.PrintLogInfo(componentConsumerMessage, methodMsg, commitPull.String())

	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "git add file")
	_, err = w.Add(completeFileName)
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error in add - File: "+completeFileName)
		return err
	}

	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "git status --porcelain")
	status, err := w.Status()
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error getting status in local repo")
		return err
	}

	utils.PrintLogInfo(componentConsumerMessage, methodMsg, status.String())

	// Commits the current staging area to the repository, with the new file just created.
	// We should provide the object.Signature of Author of the commit.
	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "git commit -m \""+event.Message+"\"")
	commit, err := w.Commit(event.Message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Gitserver.Username,
			Email: config.Gitserver.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error in commit - Message: "+event.Message)
		return err
	}

	// Prints the current HEAD to verify that all worked well.
	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "git show -s")
	obj, err := r.CommitObject(commit)
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error in showing commit for verification")
		return err
	}

	utils.PrintLogInfo(componentConsumerMessage, methodMsg, obj.String())
	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "git push")

	// push using default options
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: config.Gitserver.Username,
			Password: config.Gitserver.Password,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error in push")
		return err
	}
	return nil
}
