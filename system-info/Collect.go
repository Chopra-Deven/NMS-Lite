package system_info

import (
	"NMS-Plugins/utils"
	"fmt"
	g "github.com/gosnmp/gosnmp"
)

func Collect(snmp g.GoSNMP) (response map[string]interface{}, err error) {

	defer func() {

		if r := recover(); r != nil {
			response["message"] = fmt.Sprintf("%v", r)
			response["status"] = "Failed"
			err = fmt.Errorf("%v", r)
		}

	}()

	err2 := snmp.Connect()

	if err2 != nil {
		panic(err2)
	}

	defer func() {
		err = snmp.Conn.Close()

		if err != nil {
			err = fmt.Errorf("error While closing the snmp connection")
		}

	}()

	scalerOids := make([]string, len(utils.SCALERS))

	i := 0

	for oid := range utils.SCALERS {
		scalerOids[i] = fmt.Sprintf("%v", oid)
		i++
	}

	response = make(map[string]interface{})

	result, err := snmp.Get(scalerOids) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		panic(err)
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

	response["system.info"] = systemInfo

	return

}
