package services

// // Service is the interface to provide all API Service
// type Service interface {
// 	Generate(projectID string) (bool, error)
// }

// type service struct {
// 	projectRepository models.ProjectRepository
// }

// // NewService generate service object
// func NewService(projectRepository models.ProjectRepository) Service {
// 	return &service{
// 		projectRepository,
// 	}
// }

// // ErrorNoProjectID is returned when there is no ProjecID parameter
// var ErrorNoProjectID = errors.New("No ProjectID provided")

// func (s *service) Generate(projectID string) (bool, error) {
// 	if projectID == "" {
// 		return false, ErrorNoProjectID
// 	}

// 	// project := models.NewProjectFromID(projectID)
// 	project := models.NewProjectFromAsset()
// 	s.projectRepository.Save(project)

// 	repository := models.NewProjectOnGit(project)

// 	log.Println("Cloning...")
// 	err := repository.Clone()
// 	if err != nil {
// 		return false, err
// 	}
// 	log.Println("Clone success")

// 	log.Println("Checkout ...")
// 	err = repository.Checkout()
// 	if err != nil {
// 		return false, err
// 	}
// 	log.Println("Checkout done...")

// 	log.Println("Bash Update...")
// 	err = repository.BashUpdate()
// 	if err != nil {
// 		return false, err
// 	}
// 	log.Println("BashUpdate done")

// 	log.Println("Copying assets...")
// 	err = repository.CopyAssets()
// 	if err != nil {
// 		return false, err
// 	}
// 	log.Println("CopyAssets done")

// 	log.Println("Commit...")
// 	err = repository.Commit()
// 	if err != nil {
// 		return false, err
// 	}
// 	log.Println("Commit done")

// 	log.Println("Push....")
// 	err = repository.Push()
// 	if err != nil {
// 		return false, err
// 	}
// 	log.Println("Push done")

// 	return true, nil
// }
