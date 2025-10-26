package main

import (
	"backend/router"
)

func main() {
	router := router.Router()
	router.Run(":3000")
}
