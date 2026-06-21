package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID       int
	Title    string
	DueDate  time.Time
	Category string
	IsDone   bool
	UserID   int
}

// global
var taskStorage []Task
var userStorage []User
var authenticatedUser *User

func main() {
	fmt.Println("Hello to TODO app")

	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

	for {
		runCommand(*command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("\n**please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}

	fmt.Println("userrStorage: %+v\n")

}

func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()
	}
	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "login":
		login()
	case "exit":
		authenticatedUser = nil
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}
}

func (u User) Print() {
	fmt.Println("User: ", u.ID, u.Email, u.Name)
}
func createTask() {
	if authenticatedUser != nil {
		authenticatedUser.Print()
	}
	scanner := bufio.NewScanner(os.Stdin)
	var title, duedate, category string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category")
	scanner.Scan()
	category = scanner.Text()

	fmt.Println("please enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()

	if authenticatedUser != nil {
		task := Task{
			ID:       len(taskStorage) + 1,
			Title:    title,
			DueDate:  duedate,
			Category: category,
			IsDone:   false,
			UserID:   authenticatedUser.ID,
		}
	}
	taskStorage = append(taskStorage, task)
}

// fmt.Println("task", name, category, duedate)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category color")
	scanner.Scan()
	color = scanner.Text()

	fmt.Println("title: ", title, "-color:", color)
}

func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var id, email, name, password string

	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	id = email

	fmt.Println("user;", id, email, password)

	user := User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}

	userStorage = append(userStorage, user)
}

func login() {
	// get user authh
	scanner := bufio.NewScanner(os.Stdin)
	var username, password string

	fmt.Println("please enter your username")
	scanner.Scan()
	username = scanner.Text()

	fmt.Println("please enter your password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if user.Email == username && user.Password == password {
			authenticatedUser = &user
			fmt.Println("You'r Logged in")

			break
		} else {
			fmt.Println("The username or password is incorrect.")
		}
	}
	if authenticatedUser == nil {
		fmt.Println("Your authentication is faild.")
	}
}
