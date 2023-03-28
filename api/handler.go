package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"todolist/entity"

	"github.com/gorilla/mux"
)

type UserService interface {
	AddUser(user entity.User) error
	GetUser(id int64) (entity.User, error)
	CreateSession(login, password string) (entity.Session, error)
}

type UserHandler struct {
	userService UserService
}

func NewHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// create cookie, return cookie
	type Request struct {
		login    string `json:"login"`
		password string `json:"password"`
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
	}

	session, err := h.userService.CreateSession(req.login, req.password)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   session.ID.String() + session.ExpiredAt.String(),
		Expires: time.Now().Add(time.Minute),
	}
	http.SetCookie(w, cookie)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		sendJsonError(w, err, http.StatusNotFound)
		return
	}

	sendJson(w, cookie.Value)
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req entity.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	err = h.userService.AddUser(req)
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

	sendJson(w, user)
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
