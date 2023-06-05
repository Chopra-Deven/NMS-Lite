package discovery

import (
	"NMS-Plugins/utils"
	"errors"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
)

func Discovery(request map[string]interface{}) (systemName string, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			systemName = ""
		}

	}()

	oidSlice := []string{utils.CounterToOids["system.name"]}

	g.Default.Target = fmt.Sprintf("%v", request["ip"])
	intPort, _ := strconv.Atoi(fmt.Sprintf("%v", request["port"]))
	g.Default.Port = uint16(intPort)
	g.Default.Retries = 0
	version := fmt.Sprintf("%v", request["version"])

	if strings.EqualFold(version, "v1") {
		g.Default.Version = g.Version1
	} else {
		g.Default.Version = g.Version2c
	}

	err = g.Default.Connect()

	defer g.Default.Conn.Close()

	if err != nil {
		systemName = ""
		return
	}

	defer g.Default.Conn.Close()

	result, err := g.Default.Get(oidSlice)

	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		systemName = ""
		return
	}

	for _, variable := range result.Variables {
		return string(variable.Value.([]byte)), err

	}

	return
}
