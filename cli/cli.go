package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Bug to fix with delete...

type db struct {
	lists []list `json:"lists"`
}

type list struct {
	id    int    `json:"id"`
	name  string `json:"name"`
	tasks []Task `json:"tasks"`
}

type Task struct {
	id       int    `json:"id"`
	name     string `json:"task"`
	complete bool   `json:"complete"`
}

var nextListId = 0
var nextTaskId = 0

func Run() {
	reader := bufio.NewReader(os.Stdin)
	db := db{}
	fmt.Println("To Do list app")
	fmt.Println("--------------")

	for {
		fmt.Println("\nCommands:")
		fmt.Println("1. Create list")
		fmt.Println("2. View lists")
		fmt.Println("x. Exit")
		fmt.Println()
		fmt.Print("Choose an option: ")
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			createList(reader, &db)
		case "2":
			viewLists(reader, &db)
		case "x":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func createList(reader *bufio.Reader, db *db) {
	fmt.Println()
	fmt.Println("Create a list")
	fmt.Println("-------------")
	fmt.Println()
	fmt.Println("Name your list: ")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	newList := list{id: nextListId, name: input}
	nextListId++
	db.lists = append(db.lists, newList)
}

func viewLists(reader *bufio.Reader, db *db) {
	for {
		fmt.Println()
		fmt.Println("Lists:")
		fmt.Println("------")
		if len(db.lists) == 0 {
			fmt.Println("No lists found!")
		} else {
			for i, list := range db.lists {
				fmt.Println(i+1, list.name)
			}
		}
		fmt.Println("----------------")
		fmt.Println()
		fmt.Println("x Back to main menu")
		fmt.Println()
		fmt.Println("Select a list to view/edit: ")
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "x" {
			return
		}

		intSelection, err := strconv.ParseInt(input, 10, 0)
		if err != nil {
			fmt.Println("Please enter an integer.")
		}

		if intSelection > 0 && intSelection <= int64(len(db.lists)) {
			viewList(reader, &db.lists[intSelection-1])
		} else {
			fmt.Println(intSelection, "is not a valid number.")
		}
	}
}

func viewList(reader *bufio.Reader, list *list) {
	for {
		fmt.Println()
		fmt.Println(list.name)
		fmt.Println("----------------")
		for i, task := range list.tasks {
			fmt.Println(i+1, task.name, task.complete)
		}

		fmt.Println()
		fmt.Println("----------------")
		fmt.Println("1. Add a task")
		fmt.Println("2. Complete a task")
		fmt.Println("3. Delete a task")
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			addTask(reader, list)
		case "2":
			completeTask(reader, list)
		case "3":
			deleteTask(reader, list)
		case "x":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func addTask(reader *bufio.Reader, list *list) {
	fmt.Println()
	fmt.Println("Enter task name or press type x to go back")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "x" {
		return
	}

	taskToAdd := Task{name: input, complete: false}

	list.tasks = append(list.tasks, taskToAdd)

	fmt.Println("Task \"", input, "\" added successfully!")
}

func completeTask(reader *bufio.Reader, list *list) {
	fmt.Println()
	for i, task := range list.tasks {
		fmt.Println(i+1, task.name, task.complete)
	}
	fmt.Println()
	fmt.Println("Enter task id to complete or press type x to go back")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "x" {
		return
	}

	intSelection, err := strconv.ParseInt(input, 10, 0)
	if err != nil {
		fmt.Println("Please enter an integer.")
	}

	if intSelection > 0 && intSelection <= int64(len(list.tasks)) {
		list.tasks[intSelection-1].complete = true
	} else {
		fmt.Println(intSelection, "is not a valid number.")
	}
	fmt.Println("Task \"", list.tasks[intSelection-1], "\" completed!")
}

func deleteTask(reader *bufio.Reader, list *list) {
	fmt.Println()
	for i, task := range list.tasks {
		fmt.Println(i+1, task.name, task.complete)
	}
	fmt.Println()
	fmt.Println("Enter task id to delete or press type x to go back")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "x" {
		return
	}

	intSelection, err := strconv.ParseInt(input, 10, 0)
	if err != nil {
		fmt.Println("Please enter an integer.")
	}

	if intSelection > 0 && intSelection <= int64(len(list.tasks)) {
		if len(list.tasks) > 1 {
			list.tasks = append(list.tasks[:intSelection-1], list.tasks[intSelection:]...)
		} else {
			list.tasks = nil
		}
	} else {
		fmt.Println(intSelection, "is not a valid number.")
	}
	fmt.Println("Task \"", list.tasks[intSelection-1], "\" completed!")
}
