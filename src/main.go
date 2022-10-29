package main

import (
	"log"

	"github.com/kujilabo/cocotola/src/uuid"
)

func main() {
	id, err := uuid.Generate()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HELLO")
	log.Println(id)
}
