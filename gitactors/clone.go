package gitactors

import (
	"fmt"
	"os"
	"strings"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	utils "xqledger/gitoperator/utils"
)

const componentCloneMessage = "Git Clone"
// Example of an specific use case:
// - Clone a repository in a specific path
// - Get the HEAD reference
// - Using the HEAD reference, obtain the commit this reference is pointing to
// - Print the commit content
// - Using the commit, iterate over all its files and print them
// - Print all the commit history with commit messages, short hash and the
// first line of the commit message
func Clone(remote_repo_url string, local_repo_path string) error{
	var methodMsg = "Clone"
	// Clone the given repository, creating the remote, the local branches and fetching the objects, exactly as:
	utils.PrintLogInfo(componentCloneMessage, methodMsg, fmt.Sprintf("Clone the given repository - git clone %s %s", remote_repo_url, local_repo_path))

	r, err := git.PlainClone(local_repo_path, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: config.Gitserver.Username,
			Password: config.Gitserver.Password,
		},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		URL: remote_repo_url,
		Progress: os.Stdout,
	})
	if err != nil {
		utils.PrintLogError(err, componentCloneMessage, methodMsg, fmt.Sprintf("Error cloning the repo"))
		return err
	}

	utils.PrintLogInfo(componentCloneMessage, methodMsg, fmt.Sprintf("Getting the latest commit on the current branch"))

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		utils.PrintLogError(err, componentCloneMessage, methodMsg, fmt.Sprintf("Error retrieving the branch being pointed by HEAD"))
		return err
	}

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		utils.PrintLogError(err, componentCloneMessage, methodMsg, fmt.Sprintf("Error retrieving the commit object"))
		return err
	}
	utils.PrintLogInfo(componentCloneMessage, methodMsg, fmt.Sprintln(commit))
	utils.PrintLogInfo(componentCloneMessage, methodMsg, fmt.Sprint("List the tree from HEAD"))

	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		utils.PrintLogError(err, componentCloneMessage, methodMsg, fmt.Sprintf("Error retrieving the tree from the commit"))
		return err
	}

	// ... get the files iterator and print the file
	tree.Files().ForEach(func(f *object.File) error {
		fmt.Printf("100644 blob %s    %s\n", f.Hash, f.Name)
		return nil
	})

	//Info("git log --oneline")
	utils.PrintLogInfo(componentCloneMessage, methodMsg, fmt.Sprintf("List the history of the repository"))

	commitIter, err := r.Log(&git.LogOptions{From: commit.Hash})
	if err != nil {
		utils.PrintLogError(err, componentCloneMessage, methodMsg, fmt.Sprintf("Error retrieving the history"))
		return err
	}

	err = commitIter.ForEach(func(c *object.Commit) error {
		hash := c.Hash.String()
		line := strings.Split(c.Message, "\n")
		utils.PrintLogInfo(componentCloneMessage, methodMsg, fmt.Sprintln(hash[:7], line[0]))

		return nil
	})
	if err != nil {
		utils.PrintLogError(err, componentCloneMessage, methodMsg, fmt.Sprintf("Error making the list of commits"))
		return err
	}
	return nil
}
