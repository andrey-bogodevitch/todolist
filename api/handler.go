package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"todolist/storage"

	"github.com/gorilla/mux"
)

type UserService interface {
	AddUser(name, login, password string) error
	GetUser(id int64) (storage.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Name     string `json:"name"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	err = h.userService.AddUser(req.Name, req.Login, req.Password)
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)
	id := userID["user_id"]
	userIDInt, err := strconv.Atoi(id)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(int64(userIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
	}

	type Response struct {
		ID        int64
		Name      string
		Role      string
		CreatedAt time.Time
		Login     string
	}

	resp := Response{
		ID:        user.ID,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		Login:     user.Login,
	}
	sendJson(w, resp)
}

type jsonError struct {
	Error string `json:"error"`
}

func sendJsonError(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	sendJson(w, jsonError{Error: err.Error()}, code)
}

func sendJson(w http.ResponseWriter, data any, code ...int) {
	w.Header().Set("Content-Type", "application/json")

	if len(code) > 0 {
		w.WriteHeader(code[0])
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
	}
}
