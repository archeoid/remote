package main
import (
	"github.com/archeoid/remote/pkg/client"
	"github.com/archeoid/remote/pkg/message"
	"log"
	"sync"
)




func main() {
	c := client.Connect("192.168.0.8:12340")

	var wg sync.WaitGroup

	client.Begin(c, &wg)

	c.Send <- message.Message{1, "cpu", "pls"}

	wg.Wait()

	log.Print("")
}


