package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	http   *http.Server
}

func NewServer(port string, uh *UserHandler, mw *AuthMiddleware) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/adduser", uh.AddUser)
	router.HandleFunc("/getuser/{user_id}", mw.Auth(uh.GetUser))
	router.HandleFunc("/login", uh.Login)
	router.HandleFunc("/delete/{user_id}", mw.Auth(uh.DeleteUser))
	router.HandleFunc("/rules/{user_id}", uh.AddAdminRules)

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
