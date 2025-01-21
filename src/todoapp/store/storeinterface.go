package store

type Store interface {
	AddList(list List)
	AddTask(id string, task Task)
	GetList(id string) List
	GetAllLists() map[string]List
	CompleteTask(listId string, taskInt int, isCompleted bool)
	DeleteList(id string)
}

type List struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id       int    `json:"id"`
	Name     string `json:"task"`
	Complete bool   `json:"complete"`
}
