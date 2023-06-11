package utils

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
)

func GetSNMP(requestInput map[string]interface{}) g.GoSNMP {

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

func GetSNMP2(requestInput map[string]interface{}) (g.GoSNMP, string) {

	errorMessage := ""

	if requestInput[IP] != nil {
		g.Default.Target = fmt.Sprintf("%v", requestInput["ip"])
	} else {
		errorMessage += " Ip not provided."
	}

	if requestInput[PORT] != nil {

		intPort, _ := strconv.Atoi(fmt.Sprintf("%v", requestInput["port"]))

		g.Default.Port = uint16(intPort)

	} else {
		errorMessage += " port not provided."
	}

	g.Default.Retries = DEFAULT_RETRIES

	if requestInput[COMMUNITY] != nil {

		g.Default.Community = fmt.Sprintf("%v", requestInput["community"])

	} else {
		errorMessage += " Community not provided."
	}

	if requestInput[COMMUNITY] != nil {

		version := fmt.Sprintf("%v", requestInput["version"])

		if strings.EqualFold(version, "v1") {
			g.Default.Version = g.Version1
		} else {
			g.Default.Version = g.Version2c
		}

	} else {
		errorMessage += " Community not provided."
	}

	if strings.EqualFold(errorMessage, EMPTY_STRING) {

	}

	return *g.Default, errorMessage
}

func SetResponse(status string, message string) (response map[string]interface{}) {

	response = make(map[string]interface{})

	response[MESSAGE] = message

	response[STATUS] = status

	return
}
