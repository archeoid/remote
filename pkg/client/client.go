package client

import (
	"github.com/archeoid/remote/pkg/message"
	"crypto/tls"
	"crypto/x509"
	"log"
	"io/ioutil"
)

const (
	Disconnected = 1
	Connecting   = 2
	Connected    = 3
)

type Test struct {
	Send chan []byte
	Resv chan []byte
}

type Client struct {
	Send chan message.Message
	Resv chan message.Message
	Status int
	Conn *tls.Conn
}

func Connect(ip string) Client {

	c := Client{make(chan message.Message, 10), make(chan message.Message, 10), Disconnected, nil}

	caCert, err := ioutil.ReadFile("certs/ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}
	

	config := tls.Config {		
		RootCAs: caCertPool,
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	log.Printf("conntecting to %s", ip)

	c.Conn, err = tls.Dial("tcp", ip, &config);
	if err != nil {
		c.Status = Disconnected
		panic("failed to connect: " + err.Error())
	} else {
		c.Status = Connected
		log.Print("conntected")
	}
	return c
}

func Read(c Client) {
	var size int
	buf := make([]byte, 4)
	for {
		_, err := c.Conn.Read(buf)
		size = message.BinaryToInt(buf)
		if err != nil {
			log.Printf("server: read1: %s", err)
			break
		}
		buf = make([]byte, size)
		_, err = c.Conn.Read(buf)
		c.Resv <- message.FromBytes(buf)
		buf = make([]byte, 4)
		if err != nil {
			log.Printf("server: read2: %s", err)
			break
		}
		log.Print("response")
	}
	log.Print("server: read end")
}

func Write(c Client) {
	for {
		msg := <- c.Send
		log.Print("sending")
		buf := message.ToBytes(msg)
		_, err := c.Conn.Write(message.IntToBinary(len(buf)))
		_, err = c.Conn.Write(buf)
		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Print("server: write end")
}