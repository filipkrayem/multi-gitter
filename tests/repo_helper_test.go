package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fileName = "test.txt"

func createRepo(t *testing.T, ownerName string, repoName string, dataInFile string) vcmock.Repository {
	tmpDir, err := createDummyRepo(dataInFile, os.TempDir())
	require.NoError(t, err)

	return vcmock.Repository{
		OwnerName: ownerName,
		RepoName:  repoName,
		Path:      tmpDir,
	}
}

func createRepoWithCloneDir(t *testing.T, ownerName string, repoName string, dataInFile string, dir string) vcmock.Repository {
	err := createDirectoryIfDoesntExist(dir)
	require.NoError(t, err)

	dir, err = makeAbsolutePath(dir)
	require.NoError(t, err)

	tmpDir, err := createDummyRepo(dataInFile, dir)
	require.NoError(t, err)

	return vcmock.Repository{
		OwnerName: ownerName,
		RepoName:  repoName,
		Path:      tmpDir,
	}
}

func createDummyRepo(dataInFile string, dir string) (string, error) {
	tmpDir, err := os.MkdirTemp(dir, "multi-git-test-*.git")
	if err != nil {
		return "", err
	}

	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		return "", err
	}

	testFilePath := filepath.Join(tmpDir, fileName)

	err = os.WriteFile(testFilePath, []byte(dataInFile), 0600)
	if err != nil {
		return "", err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	if _, err = wt.Add("."); err != nil {
		return "", err
	}

	_, err = wt.Commit("First commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

func changeBranch(t *testing.T, path string, branchName string, create bool) {
	repo, err := git.PlainOpen(path)
	assert.NoError(t, err)

	wt, err := repo.Worktree()
	assert.NoError(t, err)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: create,
	})
	assert.NoError(t, err)
}

func branchExist(t *testing.T, path string, branchName string) bool {
	repo, err := git.PlainOpen(path)
	assert.NoError(t, err)

	_, err = repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName)), false)
	if err == plumbing.ErrReferenceNotFound {
		return false
	}
	assert.NoError(t, err)

	return true
}

func changeTestFile(t *testing.T, basePath string, content string, commitMessage string) {
	repo, err := git.PlainOpen(basePath)
	require.NoError(t, err)

	testFilePath := filepath.Join(basePath, fileName)

	err = os.WriteFile(testFilePath, []byte(content), 0600)
	require.NoError(t, err)

	wt, err := repo.Worktree()
	require.NoError(t, err)

	_, err = wt.Add(".")
	require.NoError(t, err)

	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
}

func addFile(t *testing.T, basePath string, fn string, content string, commitMessage string) {
	repo, err := git.PlainOpen(basePath)
	require.NoError(t, err)

	testFilePath := filepath.Join(basePath, fn)

	err = os.WriteFile(testFilePath, []byte(content), 0600)
	require.NoError(t, err)

	wt, err := repo.Worktree()
	require.NoError(t, err)

	_, err = wt.Add(".")
	require.NoError(t, err)

	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
}

func readTestFile(t *testing.T, basePath string) string {
	testFilePath := filepath.Join(basePath, fileName)

	b, err := os.ReadFile(testFilePath)
	require.NoError(t, err)

	return string(b)
}

func readFile(t *testing.T, basePath string, fn string) string {
	testFilePath := filepath.Join(basePath, fn)

	b, err := os.ReadFile(testFilePath)
	require.NoError(t, err)

	return string(b)
}

func fileExist(t *testing.T, basePath string, fn string) bool {
	_, err := os.Stat(filepath.Join(basePath, fn))
	if os.IsNotExist(err) {
		return false
	}

	require.NoError(t, err)
	return true
}

func normalizePath(path string) string {
	return strings.ReplaceAll(filepath.ToSlash(path), " ", "\\ ")
}

func createDirectoryIfDoesntExist(directoryPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		// Create the directory
		err := os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func makeAbsolutePath(dirPath string) (string, error) {
	if filepath.IsAbs(dirPath) {
		return dirPath, nil
	}

	absolutePath, err := filepath.Abs(dirPath)
	if err != nil {
		return "", err
	}

	return absolutePath, nil
}

func indexOf(arr []string, target string) int {
	for i, element := range arr {
		if element == target {
			return i
		}
	}
	return -1
}
