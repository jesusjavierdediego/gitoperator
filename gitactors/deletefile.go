package gitactors

import (
	//"errors"
	utils "xqledger/gitoperator/utils"
	"os"
	"path/filepath"
	"time"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const componentDeleteMessage = "Git Delete File Processor"


// This example receives a command to delete an existing file in the git repo
// - Remove
// - commit
// - push
func GitDeleteFile(event *utils.RecordEvent) error {
	var methodMsg = "UpdateFile"
	var fileName = event.Id + ".json"
	repoPath, err := GetLocalRepoPath(event)
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error getting path for local cloned git repository: "+repoPath)
	}
	repoPath = repoPath + "/" + event.DBName
	var completeFileName = fileName
	if len(event.Group) > 0 {
		completeFileName = event.Group + "/" + fileName
	}

	r, openErr := git.PlainOpen(repoPath)
	if openErr != nil {
		utils.PrintLogError(openErr, componentDeleteMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		/*
		Error opening the local repo -> Try to clone the remote repo
		*/
		remoteRepoURL := config.Gitserver.Url + "/" + config.Gitserver.Username + "/" + event.DBName

		utils.PrintLogInfo(componentDeleteMessage, methodMsg, "We are going to clone the remote repo if it exists - URL: " + remoteRepoURL)
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
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error getting Worktree in local Git repository: "+repoPath)
		return err
	}

	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "Lets go to delete the file - "+completeFileName)
	fileLocalPath := filepath.Join(repoPath, completeFileName)

	
	deleteFileErr := os.Remove(fileLocalPath)
	if deleteFileErr != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error deleting local file: "+fileLocalPath)
		return deleteFileErr
	}


	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "Deleted file - "+completeFileName)

	//PULL FIRST
	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "git pull origin")
	w.Pull(&git.PullOptions{RemoteName: "origin"})

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error getting HEAD reference")
		return err
	}
	commitPull, err := r.CommitObject(ref.Hash())
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error in commit - Ref Hash: "+ref.Hash().String())
		return err
	}
	utils.PrintLogInfo(componentDeleteMessage, methodMsg, commitPull.String())

	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "git add file")
	_, err = w.Add(completeFileName)
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error in add - File: "+completeFileName)
		return err
	}

	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "git status --porcelain")
	status, err := w.Status()
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error getting status in local repo")
		return err
	}

	utils.PrintLogInfo(componentDeleteMessage, methodMsg, status.String())

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "git commit -m \""+event.Message+"\"")
	commit, err := w.Commit(event.Message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Gitserver.Username,
			Email: config.Gitserver.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error in commit - Message: "+event.Message)
		return err
	}

	// Prints the current HEAD to verify that all worked well.
	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "git show -s")
	obj, err := r.CommitObject(commit)
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error in showing commit for verification")
		return err
	}

	utils.PrintLogInfo(componentDeleteMessage, methodMsg, obj.String())
	utils.PrintLogInfo(componentDeleteMessage, methodMsg, "git push")

	// push using default options
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: config.Gitserver.Username,
			Password: config.Gitserver.Password,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		utils.PrintLogError(err, componentDeleteMessage, methodMsg, "Error in push")
		return err
	}
	utils.PrintLogInfo(componentDeleteMessage, methodMsg, utils.Record_delete_git_written_ok)
	return nil
}
