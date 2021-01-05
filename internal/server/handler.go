package server

import (
	"remote/internal/message"
	//"log"
	"io/ioutil"
)

type msg = message.Message

func talkClient(client * Client) {
	for {
		resv := <-client.out

		//log.Printf("client sent: %s, %s\n", resv.Header, resv.Content)

		client.in <- getResponse(resv)
	}
}

func getResponse(m msg) msg {
	var o msg
	switch m.Subject {
	case "cpu":
		o = msg{message.ActionResponse, "cpu", getFile("/sys/class/thermal/thermal_zone3/temp")}
	}
	return o;
}

func getFile(path string) string {
	var str string;
	content, err := ioutil.ReadFile(path)
	if err != nil {
		str = "error"
	} else {
		str = string(content) 
	}
	return str
}