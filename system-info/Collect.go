package system_info

import (
	"NMS-Plugins/utils"
	"fmt"
	g "github.com/gosnmp/gosnmp"
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

	scalerOids := make([]string, len(utils.SCALERS))

	i := 0

	for oid := range utils.SCALERS {
		scalerOids[i] = fmt.Sprintf("%v", oid)
		i++
	}

	result, err := snmp.Get(scalerOids)

	if err != nil {
		response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
		return nil, err
	}

	systemInfo := make(map[string]interface{})

	for _, variable := range result.Variables {

		switch variable.Type {
		case g.OctetString:
			systemInfo[utils.SCALERS[variable.Name]] = string(variable.Value.([]byte))
		default:
			systemInfo[utils.SCALERS[variable.Name]] = fmt.Sprintf("%v", variable.Value)
		}
	}

	if err == nil {
		dataMap := make(map[string]interface{})

		dataMap[utils.SYSTEM_INFO] = systemInfo

		response[utils.DATA] = dataMap

		response[utils.STATUS] = utils.SUCCESS
	} else {
		response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
	}

	return

}
