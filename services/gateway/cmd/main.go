package main

import (
	"gateway/internal/app"
)

// @title           Twitter Clone
// @version         1.0
// @description     Application with basic functionality of twitter

// @host      localhost:7000
// @BasePath  /api/

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
