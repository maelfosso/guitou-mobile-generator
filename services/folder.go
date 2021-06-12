package services

type IFolder interface {
	SaveDownloadedJSON()
	CopyBoilerplate()
	Update()
}

type Folder struct {
}

func NewFolder() IFolder {
	return &Folder{}
}

func (f Folder) SaveDownloadedJSON() {

}

func (f Folder) CopyBoilerplate() {

}

func (f Folder) Update() {

}
