package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	//"time"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID    int
	Title string
	//DueDate  time.Time
	DueDate    string
	CategoryID int
	IsDone     bool
	UserID     int
}

type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

// global
var taskStorage []Task
var userStorage []User
var categoryStorage []Category
var authenticatedUser *User

const userStoragePath string = "users.txt"

func main() {
	// load data
	leadUserStorageFromFile()
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

		if authenticatedUser == nil {
			return
		}
	}
	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "list-task":
		listTask()
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

func leadUserStorageFromFile() {
	file, err := os.Open(userStoragePath)
	if err != nil {
		fmt.Println("can't open the file.", err)
	} else {
		var data = make([]byte, 1024)
		_, oErr := file.Read(data)
		if oErr != nil {
			fmt.Println("can't read from the file", oErr)
		} else {
			dataStr := string(data)
			dataStr = strings.Trim(dataStr, "\n")
			dataList := strings.Split(dataStr, "\n")
			for _, item := range dataList {
				if item == "" {
					continue
				}
				// fmt.Printf("[%d (len:%d)]: %s\n", index, len(item), item)
				itemList := strings.Split(item, ",")

				user := User{}
				isValidRecord := true
				for _, field := range itemList {
					values := strings.Split(field, ": ")
					if len(values) != 2 {
						fmt.Println("record is not vaild")
						isValidRecord = false
						break
					}
					fieldName := strings.TrimSpace(values[0])
					fieldValue := strings.TrimSpace(values[1])

					switch fieldName {
					case "id":
						id, err := strconv.Atoi(fieldValue)
						if err != nil {
							fmt.Println("strconv failed.", err)
							return
						}
						user.ID = id
					case "name":
						user.Name = fieldValue
					case "email":
						user.Email = fieldValue
					case "password":
						user.Password = fieldValue
					}
				}
				if isValidRecord {
					userStorage = append(userStorage, user)
				}
			}
			file.Close()
			fmt.Println(userStorage)
		}
	}
}

func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, duedate, category string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category Id")
	scanner.Scan()
	category = scanner.Text()

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("category id is not valid integer. %v\n", err)
		return
	}

	isFound := false
	for _, c := range categoryStorage {
		if c.UserID == authenticatedUser.ID && c.ID == categoryID {
			isFound = true

			break
		}
	}
	if !isFound {
		fmt.Print("category id is not found")
		return
	}

	fmt.Println("please enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()

	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		DueDate:    duedate,
		CategoryID: categoryID,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}
	taskStorage = append(taskStorage, task)
}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
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

	// fmt.Println("title: ", title, "-color:", color)

	category := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categoryStorage = append(categoryStorage, category)
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
	file, err := os.OpenFile(userStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("can't open the user.txt file", err)
		return
	}
	data := fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", user.ID, user.Name, user.Email, user.Password)
	b := []byte(data)
	file.Write(b)

	file.Close()
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
