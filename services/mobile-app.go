package services

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	repository *git.Repository
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
func (app *MobileAPP) CloneBoilerplate(projectID string) error {
	if _, err := os.Stat(projectID); !os.IsNotExist(err) {
		log.Printf("Already clone: [%s] \n", projectID)

		err := os.RemoveAll(projectID)
		if err != nil {
			return fmt.Errorf("impossible to delete [%s]", projectID)
		}

		log.Println("Last clone folder removed")
	}

	repo, err := git.PlainClone(projectID, false, &git.CloneOptions{
		URL:      GuitouURL,
		Auth:     &auth,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("error occurred when cloning for [%s]", projectID)
	}

	app.repository = repo
	log.Println("----------- [CloneBoilerplate] End : ")
	log.Println(app)

	return nil
}

// Create a new branch having the name app-{projectID}
func (app *MobileAPP) CreateBranch(projectID string) error {
	log.Println("******** [CreateBranch] ******")
	log.Println(app)

	ref, err := app.repository.Head()
	if err != nil {
		return fmt.Errorf("MAPP_CB_HEAD_ERROR")
	}
	log.Println("Getting the commit being pointed by HEAD : ", ref)

	w, err := app.repository.Worktree()
	if err != nil {
		return fmt.Errorf("MAPP_CB_WORKTREE_ERROR")
	}
	log.Println("Successful Worktree")

	err = w.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(fmt.Sprintf("app-%s", projectID)),
	})
	if err != nil {
		return fmt.Errorf("MAPP_CB_CHECKOUT_ERROR")
	}
	log.Println("Successful Checkout")

	ref, err = app.repository.Head()
	if err != nil {
		return fmt.Errorf("MAPP_CB_HEAD_END_ERROR")
	}
	fmt.Println("Get ref at end: ", ref)
	return nil
}

// Update the project folder with data from the downloaded project
func (app *MobileAPP) Update(project *protos.ProjectReply) error { // *models.Project) error {
	return nil
}

// Commit the change from Update
func (app *MobileAPP) Commit() error {
	return nil
}

// Push the new mobile application ProjectID
func (app *MobileAPP) Push() error {

	// At the end, if successful Push, delete the folder projectID
	return nil
}
