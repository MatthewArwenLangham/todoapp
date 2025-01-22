package store

import (
	"os"
	"testing"
)

var store *InMemoryStore
var testFilePath = "test_data.json"

func setup() {
	store = &InMemoryStore{
		Lists:    make(map[string]List),
		filePath: testFilePath,
	}
	list := List{Id: "1", Name: "Groceries", Tasks: []Task{
		{Id: 1, Name: "Buy milk", Complete: false},
		{Id: 2, Name: "Buy bread", Complete: false},
	}}
	store.AddList(list)
}

func teardown() {
	store = nil
	os.Remove(testFilePath)
}

func TestAddList(t *testing.T) {
	setup()
	defer teardown()

	newList := List{Id: "2", Name: "Chores"}
	store.AddList(newList)

	if len(store.Lists) != 2 {
		t.Errorf("expected 2 lists, got %d", len(store.Lists))
	}

	if store.Lists["2"].Name != "Chores" {
		t.Errorf("expected list name 'Chores', got '%s'", store.Lists["2"].Name)
	}
}

func TestAddTask(t *testing.T) {
	setup()
	defer teardown()

	task := Task{Name: "Buy eggs"}
	store.AddTask("1", task)

	if len(store.Lists["1"].Tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(store.Lists["1"].Tasks))
	}

	if store.Lists["1"].Tasks[2].Name != "Buy eggs" {
		t.Errorf("expected task name 'Buy eggs', got '%s'", store.Lists["1"].Tasks[2].Name)
	}
}

func TestCompleteTask(t *testing.T) {
	setup()
	defer teardown()

	store.CompleteTask("1", 1, true)

	if !store.Lists["1"].Tasks[0].Complete {
		t.Errorf("expected task to be complete, but it was not")
	}
}

func TestDeleteList(t *testing.T) {
	setup()
	defer teardown()

	store.DeleteList("1")

	if len(store.Lists) != 0 {
		t.Errorf("expected 0 lists, got %d", len(store.Lists))
	}
}

func TestLoadFromFile(t *testing.T) {
	setup()
	defer teardown()

	data := `[{"id": "2", "name": "Chores", "tasks": [{"id": 1, "task": "Do laundry", "complete": false}]}]`
	os.WriteFile(testFilePath, []byte(data), 0644)

	store.LoadFromFile()

	if len(store.Lists) != 2 {
		t.Errorf("expected 2 lists, got %d", len(store.Lists))
	}

	if store.Lists["2"].Name != "Chores" {
		t.Errorf("expected list name 'Chores', got '%s'", store.Lists["2"].Name)
	}
}

func TestSaveToFile(t *testing.T) {
	setup()
	defer teardown()

	store.SaveToFile()

	data, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	expected := `{
  "1": {
    "id": "1",
    "name": "Groceries",
    "tasks": [
      {
        "id": 1,
        "task": "Buy milk",
        "complete": false
      },
      {
        "id": 2,
        "task": "Buy bread",
        "complete": false
      }
    ]
  }
}`
	if string(data) != expected {
		t.Errorf("expected file content %s, got %s", expected, string(data))
	}
}
