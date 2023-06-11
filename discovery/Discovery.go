package discovery

import (
	"NMS-Plugins/utils"
	"fmt"
	g "github.com/gosnmp/gosnmp"
)

func Discovery(snmp g.GoSNMP) (response map[string]interface{}, err error) {

	//var err error

	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
		}

	}()

	oidSlice := []string{utils.CounterToOIds["system.name"]}

	err = snmp.Connect()

	if err != nil {
		err = fmt.Errorf("connection failed : %v", err)
		return
	}

	defer func() {
		err = snmp.Conn.Close()
		if err != nil {
			err = fmt.Errorf("error While closing the snmp connection")
			return
		}

	}()

	result, err := snmp.Get(oidSlice)

	if err != nil {
		response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
		return
	}

	var systemName string

	for _, variable := range result.Variables {
		systemName = string(variable.Value.([]byte))
	}

	if err == nil {
		response[utils.STATUS] = utils.SUCCESS
		response[utils.SYSTEM_NAME] = systemName
	} else {
		response = utils.SetResponse(utils.FAILED, fmt.Sprintf("%v", err))
	}

	return
}
