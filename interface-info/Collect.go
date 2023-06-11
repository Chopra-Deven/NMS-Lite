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

	response = make(map[string]interface{})

	defer func() {

		if r := recover(); r != nil {
			response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
		}

	}()

	err = snmp.Connect()

	if err != nil {

		err = fmt.Errorf("connection failed : %v", err)

		return
	}

	defer func() {
		err = snmp.Conn.Close()

		if err != nil {

			err = fmt.Errorf("error While closing the snmp connection : %v", err)

			return
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

		oid := []string{utils.CounterToOIds["system.interfaces"]}

		result, err := snmp.Get(oid)
		if err != nil {
			response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
			return
		}

		for _, variable := range result.Variables {
			*number = variable.Value.(int)
		}
	}(&numberOfInterface)

	instanceResponse := make(map[int]map[string]interface{})

	var walkOrBulkWalk = snmp.BulkWalk

	if snmp.Version == g.Version1 {
		walkOrBulkWalk = snmp.Walk
	}

	for _, oid := range instanceOids {

		err = walkOrBulkWalk(oid, func(pdu g.SnmpPDU) error {

			parts := strings.Split(pdu.Name, ".")

			index, _ := strconv.Atoi(parts[len(parts)-1])

			if instanceResponse[index] == nil {
				instanceResponse[index] = make(map[string]interface{})
				instanceResponse[index][utils.INSTANCES[oid]] = pdu.Value

				switch pdu.Type {

				case g.OctetString:
					b := pdu.Value.([]byte)

					if strings.Contains(pdu.Name, utils.CounterToOIds["interface.physical.address"]) {
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

					if strings.Contains(pdu.Name, utils.CounterToOIds["interface.physical.address"]) {
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

	}

	data := make([]map[string]interface{}, len(instanceResponse))

	i := 0
	for _, interfaceInstance := range instanceResponse {
		data[i] = interfaceInstance
		i++
	}

	if err == nil {

		dataMap := make(map[string]interface{})

		dataMap[utils.INTERFACE_INFO] = data

		response[utils.STATUS] = utils.SUCCESS
		response[utils.DATA] = dataMap
	} else {
		response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
	}

	return response, nil
}
