package main

import (
	"log"

	"github.com/michaelcosj/mvcs/app"
)

func main() {
  if err := app.Run(); err != nil {
    log.Fatal(err)
  }
}
