package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"unicode"

	"todolist/entity"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserService interface {
	AddUser(user entity.User) error
	GetUser(id int64) (entity.User, error)
	CreateSession(login, password string) (entity.Session, error)
	FindSessionByID(id uuid.UUID) (entity.Session, error)
	DeleteUser(id int64) error
	AddAdminRules(id int64) error
	AddTask(task entity.Task) error
	UpdateTask(task entity.Task) error
	GetTasks(id int64) ([]entity.Task, error)
	GetTaskByID(id int64) (entity.Task, error)
	DeleteTask(taskID int64) error
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
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
	}

	session, err := h.userService.CreateSession(req.Login, req.Password)
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   session.ID.String(),
		Expires: session.ExpiredAt,
	}
	http.SetCookie(w, cookie)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	//получаем userid из запроса
	userID := mux.Vars(r)
	id := userID["user_id"]
	userIDInt, err := strconv.Atoi(id)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	value := ctx.Value(ctxUserKey{})

	user := value.(entity.User)

	//проверить, что userid из запроса равен userid из пользователя
	if int64(userIDInt) != user.ID && user.Role != "admin" {
		sendJsonError(w, fmt.Errorf("you can't delete other users"), http.StatusForbidden)
		return
	}

	//вызываем метод deleteUser
	err = h.userService.DeleteUser(int64(userIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) AddAdminRules(w http.ResponseWriter, r *http.Request) {
	//получаем userid из запроса
	userID := mux.Vars(r)
	id := userID["user_id"]
	userIDInt, err := strconv.Atoi(id)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	err = h.userService.AddAdminRules(int64(userIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req entity.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	if isValidName(req.Name) != true {
		sendJsonError(w, fmt.Errorf("invalid name"), http.StatusBadRequest)
		return

	}

	if isValidLogin(req.Login) != true {
		sendJsonError(w, fmt.Errorf("invalid login"), http.StatusForbidden)
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

	user := userFromCtx(r)

	//проверить, что userid из запроса равен userid из пользователя
	if int64(userIDInt) != user.ID && user.Role != "admin" {
		sendJsonError(w, fmt.Errorf("you can't get other users"), http.StatusForbidden)
		return
	}

	user, err = h.userService.GetUser(int64(userIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
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
func userFromCtx(r *http.Request) entity.User {
	ctx := r.Context()
	value := ctx.Value(ctxUserKey{})
	user := value.(entity.User)
	return user
}

func (h *UserHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var req entity.Task
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	user := userFromCtx(r)

	req.UserID = user.ID

	//проверить, что userid из запроса равен userid из пользователя
	if req.UserID != user.ID && user.Role != "admin" {
		sendJsonError(w, fmt.Errorf("you can add only your tasks"), http.StatusForbidden)
		return
	}

	err = h.userService.AddTask(req)
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)
	id := taskID["id"]
	taskIDInt, err := strconv.Atoi(id)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	user := userFromCtx(r)

	task, err := h.userService.GetTaskByID(int64(taskIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusNotFound)
		return
	}
	//проверить, что userid из запроса равен userid из пользователя
	if task.UserID != user.ID && user.Role != "admin" {
		sendJsonError(w, fmt.Errorf("you can update only your tasks"), http.StatusForbidden)
		return
	}

	err = h.userService.UpdateTask(task)
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)
	id := userID["user_id"]
	userIDInt, err := strconv.Atoi(id)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	user := userFromCtx(r)

	//проверить, что userid из запроса равен userid из пользователя
	if int64(userIDInt) != user.ID && user.Role != "admin" {
		sendJsonError(w, fmt.Errorf("you can't get other users tasks"), http.StatusForbidden)
		return
	}

	tasks, err := h.userService.GetTasks(int64(userIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}

	sendJson(w, tasks)
}

func (h *UserHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)
	id := taskID["id"]
	taskIDInt, err := strconv.Atoi(id)
	if err != nil {
		sendJsonError(w, err, http.StatusBadRequest)
		return
	}

	user := userFromCtx(r)

	task, err := h.userService.GetTaskByID(int64(taskIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusNotFound)
		return
	}

	//проверить, что userid из запроса равен userid из пользователя
	if task.UserID != user.ID && user.Role != "admin" {
		sendJsonError(w, fmt.Errorf("you can't delete other users tasks"), http.StatusForbidden)
		return
	}

	err = h.userService.DeleteTask(int64(taskIDInt))
	if err != nil {
		sendJsonError(w, err, http.StatusInternalServerError)
		return
	}

}

func isValidName(str string) bool {
	if len(str) < 3 || len(str) > 40 {
		return false
	}

	for _, r := range str {
		if !unicode.Is(unicode.Latin, r) {
			return false
		}
	}

	return true
}

func isValidLogin(str string) bool {
	const ValidSymbols = `^[a-zA-Z0-9._-]{3,15}$`
	var IsLetter = regexp.MustCompile(ValidSymbols).MatchString
	return IsLetter(str)
}
