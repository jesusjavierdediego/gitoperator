package git

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

	//"github.com/go-git/go-git/v5"
	//"github.com/go-git/go-git/v5/plumbing"

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
	var fileName = event.Id + ".json"
	
	for _, unit := range config.Units {
		if unit.Name == event.Unit {
			repoPath = config.Gitserver.Localbasicpath + unit.Repo
		}
	}
	if !(len(repoPath) > 0) {
		utils.PrintLogError(nil, componentConsumerMessage, methodMsg, "Not found match with Unit in event in configuration - event.Unit: "+event.Unit)
		return errors.New("Not found match with Unit in event in configuration - event.Unit: " + event.Unit)
	}

	var completeFileName = event.Group + "/" + fileName

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		return err
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
	// TODO check if the file actually exists
	filePathAndName := filepath.Join(repoPath, completeFileName)
	err = ioutil.WriteFile(filePathAndName, prettyJSON.Bytes(), 0644)
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error writing to local file: "+filePathAndName)
		return err
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
