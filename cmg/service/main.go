package main

import (
	"github.com/varkis-ms/service-auth/internal/app"
)

const configsDir = "."

func main() {
	app.Run(configsDir)
}
