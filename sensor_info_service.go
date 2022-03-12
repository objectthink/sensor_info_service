package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type SensorStatus struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Temperature int    `json:"temperature"`
	High        int    `json:"high"`
	Low         int    `json:"low"`
}

func natsSetup() {
	info := make(map[string]string)

	//get the map from disk
	jsonString, err := ioutil.ReadFile("/home/objectthink/info.json")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(jsonString)

		//json.Unmarshal([]byte[jsonString], &info)
		json.Unmarshal(jsonString, &info)

		fmt.Println(info["626"])
	}

	// Connect to a server
	nc, _ := nats.Connect("nats://192.168.86.31:4222")

	// listen for change events
	nc.Subscribe("*.event", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s %s\n", string(m.Data), string(m.Subject))

		//get sensor id
		tokens := strings.Split(m.Subject, ".")
		sensorId := tokens[0]

		//request
		//sensors cycle is about 3 seconds
		msg, _ := nc.Request(sensorId+".get", []byte("location"), 3000*time.Millisecond)
		fmt.Printf("%s\n", string(msg.Data))

		//store location
		info[sensorId] = string(msg.Data)

		// Marshal the map into a JSON string.
		json, err := json.Marshal(info)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(string(json))

			os.WriteFile("/home/objectthink/info.json", json, 0644)
			json, err = ioutil.ReadFile("/home/objectthink/info.json")

			fmt.Println(string(json))
		}
	})

	//listen for info requests
	//assume info sevice id is 000 for now
	nc.Subscribe("000.get.location", func(m *nats.Msg) {
		fmt.Printf("info service request: %s %s\n", string(m.Data), string(m.Subject))

		//get sensor id
		//tokens := strings.Split(m.Subject, ".")
		sensorId := string(m.Data)

		//request
		//sensors cycle is about 3 seconds
		//msg, _ := nc.Request(sensorId+".get", []byte("location"), 3000*time.Millisecond)
		//fmt.Printf("%s\n", string(msg.Data))

		//store location
		fmt.Println(info[sensorId])

		//publish reply
		nc.Publish(m.Reply, []byte(info[sensorId]))
	})
}

func main() {

	natsSetup()

	done := make(chan bool)
	go forever()
	<-done // Block forever

	//fmt.Scanln()
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}
