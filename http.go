package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"guitou.cm/mobile/generator/protos"
	"guitou.cm/mobile/generator/services"
)

type HttpServer struct {
	store         services.IStore
	projectClient protos.ProjectsClient // services.IProjectClient
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
	log.Println("NewHttpServer()")

	log.Println("Connecting to gRPC Project")
	projectClient, _, err := services.NewGrpcProjectClient()
	if err != nil {
		log.Fatalf("not possible to connect to Project-API : %v", err)
	}
	log.Println("Successful connection to gRPC Project Server")
	// defer closeProjectClient()

	h := HttpServer{}
	h.store = services.NewMongoStore()
	h.mobileAPP = services.NewGitlabMobileAPP()
	h.projectClient = projectClient
	// h.Handler = r

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.NotFoundHandler = r.NewRoute().BuildOnly().HandlerFunc(http.NotFound).GetHandler()

	r.HandleFunc("/{id}", h.GenerateMobileApp).Methods(http.MethodGet)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	handler := cors.AllowAll().Handler(loggedRouter)

	h.Handler = handler
	return h
}

func (h *HttpServer) GenerateMobileApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["id"]

	// Check if the project exists : gRPC - to Project micro-services - ask to project-api msvc
	// var project *models.Project
	// if project, err := h.projectClient.IsProjectExists() isProjectExists(pid); err != nil {
	var project *protos.ProjectReply
	idRequest := &protos.IDRequest{
		Id: pid,
	}
	project, err := h.projectClient.IsProjectExists(r.Context(), idRequest)
	if err != nil {
		h.JSON(w, http.StatusInternalServerError, fmt.Errorf("error when checking project existance"))
		return
	}
	if project == nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("project does not exists"))
		return
	}
	log.Println("Grpc downloaded Project : \n\t", project)

	// Save locally the downloaded project
	if err := h.store.SaveDownloadedProject(project); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when saving the downloaded project"))
		return
	}

	// Download the repository
	if err := h.mobileAPP.CloneBoilerplate(project.Id); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when cloning boilerplate"))
		return
	}
	log.Println("Project successfully clone ", h.mobileAPP)

	if err := h.mobileAPP.CreateBranch(project.Id); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when updating boilerplate"))
		return
	}
	log.Println("Successful Checkout on new branch")

	if err := h.mobileAPP.Update(project); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when updating boilerplate"))
		return
	}
	log.Println("TODO - Update the project not done yet")

	if err := h.mobileAPP.Commit(project); err != nil {
		h.JSON(w, http.StatusBadRequest, fmt.Errorf("error when committing updated"))
		return
	}
	log.Println("TODO - Committing the update after implementing the Update")

	if err := h.mobileAPP.Push(); err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}
	log.Println("Successful Push")

	h.JSON(w, http.StatusCreated, "link of the play store app")
}
