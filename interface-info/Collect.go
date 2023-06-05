package interface_info

import (
	"NMS-Plugins/utils"
	"encoding/hex"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
)

func Collect(snmp g.GoSNMP) (response map[string]interface{}, err error) {

	defer func() {

		if r := recover(); r != nil {
			response["message"] = fmt.Sprintf("%v", r)
			response["status"] = "Failed"
			err = fmt.Errorf("%v", r)
		}

	}()

	err = snmp.Connect()

	if err != nil {
		panic("Unable to connect to snmp")
	}

	defer func() {
		err = snmp.Conn.Close()

		if err != nil {
			panic("error While closing the snmp connection")
		}
	}()

	instanceOids := make([]string, len(utils.INSTANCES))

	j := 0
	for oid := range utils.INSTANCES {
		instanceOids[j] = fmt.Sprintf("%v", oid)
		j++
	}

	numberOfInterface := 0

	func(number *int) {

		oid := []string{utils.CounterToOids["system.interfaces"]}

		result, err := snmp.Get(oid) // Get() accepts up to g.MAX_OIDS
		if err != nil {
			panic(err)
		}

		for _, variable := range result.Variables {
			*number = variable.Value.(int)
		}
	}(&numberOfInterface)

	instanceResponse := make([]map[string]interface{}, numberOfInterface)

	for _, oid := range instanceOids {

		err = snmp.BulkWalk(oid, func(pdu g.SnmpPDU) error {

			parts := strings.Split(pdu.Name, ".")

			index2, err := strconv.Atoi(parts[len(parts)-1])
			_ = err
			index := index2 - 1

			if instanceResponse[index] == nil {
				instanceResponse[index] = make(map[string]interface{})
				instanceResponse[index][utils.INSTANCES[oid]] = pdu.Value

				switch pdu.Type {

				case g.OctetString:
					b := pdu.Value.([]byte)

					if strings.Contains(pdu.Name, utils.CounterToOids["interface.physical.address"]) {
						instanceResponse[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
					} else {
						instanceResponse[index][utils.INSTANCES[oid]] = string(b)
					}
				default:
					instanceResponse[index][utils.INSTANCES[oid]] = fmt.Sprintf("%v", pdu.Value)
				}

			} else {
				switch pdu.Type {
				case g.OctetString:
					b := pdu.Value.([]byte)

					if strings.Contains(pdu.Name, utils.CounterToOids["interface.physical.address"]) {
						instanceResponse[index][utils.INSTANCES[oid]] = hex.EncodeToString(b)
					} else {
						instanceResponse[index][utils.INSTANCES[oid]] = string(b)
					}
				default:
					instanceResponse[index][utils.INSTANCES[oid]] = fmt.Sprintf("%v", pdu.Value)
				}
			}

			return nil
		})

		if err != nil {
			//panic(err)
		}

	}
	response = make(map[string]interface{})

	response["interface"] = instanceResponse

	return response, nil
}
