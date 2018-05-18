package main

import (
	"log"
	"nomad/api/src/app"
	"nomad/api/src/resources"
)

func main() {

	res, err := resources.GetInstance()
	if err != nil {
		log.Fatal(err)
	}

	app := &app.App{
		Resources: res,
	}
	app.Run()
}
