package main

import (
	"quotes/repos"
	"quotes/server"
)

func main() {

	repo := repos.NewQuotesRepo()

	srv := server.NewHttpServer(":8080", repo)

	err := srv.Start()
	if err != nil {
		panic(err)
	}

}
