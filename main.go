package main

import (
	"log"
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const destCloningRepo string = "/tmp/guitou-mobile"

func main() {
	fmt.Println("Guitou mobile generator")

	// Cloning the Guitou mobile repository
	cloningRepo()

	// Download the Project JSON


	// Download the assets related to the project
	// - icon -> ./assets/icon
	// - 

	// Run the bash script `guitou-update.sh` to update the code

	
}

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

func checkIfError(err error, msg string) {
	if err != nil {
		log.Fatalln(err, msg)
	}
	
	return
}