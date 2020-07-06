package main

import (
"fmt"
"time"
"os"
"path/filepath"
"io/ioutil"
"bytes"
"encoding/json"
"gopkg.in/src-d/go-git.v4"
"gopkg.in/src-d/go-git.v4/plumbing/object"
. "gopkg.in/src-d/go-git.v4/_examples"
"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
"github.com/sergi/go-diff/diffmatchpatch"
)


// This example receives a new file to be added to git
// - Pretty print
// - apply changes to file
// - Add
// - commit
// - push
//func UpdateFile(local_repo_path string, file_name string, newRecord string) {
func UpdateFile(local_repo_path string, fileName string, newRecordEvent RecordEvent) {
	if newRecordEvent.OperationType == "update" {
		var prettyJSON bytes.Buffer
		json.Indent(&prettyJSON, []byte(newRecordEvent.RecordContent), "", "\t")
		var prettyNewRecord = string(prettyJSON.Bytes())


		r, err := git.PlainOpen(local_repo_path)
		CheckIfError(err)

		w, err := r.Worktree()
		CheckIfError(err)

		Info("write content to file")
		fileLocalPath := filepath.Join(local_repo_path, fileName)

		replaceContentInFile(fileLocalPath, prettyNewRecord)
		Info("Written content to file " + fileName)

		//PULL FIRST
		Info("git pull origin")
		w.Pull(&git.PullOptions{RemoteName: "origin"})


		// Print the latest commit that was just pulled
		ref, err := r.Head()
		CheckIfError(err)
		commitPull, err := r.CommitObject(ref.Hash())
		CheckIfError(err)
		fmt.Println(commitPull)

		Info("git add file")
		_, err = w.Add(fileName)
		CheckIfError(err)

		Info("git status --porcelain")
		status, err := w.Status()
		CheckIfError(err)

		fmt.Println(status)

		// Commits the current staging area to the repository, with the new file
		// just created. We should provide the object.Signature of Author of the
		// commit.
		var msg = "Updated file " + fileName 
		Info("git commit -m \"" + msg + "\"")
		commit, err := w.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Jesus de Diego",
				Email: "jesus.dediego.erles@gmail.com",
				When:  time.Now(),
			},
		})

		CheckIfError(err)

		// Prints the current HEAD to verify that all worked well.
		Info("git show -s")
		obj, err := r.CommitObject(commit)
		CheckIfError(err)

		fmt.Println(obj)

		Info("git push")
		// push using default options
		err = r.Push(&git.PushOptions{
			Auth: &http.BasicAuth{
				Username: "jdediego",
				Password: "Turing_326",
			},
			Progress: os.Stdout,
		})
		CheckIfError(err)
	} else {
		fmt.Println("Operation type not correct")
	}
}


func replaceContentInFile(filepath string, newContent string) {
	dmp := diffmatchpatch.New()
	oldContentBytes, err := ioutil.ReadFile(filepath)
    if err != nil {
        fmt.Println("error opening file: ", err)
	}
	var oldContent = string(oldContentBytes)
	diffs := dmp.DiffMain(oldContent, newContent, false)
	patches := dmp.PatchMake(oldContent, diffs)
	finalText, _ := dmp.PatchApply(patches, oldContent)

    err = ioutil.WriteFile(filepath, []byte(finalText), 0644)
    if err != nil {
		fmt.Println("error writing to file: ", err)
	}
}