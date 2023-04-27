package main

import (
	"bufio"
	"fmt"
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
		fmt.Println("Write 'create' to create a User")
		fmt.Println("Write 'delete' to delete a User")
		fmt.Println("Write 'on' to turn on a challenge")
		fmt.Println("Write 'off' to turn off a challenge")
		fmt.Println("Write 'kali' to launch VM with selected challenges via vnc")

		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "create":
			fmt.Println("Write your team alias")
			scanner.Scan()
			user := scanner.Text()
			api.PostUser(user)

		case "on":
			fmt.Println("Write the name of the challenge to turn on")
			scanner.Scan()
			challengeName := scanner.Text()

			fmt.Println("Write the port of the challenge to turn on")
			scanner.Scan()
			port := scanner.Text()

			i, err1 := strconv.ParseInt(port, 10, 64)
			if err1 != nil {
				fmt.Println("hallo Write a real number!")
				panic(err1)
			}

			fmt.Println("Write the user of the challenge to turn on")
			scanner.Scan()
			user := scanner.Text()
			ports := []int32{int32(i)}
			api.PostChallenge(user, challengeName, ports)

		case "off":
			fmt.Println("Write the name of the challenge you want to turn off")
			scanner.Scan()
			challengeName := scanner.Text()

			fmt.Println("Write the user of the challenge to turn on")
			scanner.Scan()
			user := scanner.Text()

			api.DeleteChallenge(user, challengeName)
		case "kali":
			fmt.Println("Write the user of the challenge to turn on")
			scanner.Scan()
			user := scanner.Text()
			api.PostKali(user)
		default:
			fmt.Println("Invalid input")
		}
	}
}
