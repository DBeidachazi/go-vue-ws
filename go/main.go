package main

import (
	"fmt"
	"goproject/model"
	"goproject/router"
)

func main() {
	db := model.BuildGen()
	router := router.InitRouter(db)
	fmt.Println("server is running on :3000")
	err := router.Run(":3000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
