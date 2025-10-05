package main

import (
	"fmt"
	"os"

	"github.com/kaikelfalcao/github-activity/github"
)

func GetUsername() (string, error) {
	username :=  os.Args[1]

	if len(username) < 1 {
		return "", fmt.Errorf("username required!")
	}

	return username, nil
}

func main() {

	username, err := GetUsername()

	if err != nil {
		fmt.Print(err)
	}

	ghs := github.NewGitHubService(username)

	_, err = ghs.UserExists()

	if err != nil { 
		fmt.Println(err)
		return
	}

	response, err := ghs.GetActivities()

	if err != nil { 
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}