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
	router.HandleFunc("/newtask", mw.Auth(uh.AddTask))
	router.HandleFunc("/updatetask/{id}", mw.Auth(uh.UpdateTask))
	router.HandleFunc("/gettasks/{user_id}", mw.Auth(uh.GetTasks))
	router.HandleFunc("/deletetask/{id}", mw.Auth(uh.DeleteTask))

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
