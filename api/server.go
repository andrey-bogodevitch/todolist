package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	http   *http.Server
}

func NewServer(port string, uh *UserHandler) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/adduser", uh.AddUser)
	router.HandleFunc("/getuser/{user_id}", uh.GetUser)
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	server := &Server{
		router: router,
		http:   httpServer,
	}
	return server
}

func (s *Server) Run() error {
	return s.http.ListenAndServe()
}
