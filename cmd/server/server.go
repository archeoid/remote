package main
import (
	"remote/internal/server"
	"log"
)

func main() {

	var clients []server.Client

	server.StartServer("0.0.0.0:12340", clients);
	log.Print("")
}
