package collect

import (
	"NMS-Plugins/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"strconv"
	"strings"
	"sync"
)

var ch = make(chan string)

func Collect(request map[string]interface{}) (map[string]interface{}, error) {

	defer func() {

		if r := recover(); r != nil {
			request["result"] = "Failed " + fmt.Sprintf("%v", r)
			fmt.Println("Inside recover")
			response, err := json.Marshal(request)

			if err != nil {
				panic(err)
			}
			fmt.Println(response)
			log.Fatalf("%v", err)
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

	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
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

	response := make(map[string]interface{})

	go func(wg *sync.WaitGroup, response *map[string]interface{}) {

		result, err2 := g.Default.Get(scalerOids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			panic(err2)
		}

		for _, variable := range result.Variables {
			//fmt.Printf("%d: oid: %s ", i, variable.Name)

			switch variable.Type {
			case g.OctetString:
				//fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
				(*response)[utils.SCALERS[variable.Name]] = string(variable.Value.([]byte))
			default:
				//fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
				(*response)[utils.SCALERS[variable.Name]] = variable.Value
			}
		}

		wg.Done()

	}(&wg, &response)

	instanceResponce := make(map[string]map[string]interface{})

	go func(wg *sync.WaitGroup, instanceMap *map[string]map[string]interface{}) {

		for _, oid := range instanceOids {

			err2 := g.Default.BulkWalk(oid, func(pdu g.SnmpPDU) error {

				parts := strings.Split(pdu.Name, ".")

				index := fmt.Sprintf(parts[len(parts)-1])

				if (*instanceMap)[index] == nil {
					(*instanceMap)[index] = make(map[string]interface{})
					(*instanceMap)[index][utils.INSTANCES[oid]] = pdu.Value

					switch pdu.Type {
					case g.OctetString:
						b := pdu.Value.([]byte)

						if strings.Contains(pdu.Name, ".1.3.6.1.2.1.2.2.1.6") {
							(*instanceMap)[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
						} else {
							(*instanceMap)[index][utils.INSTANCES[oid]] = string(b)
						}
					default:
						//fmt.Printf(" %s\n", g.ToBigInt(pdu.Value))
						(*instanceMap)[index][utils.INSTANCES[oid]] = pdu.Value
					}

				} else {
					switch pdu.Type {
					case g.OctetString:
						b := pdu.Value.([]byte)

						if strings.Contains(pdu.Name, ".1.3.6.1.2.1.2.2.1.6") {
							(*instanceMap)[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
						} else {
							(*instanceMap)[index][utils.INSTANCES[oid]] = string(b)
						}
					default:
						(*instanceMap)[index][utils.INSTANCES[oid]] = pdu.Value
					}
				}

				return nil
			})

			if err2 != nil {
				panic(err2)
			}

		}

		//jsonObj2, err := json.Marshal(*instanceMap)

		//fmt.Println("\n\nJSON Of Instances : \n", string(jsonObj2), "\n\n")

		if err != nil {
			panic(err)
		}

		wg.Done()

	}(&wg, &instanceResponce)

	wg.Wait()

	response["interface"] = instanceResponce

	response2, err := json.Marshal(response)

	_ = response2

	if err != nil {
		panic(err)
	}

	return response, nil

}
