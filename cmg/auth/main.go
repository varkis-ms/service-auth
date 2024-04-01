package main

import (
	"service-auth/internal/app"
)

const configsDir = "."

func main() {
	app.Run(configsDir)
}
