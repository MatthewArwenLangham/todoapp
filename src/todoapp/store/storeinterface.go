package store

type Store interface {
	AddList(list List)
	AddTask(id int, task Task)
	GetList(id int) List
	GetAllLists() map[int]List
	CompleteTask(listId int, taskInt int, isCompleted bool)
	DeleteList(id int)
}

type List struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id       int    `json:"id"`
	Name     string `json:"task"`
	Complete bool   `json:"complete"`
}
