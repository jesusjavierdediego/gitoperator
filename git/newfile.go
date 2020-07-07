package git

import (
"fmt"
"time"
"os"
"bytes"
"encoding/json"
"path/filepath"
"io/ioutil"
"gopkg.in/src-d/go-git.v4"
"gopkg.in/src-d/go-git.v4/plumbing/object"
. "gopkg.in/src-d/go-git.v4/_examples"
"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// This example receives a new file to be added to git
// - Add
// - commit
// - push
func NewFile(local_repo_path string, file_name string, content string) {
	r, err := git.PlainOpen(local_repo_path)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)


	// ... we need a file to commit so let's create a new file inside of the
	// worktree of the project using the go standard library.
	Info("filename: " + file_name)
	var prettyJSON bytes.Buffer
    jsonErr := json.Indent(&prettyJSON, []byte(content), "", "\t")
	CheckIfError(jsonErr)
	// TODO check if the file actually exists
	filename := filepath.Join(local_repo_path, file_name)
	err = ioutil.WriteFile(filename, prettyJSON.Bytes(), 0644)
	CheckIfError(err)


	//PULL FIRST
	Info("git pull origin")
	w.Pull(&git.PullOptions{RemoteName: "origin"})
	// Print the latest commit that was just pulled
	ref, err := r.Head()
	CheckIfError(err)
	commitPull, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commitPull)




	Info("git add sample file")
	_, err = w.Add(file_name)
	CheckIfError(err)

	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)

	fmt.Println(status)



	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	var msg = "Added file " + file_name 
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
}