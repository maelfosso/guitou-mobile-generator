package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Project keep information
// We must add :
//	- Organisation NAME  as a Title . That title will be used in the package name
// Then the project name can be the subtitle
type Project struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewProjectFromAsset load a project from the assets
func NewProjectFromAsset() *Project {
	jsonFile, err := os.Open("assets/project.json")
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
	}
	log.Println("Asset JSON opended")

	var project Project = Project{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &project)

	return &project
}

// NewProjectFromID get project from its ID
func NewProjectFromID(id string) *Project {
	return nil
}
