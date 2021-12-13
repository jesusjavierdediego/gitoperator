package gitactors

import (
	"bytes"
	"encoding/json"
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

const componentUpdateBatchMessage = "Git Update Batch Processor"
// This example receives a command to update an existing  file into the git repo
// - Pretty print
// - apply changes to file
// - Add
// - commit
// - push
func GitUpdateFileBatch(batch *utils.RecordEventBatch) error {
	
	var methodMsg = "GitUpdateFileBatch"

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "We are going to update the file")
	
	repoPath := config.Gitserver.Localreposlocation + "/" + batch.DBName

	r, openErr := git.PlainOpen(repoPath)
	if openErr != nil {
		utils.PrintLogError(openErr, componentUpdateMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		/*
		Error opening the local repo -> Try to clone the remote repo
		*/
		remoteRepoURL := config.Gitserver.Url + "/" + config.Gitserver.Username + "/" + batch.DBName

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
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error getting Worktree in local Git repository")
		return err
	}

	//PULL FIRST
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git pull origin")
	w.Pull(&git.PullOptions{RemoteName: "origin"})


	// Compose files in batch for update
	for _, record := range batch.Records {
		processUpdateFileInBatch(&record, repoPath)
	}

	// ADD FILE
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git add file")
	_, err = w.Add(".")
	if err != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in add all files in batch")
		return err
	}

	// COMMIT
	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git commit -m \""+batch.Id+"\"")
	_, commitErr := w.Commit(batch.Id, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Gitserver.Username,
			Email: config.Gitserver.Email,
			When:  time.Now(),
		},
	})
	if commitErr != nil {
		utils.PrintLogError(err, componentUpdateMessage, methodMsg, "Error in commit - Batch ID: "+batch.Id)
		return err
	}
	
	// PUSH
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "git push")
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

func processUpdateFileInBatch(event *utils.RecordEvent, repoPath string) {
	var methodMsg = "processUpdateFileInBatch"
	var fileName = event.Id + ".json"

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "We are going to update the file")
	
	var completeFileName = ""
	if len(event.Group) > 0 {
		completeFileName = event.Group + "/" + fileName
	} else {
		completeFileName = fileName
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, []byte(event.RecordContent), "", "\t")
	var prettyNewRecord = string(prettyJSON.Bytes())

	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "write content to file - "+completeFileName)
	filePathAndName := filepath.Join(repoPath, completeFileName)
	utils.PrintLogInfo(componentNewMessage, methodMsg, "filePathAndName to process: "+filePathAndName)

	replaceContentInFileBatch(filePathAndName, prettyNewRecord)
	utils.PrintLogInfo(componentUpdateMessage, methodMsg, "Written content to file - "+completeFileName)
}

func replaceContentInFileBatch(filepath string, newContent string) {
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
