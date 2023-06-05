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

func Collect2(request map[string]interface{}) (response map[string]interface{}, err error) {

	defer func() {

		if r := recover(); r != nil {
			response["message"] = fmt.Sprintf("%v", r)
			response["status"] = "Failed"
			response, err := json.Marshal(request)

			if err != nil {
				panic(err)
			}
			fmt.Println(response)
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

	err = g.Default.Connect()
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

	response = make(map[string]interface{})

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

	numberOfInterface := 0

	func(number *int) {

		oid := []string{utils.CounterToOids["system.interfaces"]}

		result, err2 := g.Default.Get(oid) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			panic(err2)
		}

		for _, variable := range result.Variables {
			*number = variable.Value.(int)
		}
	}(&numberOfInterface)

	fmt.Printf("\n\nNumber Of Interface %v and Type : %T\n\n", numberOfInterface, numberOfInterface)

	instanceResponce := make([]map[string]interface{}, numberOfInterface)

	go func(wg *sync.WaitGroup, instanceMap *[]map[string]interface{}) {

		for _, oid := range instanceOids {

			err2 := g.Default.BulkWalk(oid, func(pdu g.SnmpPDU) error {

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

						if strings.Contains(pdu.Name, ".1.3.6.1.2.1.2.2.1.6") {
							instanceResponce[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
						} else {
							instanceResponce[index][utils.INSTANCES[oid]] = string(b)
						}
					default:
						//fmt.Printf(" %s\n", g.ToBigInt(pdu.Value))
						instanceResponce[index][utils.INSTANCES[oid]] = pdu.Value
					}

				} else {
					switch pdu.Type {
					case g.OctetString:
						b := pdu.Value.([]byte)

						if strings.Contains(pdu.Name, ".1.3.6.1.2.1.2.2.1.6") {
							instanceResponce[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
						} else {
							instanceResponce[index][utils.INSTANCES[oid]] = string(b)
						}
					default:
						instanceResponce[index][utils.INSTANCES[oid]] = pdu.Value
					}
				}

				return nil
			})

			if err2 != nil {
				panic(err2)
			}

		}

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

	return response, err

}
