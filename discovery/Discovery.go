package discovery

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"strconv"
	"strings"
)

func Discovery(request map[string]interface{}) (string, error) {

	defer func() {
		if r := recover(); r != nil {
			request["result"] = "Failed " + fmt.Sprintf("%v", r)

			response, err := json.Marshal(request)

			if err != nil {
				panic(err)
			}
			fmt.Println(response)
			log.Fatalf("%v", err)
		}
	}()

	oidSlice := []string{".1.3.6.1.2.1.1.5.0"}

	g.Default.Target = fmt.Sprintf("%v", request["ip"])
	//g.Default.Community = fmt.Sprintf("%v", request["community"])
	//g.Default.Port = request["port"].(uint16)
	intPort, _ := strconv.Atoi(fmt.Sprintf("%v", request["port"]))
	g.Default.Port = uint16(intPort)
	version := fmt.Sprintf("%v", request["version"])

	if strings.EqualFold(version, "v1") {
		g.Default.Version = g.Version1
	} else {
		g.Default.Version = g.Version2c
	}

	err := g.Default.Connect()

	defer g.Default.Conn.Close()

	if err != nil {
		return "", err
	}

	defer g.Default.Conn.Close()

	result, err2 := g.Default.Get(oidSlice) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		//log.Fatalf("Get() err: %v", err2)
		return "", err2
	}

	for _, variable := range result.Variables {

		switch variable.Type {
		case g.OctetString:
			return string(variable.Value.([]byte)), nil
		default:

			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
	return "", nil
}
