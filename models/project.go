package models

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
)

// Project keep information 
type Project struct {
	ID string `json:"_id"`
	Name string `json:"name"`
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
