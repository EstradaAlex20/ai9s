package main

import (
	"log"

	"github.com/EstradaAlex20/ai9s/internal/ui"
)

func main() {
	app := ui.New("ai9s", ui.HardcodedRows())
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
