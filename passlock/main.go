package main

import (
	"fmt"

	"github.com/krsanky/passlock-cli/model"
)

func main() {
	fmt.Println("passlock...")
	model.CreateTable()
}
