package main

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"bytes"
	"log"
	"fmt"
	"strings"
	"io"
	"os"
	"os/exec"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"guitou.cm/mobile/generator/models"
)

const destCloningRepo string = "/tmp/guitou-mobile"
/**
Project: 
	- id,
	- 
	* Download JSON (gRPC)
	
Repository
	- package
	- git.PlainOpen
	-
	* Cloning
	* Working on new branch
	* Modify : Bash + Downloaded JSON File
	* Commit
	* Push
*/
func main() {
	fmt.Println("Guitou mobile generator")

	// Load Project from Assets
	log.Println("Load the Project JSON from Asset")
	project := models.NewProjectFromAsset()
	log.Println("Project load ", project)

	projectPackage := strings.ToLower(strings.ReplaceAll(project.Name, " ", "_"))
	projectDir := fmt.Sprintf("/tmp/guitou-%s", projectPackage)
	
	
	// Cloning the Guitou mobile repository
	cloningRepo(projectDir)
	log.Println("Repo cloned")
	
	checkout(projectDir, projectPackage);
	log.Println("Repo checkout")
	
	bashUpdate(projectPackage, project.Name);
	log.Println("Project bash updated")

	nBytes, err := copyFile("./assets/project.json", fmt.Sprintf("%s/assets/project.json", projectDir))
	log.Printf("Assets copied : %d\n", nBytes);
	checkIfError(err, "")
	log.Println("Asset copied")

	commitRepo(projectDir)
	log.Println("Repo committed")

	pushRepo(projectDir)
	log.Println("Repo pushed")
}

// Should go into a packge related to web, http or api
func cloningRepo(directory string) {
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		log.Println(fmt.Sprintf("Remove the destination cloning repository"))
		err := os.RemoveAll(directory)
		checkIfError(err, fmt.Sprintf("Impossible to remove %s", destCloningRepo))
	}
	
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      "https://gitlab.com/guitou-daco-tool/guitou-mobile.git",
		Auth: &http.BasicAuth{
			Username: "maelfosso",
			Password: "f170892m",
		},
		// ReferenceName: plumbing.ReferenceName("refs/heads/develop"),
		// SingleBranch: true,
    Progress: os.Stdout,
	})
	checkIfError(err, "An error occurred during git clone")
}

func checkout(directory, project string) {
	r, err := git.PlainOpen(directory)
	checkIfError(err, "")
	log.Println("[Checkout] PlainOpen")

	// headRef, err := r.Head()
	// checkIfError(err, "")

	branch := plumbing.ReferenceName(fmt.Sprintf("refs/heads/deploy/%s-%v", project, time.Now().Unix()))
	log.Println("[Checkout] ReferenceName")

	w, err := r.Worktree()
	checkIfError(err, "")
	log.Println("[Checkout] Worktree")

	err = w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force: false,
		Branch: branch,
	})
	if err != nil {
		log.Println("Checkout Failed")
		err = w.Checkout(&git.CheckoutOptions{
			Create: true,
			Force: false,
			Branch: branch,
		})
		checkIfError(err, "")
		
	}
	log.Println("[Checkout] Checkout Final")
}

func bashUpdate(projectPackage, projectName string) {
	var out bytes.Buffer
	cmd := exec.Command("bash", "guitou-update.sh", projectPackage, projectName)
	cmd.Stdout = &out
	err := cmd.Start()
	checkIfError(err, "Error during the update of the app")
	log.Println("Waiting for the update to finish ....")
	err = cmd.Wait()
	log.Println(out.String())
	log.Printf("Command finish with error: %v", err)
	checkIfError(err, "Error during the update of the app")
}

func commitRepo(directory string) {
	r, err := git.PlainOpen(directory)
	checkIfError(err, "")

	w, err := r.Worktree()
	checkIfError(err, "")

	_, err = w.Add(".")
	checkIfError(err, "")

	status, err := w.Status()
	checkIfError(err, "")

	log.Println(status)

	commit, err := w.Commit(fmt.Sprintf("project modified committed"), &git.CommitOptions{
		Author: &object.Signature{
			Name: "Mael FOSSO",
			Email: "fosso.mael.elvis@gmail.com",
			When: time.Now(),
		},
	})
	checkIfError(err, "")

	obj, err := r.CommitObject(commit)
	checkIfError(err, "")

	fmt.Println(obj)
}

func pushRepo(dir string) {
	r, err := git.PlainOpen(dir)
	checkIfError(err, "Git Plain Open")

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "maelfosso",
			Password: "f170892m",
		},
	})
	checkIfError(err, "")
}

func checkIfError(err error, msg string) {
	if err != nil {
		log.Fatalln(err, msg)
	}
	
	return
}

func copyFile(src, dest string) (int64, error) {
	source, err := os.Open(src)
	checkIfError(err, "Impossible to open the source file")
	defer source.Close()

	destination, err := os.Create(dest)
	checkIfError(err, "Impossible to create the destination file")
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	checkIfError(err, "Impossible to copy file")

	return nBytes, nil
}