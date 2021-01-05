package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"github.com/archeoid/remote/pkg/message"
	"sync"
	"io/ioutil"
)

type Client struct {
	in chan message.Message
	out chan message.Message
	ip string
}

func newClient(ip string) Client {
	return Client{make(chan message.Message, 10), make(chan message.Message, 10), ip}
}

func StartServer(bind string, clients []Client) {
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	caCert, _ := ioutil.ReadFile("certs/ca.crt")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := tls.Config {		
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAnyClientCert,
		RootCAs: caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	config.BuildNameToCertificate()

	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", bind, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}

		clients = append(clients, newClient("someshit"))

		c := &clients[len(clients)-1]

		go talkClient(c)
		
		go handleClient(conn, c)
	}
}

func handleClient(conn net.Conn, client * Client) {
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go clientRead(conn, &wg, client.out)
	go clientWrite(conn, client.in)

	wg.Wait()

	client.in <- message.Message{0,"close", "close"}

	log.Printf("server: disconnected %s", client.ip)
}

func clientRead(conn net.Conn, wg * sync.WaitGroup, out chan message.Message) {
	defer wg.Done()

	var size int
	buf := make([]byte, 4)
	for {
		_, err := conn.Read(buf)
		size = message.BinaryToInt(buf)
		if err != nil {
			log.Printf("server: read1: %s", err)
			break
		}
		buf = make([]byte, size)
		_, err = conn.Read(buf)
		out <- message.FromBytes(buf)
		buf = make([]byte, 4)
		if err != nil {
			log.Printf("server: read2: %s", err)
			break
		}
	}
	log.Print("server: read end")
}

func clientWrite(conn net.Conn, in chan message.Message) {
	for {
		msg := <-in
		buf := message.ToBytes(msg)
		_, err := conn.Write(message.IntToBinary(len(buf)))
		_, err = conn.Write(buf)
		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Print("server: write end")
}