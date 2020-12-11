package main

import (
	"bytes"
	"log"
	"fmt"
	"strings"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"guitou.cm/mobile/generator/models"
)

const destCloningRepo string = "/tmp/guitou-mobile"

func main() {
	fmt.Println("Guitou mobile generator")

	// Cloning the Guitou mobile repository
	cloningRepo()

	log.Println("Load the Project JSON from Asset")
	project := models.NewProjectFromAsset()
	log.Println("Project load ", project)


	projectPackage := strings.ToLower(strings.ReplaceAll(project.Name, " ", "_"))
	
	// Download the assets related to the project
	// - icon -> ./assets/icon
	// - 

	// Run the bash script `guitou-update.sh` to update the code
	var out bytes.Buffer
	cmd := exec.Command("bash", "guitou-update.sh", projectPackage)
	cmd.Stdout = &out
	err := cmd.Start()
	checkIfError(err, "Error during the update of the app")
	log.Println("Waiting for the update to finish ....")
	err = cmd.Wait()
	log.Println(out.String())
	log.Printf("Command finish with error: %v", err)
	checkIfError(err, "Error during the update of the app")

}

// Should go into a packge related to web, http or api
func cloningRepo() {
	if _, err := os.Stat(destCloningRepo); !os.IsNotExist(err) {
		log.Println(fmt.Sprintf("Remove the destination cloning repository"))
		err := os.RemoveAll(destCloningRepo)
		checkIfError(err, fmt.Sprintf("Impossible to remove %s", destCloningRepo))
	}
	
	_, err := git.PlainClone(destCloningRepo, false, &git.CloneOptions{
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

func loadIcon() {

}

func checkIfError(err error, msg string) {
	if err != nil {
		log.Fatalln(err, msg)
	}
	
	return
}