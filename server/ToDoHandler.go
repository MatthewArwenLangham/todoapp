package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/MatthewArwenLangham/todoapp/store"
	"github.com/teris-io/shortid"
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
	//Switch this to parameters???
	ToDoRe       = regexp.MustCompile(`^/api/v1/lists/*$`)
	ToDoReWithID = regexp.MustCompile(`^/api/v1/lists/[a-zA-Z0-9_-]+$`)
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
	id := strings.Split(r.URL.Path, "/api/v1/lists/")[1]
	response := h.store.GetList(id)
	jsonData, _ := json.Marshal(response)
	r.Header.Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *ToDoHandler) ViewLists(w http.ResponseWriter, r *http.Request) {
	lists := h.getAllLists()
	jsonData, _ := json.Marshal(lists)
	r.Header.Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *ToDoHandler) AddList(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("name") == "" {
		return
	}

	id, err := shortid.Generate()
	if err != nil {
		panic(err)
	}

	newList := store.List{
		Id:    id,
		Name:  r.FormValue("name"),
		Tasks: []store.Task{},
	}

	go h.store.AddList(newList)
	response := fmt.Sprintf("%s successfully added.", newList.Name)
	w.Write([]byte(response))
}

func (h *ToDoHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	newTask := store.Task{
		Id:       0,
		Name:     r.FormValue("name"),
		Complete: false,
	}

	id := strings.Split(r.URL.Path, "/api/v1/lists/")[1]
	go h.store.AddTask(id, newTask)
	response := fmt.Sprintf("%v successfully added.\n", newTask)
	w.Write([]byte(response))
}

func (h *ToDoHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	listId := strings.Split(r.URL.Path, "lists/")[1]

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

	go h.store.CompleteTask(listId, int(taskIdInt), isTaskCompleted)
}

func (h *ToDoHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	listId := strings.Split(r.URL.Path, "/api/v1/lists/")[1]
	go h.store.DeleteList(listId)
	response := fmt.Sprintf("List with ID %v successfully removed.\n", listId)
	w.Write([]byte(response))
}

// actor model
type operation struct {
	action       string
	responseChan chan map[string]store.List
}

var bufferChan = make(chan operation, 100)

func (h *ToDoHandler) getAllLists() map[string]store.List {
	responseChan := make(chan map[string]store.List, 1)
	op := operation{
		action:       "getAll",
		responseChan: responseChan,
	}
	bufferChan <- op
	lists, ok := <-responseChan
	if !ok {
		fmt.Println("Closed")
	}
	return lists
}

func (h *ToDoHandler) actor() {
	for op := range bufferChan {
		switch op.action {
		case "getAll":
			lists := h.store.GetAllLists()
			if len(lists) > 0 {
				op.responseChan <- lists
			}
		default:
			fmt.Println("You shouldn't see me... stop looking.")
		}
		close(op.responseChan)
	}
}
