package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"guitou.cm/mobile/generator/models"
	"guitou.cm/mobile/generator/services"
)

type HttpServer struct {
	store         services.IStore
	projectClient services.IProjectClient
	mobileAPP     services.IMobileAPP

	http.Handler
}

func (h HttpServer) JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHttpServer() http.Handler {
	r := mux.NewRouter()

	h := HttpServer{}
	h.store = services.NewMongoStore()
	h.projectClient = services.NewGrpcProjectClient()
	h.mobileAPP = services.NewGitlabMobileAPP()

	h.Handler = r

	r.HandleFunc("/generate/{id}", h.GenerateMobileApp)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	handler := cors.AllowAll().Handler(loggedRouter)

	return handler
}

func (h *HttpServer) GenerateMobileApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["id"]

	// Check if the project exists : gRPC - to Project micro-services - ask to project-api msvc
	var project *models.Project
	if project, err := h.projectClient.isProjectExists(pid); err != nil {
		h.JSON(w, http.StatusInternalServerError, fmt.Errorf("error when checking project existance"))
		return
	}
	if project == nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("project does not exists"))
		return
	}

	// Save locally the downloaded project
	err := h.store.SaveDownloadedProject(project)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when saving the downloaded project"))
		return
	}

	// Download the repository
	if err := h.mobileAPP.CloneBoilerplate(project.ID.String()); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when cloning boilerplate"))
		return
	}

	if err := h.mobileAPP.CreateBranch(project.ID.String()); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when updating boilerplate"))
		return
	}

	if err := h.mobileAPP.Update(project); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when updating boilerplate"))
		return
	}

	if err := h.mobileAPP.Commit(); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when committing updated"))
		return
	}

	if err := h.mobileAPP.Push(); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when pushing the mobile application"))
		return
	}

	h.JSON(w, http.StatusCreated, "link of the play store app")
}
