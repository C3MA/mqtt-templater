package main

import (
	"os"
	"fmt"
	"flag"
	"time"
)

func getTimestamp() []byte {
	return []byte(fmt.Sprintf("%d", time.Now().Unix()))
}

func main() {
	/* Parse the arguments */
	templateFilename := flag.String("template", "", "Input template filename")
	outputFilename := flag.String("output", "", "Output filename")
	brokerAddress := flag.String("broker", "localhost:1883", "Address of MQTT broker")
	useEncryption := flag.Bool("tls", false, "Use TLS encryption for MQTT connection")

	flag.Parse()

	if *templateFilename == "" {
		fmt.Fprintln(os.Stderr, "Missing parameter 'template'!")
		flag.PrintDefaults()
		return
	}
	if *outputFilename == "" {
		fmt.Fprintln(os.Stderr, "Missing parameter 'output'!")
		flag.PrintDefaults()
		return
	}

	/* Read the template */
	t := NewTemplate()
	err := t.ReadTemplate(*templateFilename)
	if err != nil {
		panic(fmt.Errorf("ReadTemplate: %v", err))
	}

	messages := make(chan Message)

	/* Connect to MQTT broker */
	cli := connectToBroker(*brokerAddress, *useEncryption)
	receiveMessages(cli, messages)

	/* Handle messages */
	for msg := range messages {
		isUpdate := t.SetVariable(msg.Key, msg.Value)
		if isUpdate == true {
			fmt.Println(msg.Key, "->", string(msg.Value))
			t.SetVariable("now", getTimestamp())

			err := t.WriteOutput(*outputFilename)
			if err != nil {
				panic(fmt.Errorf("WriteOutput: %v", err))
			}
		}
	}
}
