package services

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"guitou.cm/mobile/generator/protos"
)

type IMobileAPP interface {
	CloneBoilerplate(projectID string) error
	CreateBranch(projectID string) error
	Update(project *protos.ProjectReply) error // *models.Project) error
	Commit() error
	Push() error
}

type MobileAPP struct {
}

const (
	GuitouURL    = "https://gitlab.com/guitou-app/mobile-app-boilerplate"
	MobileAppDir = "mobile-app-boilerplate"
	Username     = "maelfosso" // From K8s ENV
	Password     = "f170892m"  // From K8s ENV
)

var auth = http.BasicAuth{
	Username: Username,
	Password: Password,
}

func NewGitlabMobileAPP() IMobileAPP {
	return &MobileAPP{}
}

// Clone the boilerplate of the mobile application
// Renamed the clone repository with projectID
func (r MobileAPP) CloneBoilerplate(projectID string) error {
	if _, err := os.Stat(MobileAppDir); !os.IsNotExist(err) {
		err := os.RemoveAll(MobileAppDir)
		if err != nil {
			return fmt.Errorf("impossible to delete [%s]", MobileAppDir)
		}
	}

	_, err := git.PlainClone(MobileAppDir, false, &git.CloneOptions{
		URL:      GuitouURL,
		Auth:     &auth,
		Progress: os.Stdout,
	})

	// MobileAPP cloned

	if err != nil {
		return fmt.Errorf("Error occurred when cloning [%s] for [%s]")
	}

	return nil
}

// Create a new branch having the name app-{projectID}
func (r MobileAPP) CreateBranch(projectID string) error {
	return nil
}

// Update the project folder with data from the downloaded project
func (r MobileAPP) Update(project *protos.ProjectReply) error { // *models.Project) error {
	return nil
}

// Commit the change from Update
func (r MobileAPP) Commit() error {
	return nil
}

// Push the new mobile application ProjectID
func (r MobileAPP) Push() error {

	// At the end, if successful Push, delete the folder projectID
	return nil
}
