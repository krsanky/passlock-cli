package main

import (
	"fmt"
	"os"

	"github.com/krsanky/passlock-cli/model"
	"github.com/krsanky/passlock-cli/passlock/ui"
	"github.com/krsanky/passlock-cli/passlock/ui2"
)

func main() {

	if len(os.Args) >= 2 {
		switch arg1 := os.Args[1]; arg1 {
		case "alt":
			ui.Ui2()
		case "cellview":
			ui2.CellView()
		case "simple":
			ui2.Simple()
		default:
			usage()
		}
	} else {
		usage()
	}
	model.Close()
}

func usage() {
	fmt.Println()
	fmt.Println(`passlock [ui2|cellview]`)
	fmt.Println()
}
