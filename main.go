package main

import (
	"fmt"
	"log"

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

	repository := models.NewProjectRepository(project)

	log.Println("Cloning...")
	err := repository.Clone()
	checkIfError(err)
	log.Println("Clone success")

	log.Println("Checkout ...")
	err = repository.Checkout()
	checkIfError(err)
	log.Println("Checkout done...")

	log.Println("Bash Update...")
	err = repository.BashUpdate()
	checkIfError(err)
	log.Println("BashUpdate done")

	log.Println("Copying assets...")
	err = repository.CopyAssets()
	checkIfError(err)
	log.Println("CopyAssets done")

	log.Println("Commit...")
	err = repository.Commit()
	checkIfError(err)
	log.Println("Commit done")

	log.Println("Push....")
	err = repository.Push()
	checkIfError(err)
	log.Println("Push done")
}

func checkIfError(err error) {
	if err != nil {
		log.Fatalln(err)
	}

	return
}
