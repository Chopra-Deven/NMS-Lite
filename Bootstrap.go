package main

import (
	"NMS-Plugins/collect"
	"NMS-Plugins/discovery"
	"encoding/json"
	"flag"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"os"
	"strconv"
	"strings"
)

var request string

func init() {
	flag.StringVar(&request, "json", os.Args[1], "credential")
}

func main() {

	flag.Parse()

	var requestMap map[string]interface{}

	var err = json.Unmarshal([]byte(request), &requestMap)

	defer func() {

		if r := recover(); r != nil {
			requestMap["result"] = "Failed " + fmt.Sprintf("%v", r)

			response, err := json.Marshal(requestMap)

			if err != nil {
				panic(err)
			}
			fmt.Println(response)
			log.Fatalf("%v", err)
		}

	}()

	if err != nil {
		panic(err)
	}

	g.Default.Target = fmt.Sprintf("%v", requestMap["ip"])
	intPort, _ := strconv.Atoi(fmt.Sprintf("%v", requestMap["port"]))
	g.Default.Port = uint16(intPort)
	version := fmt.Sprintf("%v", requestMap["version"])

	if strings.EqualFold(version, "v1") {
		g.Default.Version = g.Version1
	} else {
		g.Default.Version = g.Version2c
	}

	err = g.Default.Connect()

	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	if requestMap["type"] == "discovery" {

		id, err := discovery.Discovery(requestMap)

		if err != nil {
			panic(err)

		} else {

			result := make(map[string]string)

			result["status"] = "success"
			result["system.name"] = id

			requestMap["result"] = result

			response, err := json.Marshal(requestMap)

			if err != nil {
				panic(err)
			}

			fmt.Println(string(response))
		}
		fmt.Println("\nDisovery end\n")

	} else if requestMap["type"] == "provision" {

		//fmt.Println("\nProvision Start\n")

		_, err := discovery.Discovery(requestMap)

		if err != nil {
			fmt.Println("Discovery error")
			panic(err)

		} else {

			response, err := collect.Collect(requestMap)

			if err != nil {
				panic(err)
			}

			requestMap["result"] = response

			response2, err := json.Marshal(requestMap)

			if err != nil {
				panic(err)
			}

			fmt.Println(string(response2))

			//fmt.Println(string(response))
		}
		//fmt.Println("\nProvision Ended\n")

	} else {

	}

}
