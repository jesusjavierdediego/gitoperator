package gitactors

import (
	"errors"
	utils "xqledger/gitoperator/utils"
	"os"
	"path/filepath"
	"time"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// This example receives a command to delete an existing file in the git repo
// - Remove
// - commit
// - push
func GitDeleteFile(event *utils.RecordEvent) error {
	var methodMsg = "UpdateFile"
	repoPath, err := os.Getwd()
	if err != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Not found current directory")
	}
	var repoName = config.Gitserver.Repository
	var fileName = event.Id + ".json"

	if !(len(repoPath) > 0) {
		utils.PrintLogError(nil, componentConsumerMessage, methodMsg, "Not found match with Unit in event in configuration - event.Unit: "+event.DBName)
		return errors.New("Not found match with Unit in event in configuration - event.Unit: " + event.DBName)
	}
	repoPath = repoPath + "/" + repoName
	var completeFileName = fileName
	if len(event.Group) > 0 {
		completeFileName = event.Group + "/" + fileName
	}

	r, openErr := git.PlainOpen(repoPath)
	if openErr != nil {
		utils.PrintLogError(openErr, componentConsumerMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		/*
		Error opening the local repo -> Try to clone the remote repo
		*/
		remoteRepoURL := config.Gitserver.Url + "/" + config.Gitserver.Username + "/" + repoName

		utils.PrintLogInfo(componentConsumerMessage, methodMsg, "We are going to clone the remote repo if it exists - URL: " + remoteRepoURL)
		cloneErr := Clone(remoteRepoURL, repoPath)
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

	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "Lets go to delete the file - "+completeFileName)
	fileLocalPath := filepath.Join(repoPath, completeFileName)

	
	deleteFileErr := os.Remove(fileLocalPath)
	if deleteFileErr != nil {
		utils.PrintLogError(err, componentConsumerMessage, methodMsg, "Error deleting local file: "+fileLocalPath)
		return deleteFileErr
	}


	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "Deleted file - "+completeFileName)

	//PULL FIRST
	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "git pull origin")
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

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
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
