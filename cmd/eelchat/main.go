package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apachejuice/eelchat/client/api"
)

func main() {
	client, err := api.NewClient("https://apachejuice.dev:5555")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.CreateUser(context.TODO(), api.User{
		Username: "apachejuice",
		Password: "this is an example password",
	})

	if err != nil {
		log.Fatal(err)
	}

	switch resp := resp.(type) {
	case *api.CreateUserNoContent:
		fmt.Println("successfully created user")
	case *api.CreateUserApplicationJSONBadRequest:
		fmt.Println("bad request:", resp.Message)
	case *api.CreateUserApplicationJSONInternalServerError:
		fmt.Println("internal server error:", resp.Message)
	default:
		fmt.Println("response:", resp)
	}
}
