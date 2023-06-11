package main

import (
	"NMS-Plugins/discovery"
	interfaceinfo "NMS-Plugins/interface-info"
	systeminfo "NMS-Plugins/system-info"
	"NMS-Plugins/utils"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var request string

func init() {
	flag.StringVar(&request, "json", os.Args[1], "user input")
}

func main() {

	flag.Parse()

	var requestMap map[string]interface{}

	var err = json.Unmarshal([]byte(request), &requestMap)

	defer func() {

		if r := recover(); r != nil {

			requestMap[utils.RESULT] = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", r))

			result, _ := json.Marshal(requestMap)

			fmt.Println(string(result))
		}

	}()

	if err != nil {
		panic(err)
	}

	snmp, errorMessage := utils.GetSNMP2(requestMap)

	if !strings.EqualFold(errorMessage, utils.EMPTY_STRING) {
		panic(errorMessage)
	}

	if requestMap[utils.TYPE] == utils.DISCOVERY {

		data, err := discovery.Discovery(snmp)

		if err != nil || data == nil {

			requestMap[utils.RESULT] = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
			response, _ := json.Marshal(requestMap)

			fmt.Println(string(response))

		} else {

			requestMap[utils.RESULT] = data

			response, err := json.Marshal(requestMap)

			if err != nil {
				panic(err)
			}

			fmt.Println(string(response))
		}

	} else {

		metrics := fmt.Sprintf("%v", requestMap[utils.METRICS])

		switch {

		case strings.EqualFold(metrics, utils.SYSTEM_INFO):

			data, err := systeminfo.Collect(snmp)

			if err == nil {

				requestMap[utils.RESULT] = data

				response, err := json.Marshal(requestMap)

				if err == nil {
					fmt.Println(string(response))
				} else {
					panic(err)
				}
			} else {
				requestMap[utils.RESULT] = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))

				result, _ := json.Marshal(requestMap)

				fmt.Println(string(result))
			}

		case strings.EqualFold(metrics, utils.INTERFACE_INFO):

			data, err := interfaceinfo.Collect(snmp)

			if err == nil {

				requestMap[utils.RESULT] = data

				response, err := json.Marshal(requestMap)

				if err == nil {
					fmt.Println(string(response))
				} else {
					panic(err)
				}

			} else {
				requestMap[utils.RESULT] = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))

				result, _ := json.Marshal(requestMap)

				fmt.Println(string(result))
			}
		}

	}

}
