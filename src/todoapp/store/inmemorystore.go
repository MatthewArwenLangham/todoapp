package store

import (
	"encoding/json"
	"fmt"
	"os"
)

type InMemoryStore struct {
	Lists    map[int]List
	filePath string
}

func NewMemStore() *InMemoryStore {
	list := make(map[int]List)
	return &InMemoryStore{
		Lists:    list,
		filePath: "data.json",
	}
}

func (s *InMemoryStore) AddList(list List) {
	fmt.Println(list.Id)
	s.Lists[list.Id] = list
	fmt.Println(s.Lists)
}

func (s *InMemoryStore) AddTask(id int, task Task) {
	lists := s.Lists[id]
	task.Id = len(lists.Tasks) + 1
	lists.Tasks = append(lists.Tasks, task)
	s.Lists[id] = lists
}

func (s *InMemoryStore) GetAllLists() map[int]List {
	return s.Lists
}

func (s *InMemoryStore) GetList(id int) List {
	return s.Lists[id]
}

func (s *InMemoryStore) CompleteTask(listId int, taskId int, isCompleted bool) {
	list := s.Lists[listId]
	list.Tasks[taskId-1].Complete = isCompleted
	s.Lists[listId] = list
}

func (s *InMemoryStore) DeleteList(id int) {
	delete(s.Lists, id)
}

func (s *InMemoryStore) LoadFromFile() {
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
	file, err := os.Create(s.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, _ := json.MarshalIndent(s.Lists, "", "  ")
	os.WriteFile(s.filePath, data, 0644)
}
