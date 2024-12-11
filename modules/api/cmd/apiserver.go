package main

import (
	"log"

	"github.com/kubesphere-extensions/gateway-api/cmd/app"
)

func main() {

	cmd := app.NewAPIServerCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
