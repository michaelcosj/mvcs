package main

import (
	"log"
	"michaelcosj/mvcs/app"
)

func main() {
  if err := app.Run(); err != nil {
    log.Fatal(err)
  }
}
