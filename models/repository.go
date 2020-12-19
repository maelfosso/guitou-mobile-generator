package models

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const guitouURL = "https://gitlab.com/guitou-daco-tool/guitou-mobile.git"

var auth http.BasicAuth = http.BasicAuth{
	Username: "xxxxx",
	Password: "xxxxx",
}

// ProjectRepository of Gitlab gitou-mobile for the project
type ProjectRepository struct {
	name       string
	directory  string
	pkg        string
	repository *git.Repository
}

// NewProjectRepository create an instance the ProjectRepository struct
func NewProjectRepository(project *Project) *ProjectRepository {
	aux := strings.ToLower(strings.ReplaceAll(project.Name, " ", "_"))

	return &ProjectRepository{
		name:       project.Name,
		pkg:        aux,
		directory:  fmt.Sprintf("/tmp/guitou-%s", aux),
		repository: nil,
	}
}

// Clone the project from Gitlab
func (repo *ProjectRepository) Clone() error {
	if _, err := os.Stat(repo.directory); !os.IsNotExist(err) {
		log.Println(fmt.Sprintf("Remove the destination cloning repository"))

		err := os.RemoveAll(repo.directory)
		if err != nil {
			return &ErrorProjectRepository{
				Code:    11,
				Message: fmt.Sprintf("Impossible to delete [%s] folder", repo.directory),
				Err:     err,
			}
		}
	}

	r, err := git.PlainClone(repo.directory, false, &git.CloneOptions{
		URL:  guitouURL,
		Auth: &auth,
		// ReferenceName: plumbing.ReferenceName("refs/heads/develop"),
		// SingleBranch: true,
		Progress: os.Stdout,
	})

	if err != nil {
		return &ErrorProjectRepository{
			Code:    12,
			Message: fmt.Sprintf("Error occurred when cloning [%s] for [%s]", guitouURL, repo.name),
			Err:     err,
		}
	}
	repo.repository = r

	return nil
}

// Checkout create a new branch in which the change will be done
func (repo *ProjectRepository) Checkout() error {
	branch := plumbing.ReferenceName(fmt.Sprintf("refs/heads/deploy/%s-%v", repo.pkg, time.Now().Unix()))
	log.Println("[Checkout] ReferenceName")

	w, err := repo.repository.Worktree()
	if err != nil {
		return &ErrorProjectRepository{
			Code:    21,
			Message: fmt.Sprintf("Error occurred when Worktree for [%s]", repo.name),
			Err:     err,
		}
	}
	log.Println("[Checkout] Worktree")

	err = w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  false,
		Branch: branch,
	})
	if err != nil {
		log.Println("Checkout Failed")
		err = w.Checkout(&git.CheckoutOptions{
			Create: true,
			Force:  false,
			Branch: branch,
		})

		if err != nil {
			return &ErrorProjectRepository{
				Code:    22,
				Message: fmt.Sprintf("Error occurred when Checking our for [%s]", repo.name),
				Err:     err,
			}
		}
	}

	return nil
}

// BashUpdate run the bash script for updating the package
func (repo *ProjectRepository) BashUpdate() error {
	var out bytes.Buffer

	cmd := exec.Command("bash", "guitou-update.sh", repo.pkg, repo.name)
	cmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		return &ErrorProjectRepository{
			Code:    31,
			Message: fmt.Sprintf("Error occurred.when `Cmd.Start()` for [%s]", repo.name),
			Err:     err,
		}
	}

	// checkIfError(err, "Error during the update of the app")

	log.Println("Waiting for the update to finish ....")
	err = cmd.Wait()
	log.Println(out.String())
	log.Printf("Command finish with error: %v", err)

	if err != nil {
		return &ErrorProjectRepository{
			Code:    32,
			Message: fmt.Sprintf("Error occurred when `Cmd.Wait()` for [%s]", repo.name),
			Err:     err,
		}
	}

	return nil
}

// CopyAssets copy all the assets into the project
func (repo *ProjectRepository) CopyAssets() error {
	nBytes, err := copyFile("./assets/project.json", fmt.Sprintf("%s/assets/project.json", repo.directory))
	log.Printf("Assets copied : %d\n", nBytes)
	if err != nil {
		return &ErrorProjectRepository{
			Code:    41,
			Message: fmt.Sprintf("Error occurred when copying files for [%s]", repo.name),
			Err:     err,
		}
	}

	return nil
}

// Commit the changes done in the repository
func (repo *ProjectRepository) Commit() error {
	w, err := repo.repository.Worktree()

	if err != nil {
		return &ErrorProjectRepository{
			Code:    51,
			Message: fmt.Sprintf("Error occurred when `Worktree()` for [%s]", repo.name),
			Err:     err,
		}
	}

	_, err = w.Add(".")
	if err != nil {
		return &ErrorProjectRepository{
			Code:    52,
			Message: fmt.Sprintf("Error occurred when adding files into index `.Add('.')` for [%s]", repo.name),
			Err:     err,
		}
	}

	_, err = w.Status()
	if err != nil {
		return &ErrorProjectRepository{
			Code:    53,
			Message: fmt.Sprintf("Error occurred when getting `w.Status()` for [%s]", repo.name),
			Err:     err,
		}
	}

	commit, err := w.Commit(fmt.Sprintf("project modified committed"), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mael FOSSO",
			Email: "fosso.mael.elvis@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return &ErrorProjectRepository{
			Code:    54,
			Message: fmt.Sprintf("Error occurred when committing for [%s]", repo.name),
			Err:     err,
		}
	}

	_, err = repo.repository.CommitObject(commit)
	if err != nil {
		return &ErrorProjectRepository{
			Code:    55,
			Message: fmt.Sprintf("Error occurred when `Cmd.Wait()` for [%s]", repo.name),
			Err:     err,
		}
	}

	return nil
}

// Push the changes into the Gitlab
func (repo *ProjectRepository) Push() error {
	err := repo.repository.Push(&git.PushOptions{
		Auth: &auth,
	})
	if err != nil {
		return &ErrorProjectRepository{
			Code:    61,
			Message: fmt.Sprintf("Error occurred when pushing for [%s]", repo.name),
			Err:     err,
		}
	}

	return nil
}

func copyFile(src, dest string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return -1, &ErrorProjectRepository{
			Code:    71,
			Message: fmt.Sprintf("Error occurred. CopyFile. Impossile to open the source file [%s]", src),
			Err:     err,
		}
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return -1, &ErrorProjectRepository{
			Code:    72,
			Message: fmt.Sprintf("Error occurred. CopyFile. Impossible to create the destination file [%s]", dest),
			Err:     err,
		}
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return -1, &ErrorProjectRepository{
			Code:    73,
			Message: fmt.Sprintf("Error occurred. CopyFile. Impossible to copy file : [%s] -> [%s]", src, dest),
			Err:     err,
		}
	}

	return nBytes, nil
}

// ErrorProjectRepository collect information from ProjectRepository operation
type ErrorProjectRepository struct {
	Code    int
	Message string
	Err     error
}

func (e *ErrorProjectRepository) Error() string {
	return fmt.Sprintf("[%v] - [%s] - %v", e.Code, e.Message, e.Err)
}
