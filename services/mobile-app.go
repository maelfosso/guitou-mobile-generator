package services

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"

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
	MobileAppDir = "mobile-app-boilerplate"
)

var (
	GuitouURL = os.Getenv("GIT_REPO_URL")
	Username  = os.Getenv("GIT_AUTH_USERNAME")
	Password  = os.Getenv("GIT_AUTH_PASSWORD")
	auth      = http.BasicAuth{
		Username: Username,
		Password: Password,
	}
)

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
	log.Println("Update start...")

	var wg sync.WaitGroup

	files, err := WalkMatch(".", "*.tmpl")
	if err != nil {
		log.Println("error when WalkMatch, ", err)
		return fmt.Errorf("MAPP_UPDATE_WALKMATCH_ERROR")
	}
	log.Println("WalkMatch : ", files)

	funcMap := template.FuncMap{
		"toId": func(str string) string {
			value := strings.ToLower(str)

			if _, err := mail.ParseAddress(value); err == nil {
				return strings.Split(value, "@")[0]
			} else {
				reg, err := regexp.Compile("[^a-zA-Z0-9]+")
				if err != nil {
					log.Fatal(err)
				}
				return reg.ReplaceAllString(value, "_")
			}
		},
	}

	wg.Add(len(files))

	for _, file := range files {

		go func(file string) {
			defer wg.Done()

			filename := filepath.Base(file)
			log.Printf("\n\n************* %s\n", filename)

			// Create a new file without the extension .tmpl
			newFilename := strings.TrimSuffix(filename, filepath.Ext((filename)))
			newFile, err := os.Create(filepath.Join(filepath.Dir(file), newFilename))
			if err != nil {
				log.Panic("error when creating the file, ", err)
			}

			// Run the template
			t := template.Must(template.New(filename).Funcs(funcMap).ParseFiles(file))
			if err != nil {
				log.Panic("error occured", err)
				// return fmt.Errorf("MAPP_UPDATE_PARSEGLOB_ERROR")
			}

			err = t.Execute(newFile, project)
			if err != nil {
				log.Println("error occured when executing, ", err)
				return
			}

			// Delete the template file
			err = os.Remove(file)
			if err != nil {
				log.Panic("error when deleting, ", err)
			}
		}(file)
	}

	wg.Wait()
	// t := template.Must(
	// 	template.New("pubspec.yaml.tmpl").Funcs(template.FuncMap{
	// 		"toId": func(str string) string {
	// 			value := strings.ToLower(str)

	// 			if _, err := mail.ParseAddress(value); err == nil {
	// 				return strings.Split(value, "@")[0]
	// 			} else {
	// 				reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	// 				if err != nil {
	// 					log.Fatal(err)
	// 				}
	// 				return reg.ReplaceAllString(value, "_")
	// 			}
	// 		},
	// 	}).ParseFiles(files...))
	// if err != nil {
	// 	log.Println("error occured", err)
	// 	return fmt.Errorf("MAPP_UPDATE_PARSEGLOB_ERROR")
	// }
	// fmt.Println("ParseGlob result: ", t)
	// // fmt.Println(t.Root.Nodes)
	// // fmt.Println(t.ParseName)
	// fmt.Println(t.Templates(), len(t.Templates()))
	// fmt.Println(t.Templates()[0].Name, t.Templates()[0].ParseName)

	// err = t.Execute(os.Stdout, project)
	// // err = t.ExecuteTemplate(os.Stdout, project.Id, project)
	// if err != nil {
	// 	log.Println("error occured when executing, ", err)
	// 	return fmt.Errorf("MAPP_UPDATE_EXECUTE_ERROR")
	// }

	return nil
}

// Commit the change from Update
func (app *MobileAPP) Commit() error {
	return nil
}

// Push the new mobile application ProjectID
func (app *MobileAPP) Push() error {
	log.Println("Push start ..................")

	err := app.repository.Push(&git.PushOptions{
		Auth:       &auth,
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("MAPP_PUSH_PUSH_ERROR")
	}

	// At the end, if successful Push, delete the folder projectID
	return nil
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
