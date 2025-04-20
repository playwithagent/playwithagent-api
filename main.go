package main

import (
	"fmt"
	"playwithagent-xo/cmd/httpserver"
	"playwithagent-xo/config"

)

func main(){

	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Port: 8080,
		},
	}
	
	
	server := httpserver.NewServer(cfg)
	fmt.Println("Start Echo server")

	server.Serve()
	
}