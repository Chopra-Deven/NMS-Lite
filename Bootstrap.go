package main

import (
	"NMS-Plugins/collect"
	"NMS-Plugins/discovery"
	"encoding/json"
	"flag"
	"fmt"
	"os"
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
			requestMap["status"] = "Failed"

			requestMap["message"] = fmt.Sprintf("%v", r)

			response, err := json.Marshal(requestMap)

			if err != nil {
				panic(err)
			}

			fmt.Println(string(response))

		}

	}()

	if err != nil {
		panic(err)
	}

	if requestMap["type"] == "discovery" {

		name, err := discovery.Discovery(requestMap)

		if err != nil && name == "" {
			panic(err)

		} else {

			result := make(map[string]string)

			result["status"] = "success"
			result["system.name"] = name

			requestMap["result"] = result

			response, err := json.Marshal(requestMap)

			if err != nil {
				panic(err)
			}

			fmt.Println(string(response))
		}

	} else {

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

		}

	}

}
