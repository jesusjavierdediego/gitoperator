package gitactors

import (
	"bytes"
	"encoding/json"
	//"errors"
	"io/ioutil"
	utils "xqledger/gitoperator/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const componentUpdateMessage = "Git Update File Processor"
// This example receives a command to update an existing  file into the git repo
// - Pretty print
// - apply changes to file
// - Add
// - commit
// - push
func GitUpdateFile(event *utils.RecordEvent) error {
	var methodMsg = "UpdateFile"
	var fileName = event.Id + ".json"
	repoPath, err := GetLocalRepoPath(event)
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error getting path for local clones git repository: "+repoPath)
	}
	repoPath = repoPath + "/" + event.DBName
	var completeFileName = fileName
	if len(event.Group) > 0 {
		completeFileName = event.Group + "/" + fileName
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, []byte(event.RecordContent), "", "\t")
	var prettyNewRecord = string(prettyJSON.Bytes())

	r, openErr := git.PlainOpen(repoPath)
	if openErr != nil {
		utils.PrintLogError(openErr, componentUpdateMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		/*
		Error opening the local repo -> Try to clone the remote repo
		*/
		remoteRepoURL := config.Gitserver.Url + "/" + config.Gitserver.Username + "/" + event.DBName

		utils.PrintLogInfo(componentUpdateMessage, methodMsg, "We are going to clone the remote repo if it exists - URL: " + remoteRepoURL)
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
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error getting Worktree in local Git repository: "+repoPath)
		return err
	}

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "write content to file - "+completeFileName)
	fileLocalPath := filepath.Join(repoPath, completeFileName)

	replaceContentInFile(fileLocalPath, prettyNewRecord)
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "Written content to file - "+completeFileName)

	//PULL FIRST
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git pull origin")
	w.Pull(&git.PullOptions{RemoteName: "origin"})

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error getting HEAD reference")
		return err
	}
	commitPull, err := r.CommitObject(ref.Hash())
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in commit - Ref Hash: "+ref.Hash().String())
		return err
	}
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, commitPull.String())

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git add file")
	_, err = w.Add(completeFileName)
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in add - File: "+completeFileName)
		return err
	}

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git status --porcelain")
	status, err := w.Status()
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error getting status in local repo")
		return err
	}

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, status.String())

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git commit -m \""+event.Message+"\"")
	commit, err := w.Commit(event.Message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Gitserver.Username,
			Email: config.Gitserver.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in commit - Message: "+event.Message)
		return err
	}

	// Prints the current HEAD to verify that all worked well.
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git show -s")
	obj, err := r.CommitObject(commit)
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in showing commit for verification")
		return err
	}

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, obj.String())
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git push")

	// push using default options
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: config.Gitserver.Username,
			Password: config.Gitserver.Password,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in push")
		return err
	}
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, utils.Record_update_git_written_ok)
	return nil
}

func replaceContentInFile(filepath string, newContent string) {
	var methodMsg = "replaceContentInFile"
	dmp := diffmatchpatch.New()
	oldContentBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error opening local file: "+filepath)
		return
	}
	var oldContent = string(oldContentBytes)
	diffs := dmp.DiffMain(oldContent, newContent, false)
	patches := dmp.PatchMake(oldContent, diffs)
	finalText, _ := dmp.PatchApply(patches, oldContent)

	err = ioutil.WriteFile(filepath, []byte(finalText), 0644)
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error writing to local file: "+filepath)
		return
	}
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "Content in local file updated - File: "+filepath)
}
