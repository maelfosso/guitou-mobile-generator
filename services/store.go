package services

import (
	"go.mongodb.org/mongo-driver/mongo"
	"guitou.cm/mobile/generator/models"
)

type IStore interface {
	SaveDownloadedProject(project *models.Project) error
	LogAppGeneration()
	GetAllGenerations(id string)
}

type mongoStore struct {
	db mongo.Client
}

func (s mongoStore) SaveDownloadedProject(project *models.Project) error {
	return nil
}

func (s mongoStore) LogAppGeneration() {

}

func (s mongoStore) GetAllGenerations(id string) {

}

func NewMongoStore() IStore {
	return &mongoStore{}
}
