package main

import (
	"NMS-Plugins/discovery"
	interfaceinfo "NMS-Plugins/interface-info"
	systeminfo "NMS-Plugins/system-info"
	"encoding/json"
	"flag"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"os"
	"strconv"
	"strings"
)

var request string

func init() {
	flag.StringVar(&request, "json", os.Args[1], "credential")
}

func getSNMP(requestInput map[string]interface{}) g.GoSNMP {

	g.Default.Target = fmt.Sprintf("%v", requestInput["ip"])
	intPort, _ := strconv.Atoi(fmt.Sprintf("%v", requestInput["port"]))
	g.Default.Port = uint16(intPort)
	g.Default.Retries = 0
	g.Default.Community = fmt.Sprintf("%v", requestInput["community"])
	version := fmt.Sprintf("%v", requestInput["version"])

	if strings.EqualFold(version, "v1") {
		g.Default.Version = g.Version1
	} else {
		g.Default.Version = g.Version2c
	}

	return *g.Default
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

		name, err := discovery.Discovery(getSNMP(requestMap))

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

		_, err := discovery.Discovery(getSNMP(requestMap))

		if err != nil {
			fmt.Println("Discovery Failed")
			err = fmt.Errorf("discovery failed (maybe device is unreachable)")
			panic(err)

		} else {

			metrics := fmt.Sprintf("%v", requestMap["metrics"])

			switch {

			case strings.EqualFold(metrics, "system.info"):
				response, err := systeminfo.Collect(getSNMP(requestMap))

				if err == nil {
					requestMap["result"] = response

					response, err := json.Marshal(requestMap)

					if err == nil {
						fmt.Println(string(response))
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}

			case strings.EqualFold(metrics, "interface.info"):
				response, err := interfaceinfo.Collect(getSNMP(requestMap))

				if err == nil {
					requestMap["result"] = response

					response, err := json.Marshal(requestMap)

					if err == nil {
						fmt.Println(string(response))
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}
			}

		}

	}

}
