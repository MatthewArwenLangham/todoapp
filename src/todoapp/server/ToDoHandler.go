package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/MatthewArwenLangham/todoapp/store"
)

type ToDoHandler struct {
	store store.Store
}

func NewToDoHandler(s store.Store) *ToDoHandler {
	return &ToDoHandler{
		store: s,
	}
}

var (
	ToDoRe       = regexp.MustCompile(`^/lists/*$`)
	ToDoReWithID = regexp.MustCompile(`^/lists/[0-9]+$`)
)

func (h *ToDoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && ToDoReWithID.MatchString(r.URL.Path):
		h.ViewList(w, r)
		return

	case r.Method == http.MethodGet && ToDoRe.MatchString(r.URL.Path):
		h.ViewLists(w, r)
		return

	case r.Method == http.MethodPost && ToDoRe.MatchString(r.URL.Path):
		h.AddList(w, r)
		return

	case r.Method == http.MethodPost && ToDoReWithID.MatchString(r.URL.Path):
		h.AddTask(w, r)
		return

	case r.Method == http.MethodPatch && ToDoReWithID.MatchString(r.URL.Path):
		h.CompleteTask(w, r)
		return

	case r.Method == http.MethodDelete && ToDoReWithID.MatchString(r.URL.Path):
		h.DeleteList(w, r)
		return
	}
}

func (h *ToDoHandler) ViewList(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "lists/")[1]
	idInt, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		w.Write([]byte("Invalid id"))
	}

	response := fmt.Sprintf("%v\n", h.store.GetList(int(idInt)))
	w.Write([]byte(response))
}

func (h *ToDoHandler) ViewLists(w http.ResponseWriter, r *http.Request) {
	lists := h.store.GetAllLists()
	response := fmt.Sprintf("%v\n", lists)
	fmt.Println(response)
	w.Write([]byte(response))
}

func (h *ToDoHandler) AddList(w http.ResponseWriter, r *http.Request) {
	newList := store.List{
		Id:    len(h.store.GetAllLists()) + 1,
		Name:  r.FormValue("name"),
		Tasks: []store.Task{},
	}

	h.store.AddList(newList)
	response := fmt.Sprintf("%s successfully added.", newList.Name)
	w.Write([]byte(response))
}

func (h *ToDoHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	newTask := store.Task{
		Id:       0,
		Name:     r.FormValue("name"),
		Complete: false,
	}

	id := strings.Split(r.URL.Path, "lists/")[1]
	idInt, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		panic(err)
	}

	h.store.AddTask(int(idInt), newTask)
	response := fmt.Sprintf("%v successfully added.\n", newTask)
	w.Write([]byte(response))
}

func (h *ToDoHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	listId := strings.Split(r.URL.Path, "lists/")[1]
	listIdInt, err := strconv.ParseInt(listId, 10, 0)
	if err != nil {
		w.Write([]byte("Invalid id"))
	}

	taskId := r.FormValue("taskId")
	taskIdInt, err := strconv.ParseInt(taskId, 10, 0)
	if err != nil {
		w.Write([]byte("Invalid id"))
	}

	completed := r.FormValue("completed")
	isTaskCompleted, err := strconv.ParseBool(completed)
	if err != nil {
		panic(err)
	}

	h.store.CompleteTask(int(listIdInt), int(taskIdInt), isTaskCompleted)
}

func (h *ToDoHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	listId := strings.Split(r.URL.Path, "lists/")[1]
	listIdInt, err := strconv.ParseInt(listId, 10, 0)
	if err != nil {
		w.Write([]byte("Invalid id"))
	}

	h.store.DeleteList(int(listIdInt))
	response := fmt.Sprintf("List with ID %v successfully removed.\n", listIdInt)
	w.Write([]byte(response))
}
