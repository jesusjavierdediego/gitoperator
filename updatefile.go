package main

import (
"fmt"
"time"
"os"
"strings"
"path/filepath"
"io/ioutil"
//"bytes"
//"encoding/json"
"gopkg.in/src-d/go-git.v4"
"gopkg.in/src-d/go-git.v4/plumbing/object"
. "gopkg.in/src-d/go-git.v4/_examples"
"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)


// This example receives a new file to be added to git
// - apply changes to file
// - Add
// - commit
// - push
func UpdateFile(local_repo_path string, file_name string, event RecordEvent) {
	r, err := git.PlainOpen(local_repo_path)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	Info("write content to file")
	filename := filepath.Join(local_repo_path, file_name)
	// OPTION 1: write all to file
	// var prettyJSON bytes.Buffer
	// jsonErr := json.Indent(&prettyJSON, []byte(content), "", "\t")
	// CheckIfError(jsonErr)
	// err = ioutil.WriteFile(filename, prettyJSON.Bytes(), 0644)
	// CheckIfError(err)

	// OPTION 2: replace by line
	//replaceContentInFile(filename)

	// OPTION 3:
	// Pretty ptrint incoming payload  (whole record is needed)
	// Iterate lines in file, compare lines in file with incoming. If ot does not match, replace with the new line
	if len(event.Fields) > 0 {
		for k, v := range event.Fields {
			replaceJSONContentInFileByLine(filename, k, v)
		}
	} 
	
	Info("Written content to file " + file_name)


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
	_, err = w.Add(file_name)
	CheckIfError(err)

	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)

	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	var msg = "Updated file " + file_name 
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

// TODO
func replaceJSONContentInFileByLine(filepath string, key string, value string) {
}

func replaceContentInFile(filepath string) {
	input, err := ioutil.ReadFile(filepath)
        if err != nil {
            fmt.Println("error opening file: ", err)
        }

        lines := strings.Split(string(input), "\n")

        for i, line := range lines {
            if strings.Contains(line, "Wed Mar 03 11:00:08 +0000 2020") {
                lines[i] = "\"created_at\": \"Wed Mar 15 15:00:08 +0000 2020\"," //replace?
            }
        }
        output := strings.Join(lines, "\n")
        err = ioutil.WriteFile(filepath, []byte(output), 0644)
        if err != nil {
			fmt.Println("error writing to file: ", err)
        }
}