package services

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"guitou.cm/mobile/generator/protos"
)

type IStore interface {
	SaveDownloadedProject(project *protos.ProjectReply) error // *models.Project) error
	LogAppGeneration()
	GetAllGenerations(id string)
}

type mongoStore struct {
	db mongo.Client
}

func (s mongoStore) SaveDownloadedProject(project *protos.ProjectReply) error { // *models.Project) error {
	log.Printf("Save Download Project : \n\t%v\n", project)
	return nil
}

func (s mongoStore) LogAppGeneration() {

}

func (s mongoStore) GetAllGenerations(id string) {

}

func NewMongoStore() IStore {
	return &mongoStore{}
}
