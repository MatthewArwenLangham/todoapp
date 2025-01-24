package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type InMemoryStore struct {
	Lists    map[string]List
	filePath string
	mu       sync.Mutex
}

func NewMemStore() *InMemoryStore {
	list := make(map[string]List)
	return &InMemoryStore{
		Lists:    list,
		filePath: "data.json",
	}
}

func (s *InMemoryStore) AddList(list List) {
	s.Lists[list.Id] = list
}

func (s *InMemoryStore) AddTask(id string, task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lists := s.Lists[id]
	task.Id = len(lists.Tasks) + 1
	lists.Tasks = append(lists.Tasks, task)
	s.Lists[id] = lists
}

func (s *InMemoryStore) GetAllLists() map[string]List {
	return s.Lists
}

func (s *InMemoryStore) GetList(id string) List {
	return s.Lists[id]
}

func (s *InMemoryStore) CompleteTask(listId string, taskId int, isCompleted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := s.Lists[listId]
	list.Tasks[taskId-1].Complete = isCompleted
	s.Lists[listId] = list
}

func (s *InMemoryStore) DeleteList(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Lists, id)
	// go s.SaveToFile()
}

func (s *InMemoryStore) LoadFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic(err)
	}

	defer file.Close()

	var lists []List
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&lists); err != nil {
		fmt.Println("error decoding JSON:", err)
	}

	for _, list := range lists {
		s.Lists[list.Id] = list
	}
}

func (s *InMemoryStore) SaveToFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var lists []List
	for _, list := range s.Lists {
		lists = append(lists, list)
	}

	data, err := json.MarshalIndent(lists, "", " ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
