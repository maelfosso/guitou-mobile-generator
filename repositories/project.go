package repositories

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"guitou.cm/mobile/generator/db"
	"guitou.cm/mobile/generator/models"
)

// const collection = "projects"

type projectRepository struct {
	// dbConnexion *db.MongoDBClient
	collection *mongo.Collection
}

// NewProjectRepository generates a new repository for projects
func NewProjectRepository(dbConnexion *db.MongoDBClient) models.ProjectRepository {
	collection := dbConnexion.GetCollection("projects")

	return &projectRepository{
		collection,
	}
}

func (pr *projectRepository) Save(p *models.Project) error {

	_, err := pr.collection.InsertOne(context.TODO(), p)
	if err != nil {
		return err
	}

	return nil
}

func (pr *projectRepository) Find(id string) (*models.Project, error) {
	var project *models.Project

	err := pr.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pr *projectRepository) FindAll() []*models.Project {
	ctx := context.TODO()

	cursor, err := pr.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var projects []*models.Project
	if err = cursor.All(ctx, &projects); err != nil {
		log.Fatal(err)
		return []*models.Project{}
	}

	return projects
}
