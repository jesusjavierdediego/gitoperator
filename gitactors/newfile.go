package gitactors

import (
	"bytes"
	"encoding/json"
	"errors"
	//"io/fs"
	"io/ioutil"
	"os"
	//"os/exec"
	"path/filepath"
	"time"
	configuration "xqledger/gitoperator/configuration"
	utils "xqledger/gitoperator/utils"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	//"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	// billy "github.com/go-git/go-billy/v5"
	// memfs "github.com/go-git/go-billy/v5/memfs"
	// "github.com/go-git/go-git/v5"
	// git "github.com/go-git/go-git/v5"
	// http "github.com/go-git/go-git/v5/plumbing/transport/http"
	// memory "github.com/go-git/go-git/v5/storage/memory"
)

const componentNewMessage = "Git New File Processor"

var config = configuration.GlobalConfiguration

// This example receives a new file to be added to git
// - Add
// - commit
// - push
func GitProcessNewFile(event *utils.RecordEvent) error {
	var methodMsg = "ProcessNewFile"
	var fileName = event.Id + ".json"
	repoPath, err := GetLocalRepoPath(event)
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error getting path for local clones git repository: "+repoPath)
	}
	var completeFileName = ""
	if len(event.Group) > 0 {
		completeFileName = event.Group + "/" + fileName
	} else {
		completeFileName = fileName
	}
	

	var r *git.Repository
	var openErr error
	r, openErr = git.PlainOpen(repoPath)
	if openErr != nil {
		utils.PrintLogError(openErr, componentNewMessage, methodMsg, "Error opening local Git repository: "+repoPath)
		/*
		Error opening the local repo -> Try to clone the remote repo
		*/
		remoteRepoURL := config.Gitserver.Url + "/" + config.Gitserver.Username + "/" + event.DBName
		utils.PrintLogInfo(componentNewMessage, methodMsg, "remoteRepoURL: " + remoteRepoURL)
		utils.PrintLogInfo(componentNewMessage, methodMsg, "We are going to clone the remote repo if it exists - URL: " + remoteRepoURL)
		cloneErr := Clone(remoteRepoURL, repoPath)
		if cloneErr != nil {
			utils.PrintLogError(cloneErr,  componentNewMessage, methodMsg, "Error cloning the repo: "+repoPath)
			return  cloneErr
			/*
			Error cloning a repo - It does not exist - Init a new repo
			*/
			// var storer *memory.Storage
 			// var fs billy.Filesystem
			// fs.Create(repoPath)
			// newRepo, initErr := git.Init(storer, fs.DirEntry)
			// if initErr != nil {
			// 	utils.PrintLogError(initErr,  componentNewMessage, methodMsg, "Error initializing a new repo: "+repoPath)
			// }
			// newRepo.Push(&git.PushOptions{
			// 	Auth: &http.BasicAuth{
			// 		Username: config.Gitserver.Username,
			// 		Password: config.Gitserver.Password,
			// 	},
			// 	Progress: os.Stdout,
			// })
		}
		r, openErr = git.PlainOpen(repoPath)
		if openErr != nil {
			return openErr
		}
	}

	w, err := r.Worktree()
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error getting Worktree in local Git repository: "+repoPath)
		return err
	}

	// We need a file to commit so let's create a new file inside of the
	// worktree of the project using the go standard library.
	utils.PrintLogInfo(componentNewMessage, methodMsg, "File to process: "+completeFileName)
	var prettyJSON bytes.Buffer
	jsonErr := json.Indent(&prettyJSON, []byte(event.RecordContent), "", "\t")
	if jsonErr != nil {
		utils.PrintLogError(jsonErr, componentNewMessage, methodMsg, "Error in JSON pretty printing")
		return jsonErr
	}

	filePathAndName := filepath.Join(repoPath, completeFileName)
	utils.PrintLogInfo(componentNewMessage, methodMsg, "filePathAndName to process: "+filePathAndName)
	writeFileErr := ioutil.WriteFile(filePathAndName, prettyJSON.Bytes(), 0644)
	if writeFileErr != nil {
		utils.PrintLogError(writeFileErr, componentNewMessage, methodMsg, "Error writing to local file: "+filePathAndName)
		if len(event.Group) > 0 {
			utils.PrintLogInfo(componentNewMessage, methodMsg, "We are going to make the tree if it does not exist")
			makedirErr := os.Mkdir(filepath.Join(repoPath, event.Group), 0755)
			if makedirErr != nil {
				utils.PrintLogError(makedirErr, componentNewMessage, methodMsg, "Error making new dir: "+event.Group)
				return makedirErr
			}
			writeFileErr := ioutil.WriteFile(filePathAndName, prettyJSON.Bytes(), 0644)
			if writeFileErr != nil {
				utils.PrintLogError(writeFileErr, componentNewMessage, methodMsg, "Error writing file in new tree: "+filePathAndName)
				return writeFileErr
			}
			return nil
		}else{
			return errors.New("Event group is empty")
		}
	}

	//PULL FIRST
	w.Pull(&git.PullOptions{RemoteName: "origin"})
	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		utils.PrintLogError(err,componentNewMessage, methodMsg, "Error getting HEAD reference")
		return err
	}
	commitPull, err := r.CommitObject(ref.Hash())
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error in commit - Ref Hash: "+ref.Hash().String())
		return err
	}

	utils.PrintLogInfo(componentNewMessage, methodMsg, commitPull.String())

	utils.PrintLogInfo(componentNewMessage, methodMsg, "File: "+completeFileName)

	utils.PrintLogInfo(componentNewMessage, methodMsg, "git add file")
	_, err = w.Add(completeFileName)
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error in add - File: "+completeFileName)
		return err
	}

	utils.PrintLogInfo(componentNewMessage, methodMsg, "git status --porcelain")
	status, err := w.Status()
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error getting status in local repo")
		return err
	}

	utils.PrintLogInfo(componentNewMessage, methodMsg, status.String())

	// Commits the current staging area to the repository, with the new file just created.
	// We should provide the object.Signature of Author of the commit.
	utils.PrintLogInfo(componentNewMessage, methodMsg, "git commit -m \""+completeFileName+"\"")
	commit, err := w.Commit(completeFileName, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.Gitserver.Username,
			Email: config.Gitserver.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error in commit - ID: "+event.Id)
		return err
	}

	// Prints the current HEAD to verify that all worked well.
	utils.PrintLogInfo(componentNewMessage, methodMsg, "git show -s")
	obj, err := r.CommitObject(commit)
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error in showing commit for verification")
		return err
	}

	utils.PrintLogInfo(componentNewMessage, methodMsg, obj.String())
	utils.PrintLogInfo(componentNewMessage, methodMsg, "git push")

	//utils.PrintLNew, methodMsg, config.Gitserver.Username)
	//utils.PrintLNew, methodMsg, config.Gitserver.Password)

	// push using default options
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: config.Gitserver.Username,
			Password: config.Gitserver.Password,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		utils.PrintLogError(err, componentNewMessage, methodMsg, "Error in push")
		return err
	} 
	return nil
}
