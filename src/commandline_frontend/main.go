package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"k8s-project/commandline_frontend/api"
	"os"
	"strconv"
)

func main() {

	fmt.Println("--------------------")
	fmt.Println("Hello from cyberspace! Welcome to haaukins")
	fmt.Println("--------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("")
		fmt.Println("Enter 'create' to create a user")
		fmt.Println("Enter 'delete' to delete a user")
		fmt.Println("Enter 'users' to see users")
		fmt.Println("Enter 'on' to turn on a challenge")
		fmt.Println("Enter 'off' to turn off a challenge")
		fmt.Println("Enter 'challenges' to see the challenges/kali/wireguard running for a user")
		fmt.Println("Enter 'kali' to launch Kali VM that can be accessed in-browser")
		fmt.Println("Enter 'wg' to launch Wireguard")

		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "create":
			fmt.Println("Enter your username")
			scanner.Scan()
			user := scanner.Text()
			api.PostUser(user)
		case "delete":
			fmt.Println("Enter your username")
			scanner.Scan()
			user := scanner.Text()
			api.DeleteUser(user)
		case "users":
			api.GetUsers()

		case "on":
			fmt.Println("Enter the name of the challenge to turn on")
			scanner.Scan()
			challengeName := scanner.Text()

			fmt.Println("Enter the port of the challenge to turn on")
			scanner.Scan()
			port := scanner.Text()

			i, err1 := strconv.ParseInt(port, 10, 64)
			if err1 != nil {
				fmt.Println("hallo Write a real number!")
				panic(err1)
			}

			fmt.Println("Enter your username")
			scanner.Scan()
			user := scanner.Text()
			ports := []int32{int32(i)}
			api.PostChallenge(user, challengeName, ports)

		case "off":
			fmt.Println("Enter the name of the challenge you want to turn off")
			scanner.Scan()
			challengeName := scanner.Text()

			fmt.Println("Enter your username")
			scanner.Scan()
			user := scanner.Text()

			api.DeleteChallenge(user, challengeName)
		case "challenges":
			fmt.Println("Enter your username")
			scanner.Scan()
			name := scanner.Text()
			api.GetUserChallenges(name)
		case "kali":
			fmt.Println("Enter your username")
			scanner.Scan()
			user := scanner.Text()
			password := getUserPassword()
			api.PostKali(user, password)
		case "wg":
			fmt.Println("Enter your username")
			scanner.Scan()
			user := scanner.Text()
			api.PostWireguard(user)
		default:
			fmt.Println("Invalid input")
		}
	}
}

func getUserPassword() string {
	for {
		fmt.Print("Enter your password: ")
		bytePassword, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Error reading password")
		}
		password := string(bytePassword)

		fmt.Println("\nRe-enter your password: ")
		byteConfirm, err := term.ReadPassword(0)
		if err != nil {
			fmt.Println("Error reading password")
		}
		confirm := string(byteConfirm)

		if password == confirm {
			return password
		} else {
			fmt.Println("\nPasswords do not match! Please try again.")
		}
	}
}
