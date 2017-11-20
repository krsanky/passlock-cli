package main

import (
	"fmt"

	_ "github.com/krsanky/passlock-cli/model"
	"github.com/krsanky/passlock-cli/passlock/ui"
)

func main() {
	fmt.Println("passlock...")
	ui.Ui()
}
