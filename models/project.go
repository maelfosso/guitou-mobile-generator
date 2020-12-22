package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project keep information
// We must add :
//	- Organisation NAME  as a Title . That title will be used in the package name
// Then the project name can be the subtitle
type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
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

// ProjectRepository provides access to a Project store
type ProjectRepository interface {
	Save(project *Project) error
	Find(id string) (*Project, error)
	FindAll() []*Project
}

// ErrUnknown is used when a project could not be found
var ErrUnknown = errors.New("unknown project")
