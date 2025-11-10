package server

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/ReanSn0w/go-docker-status/internal/docker"
	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/lgr"
)

func New(log lgr.L) (*Server, error) {
	docker, err := docker.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		log:    log,
		srv:    web.New(log),
		tmpl:   template.Must(template.ParseGlob("./tmpl/*.html")),
		docker: docker,
	}, nil
}

type Server struct {
	log    lgr.L
	srv    *web.Server
	tmpl   *template.Template
	docker *docker.Docker
}

func (s *Server) Start(cancel context.CancelCauseFunc, port int) {
	s.log.Logf("[INFO] Starting server on port %d", port)
	s.srv.Run(cancel, port, s.handler())
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Logf("[ERROR] Stopping server err: %v", err)
	}
}

func (s *Server) handler() http.Handler {
	router := chi.NewRouter()

	router.Get("/", s.statusPageHandler)

	return router
}

type pageData struct {
	Containers []docker.ContainerInfo
	TotalCount int
	Error      string
}

func (s *Server) statusPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := pageData{
		Containers: []docker.ContainerInfo{},
	}

	containers, err := s.docker.GetContainers(ctx)
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Containers = containers
		data.TotalCount = len(containers)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, data); err != nil {
		log.Printf("template execution error: %v", err)
	}
}
