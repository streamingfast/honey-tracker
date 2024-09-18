package main

import "github.com/streamingfast/honey-tracker/web"

func main() {
	server := &web.Server{}
	server.ServeHttp()
}
