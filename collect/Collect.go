package collect

import (
	"NMS-Plugins/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
	"sync"
)

func Collect(request map[string]interface{}) (response map[string]interface{}, err error) {

	defer func() {

		if r := recover(); r != nil {
			response["message"] = fmt.Sprintf("%v", r)
			response["status"] = "Failed"
			err = fmt.Errorf("%v", r)
		}

	}()

	g.Default.Target = fmt.Sprintf("%v", request["ip"])
	intPort, _ := strconv.Atoi(fmt.Sprintf("%v", request["port"]))
	g.Default.Port = uint16(intPort)
	version := fmt.Sprintf("%v", request["version"])

	if strings.EqualFold(version, "v1") {
		g.Default.Version = g.Version1
	} else {
		g.Default.Version = g.Version2c
	}

	err2 := g.Default.Connect()
	if err2 != nil {
		panic(err2)
	}
	defer g.Default.Conn.Close()

	var wg sync.WaitGroup

	wg.Add(2)

	scalerOids := make([]string, len(utils.SCALERS))
	instanceOids := make([]string, len(utils.INSTANCES))

	i := 0

	for oid := range utils.SCALERS {
		scalerOids[i] = fmt.Sprintf("%v", oid)
		i++
	}

	j := 0
	for oid := range utils.INSTANCES {
		instanceOids[j] = fmt.Sprintf("%v", oid)
		j++
	}

	response = make(map[string]interface{})

	go func(wg *sync.WaitGroup, response *map[string]interface{}) {

		defer wg.Done()

		defer func() {

			if r := recover(); r != nil {
				(*response)["message"] = fmt.Sprintf("%v", r)
			}

		}()

		result, err := g.Default.Get(scalerOids) // Get() accepts up to g.MAX_OIDS
		if err != nil {
			panic(err)
		}

		for _, variable := range result.Variables {

			switch variable.Type {
			case g.OctetString:
				(*response)[utils.SCALERS[variable.Name]] = string(variable.Value.([]byte))
			default:
				(*response)[utils.SCALERS[variable.Name]] = fmt.Sprintf("%v", variable.Value)
			}
		}

	}(&wg, &response)

	numberOfInterface := 0

	func(number *int) {

		oid := []string{utils.CounterToOids["system.interfaces"]}

		result, err := g.Default.Get(oid) // Get() accepts up to g.MAX_OIDS
		if err != nil {
			panic(err)
		}

		for _, variable := range result.Variables {
			*number = variable.Value.(int)
		}
	}(&numberOfInterface)

	instanceResponce := make([]map[string]interface{}, numberOfInterface)

	go func(wg *sync.WaitGroup, instanceMap *[]map[string]interface{}) {
		defer wg.Done()

		defer func() {

			if r := recover(); r != nil {
				response["message"] = fmt.Sprintf("%v", r)
			}

		}()
		for _, oid := range instanceOids {

			err = g.Default.BulkWalk(oid, func(pdu g.SnmpPDU) error {

				parts := strings.Split(pdu.Name, ".")

				index2, err := strconv.Atoi(parts[len(parts)-1])
				_ = err
				index := index2 - 1

				if instanceResponce[index] == nil {
					instanceResponce[index] = make(map[string]interface{})
					instanceResponce[index][utils.INSTANCES[oid]] = pdu.Value

					switch pdu.Type {

					case g.OctetString:
						b := pdu.Value.([]byte)

						if strings.Contains(pdu.Name, utils.CounterToOids["interface.physical.address"]) {
							instanceResponce[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
						} else {
							instanceResponce[index][utils.INSTANCES[oid]] = string(b)
						}
					default:
						//fmt.Printf(" %s\n", g.ToBigInt(pdu.Value))
						instanceResponce[index][utils.INSTANCES[oid]] = fmt.Sprintf("%v", pdu.Value)
					}

				} else {
					switch pdu.Type {
					case g.OctetString:
						b := pdu.Value.([]byte)

						if strings.Contains(pdu.Name, utils.CounterToOids["interface.physical.address"]) {
							instanceResponce[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
						} else {
							instanceResponce[index][utils.INSTANCES[oid]] = string(b)
						}
					default:
						instanceResponce[index][utils.INSTANCES[oid]] = fmt.Sprintf("%v", pdu.Value)
					}
				}

				return nil
			})

			if err != nil {
				//panic(err)
			}

		}

		if err != nil {
			panic(err)
		}

	}(&wg, &instanceResponce)

	wg.Wait()

	response["interface"] = instanceResponce

	response2, err := json.Marshal(response)

	_ = response2

	if err != nil {
		panic(err)
	}

	return response, err

}
