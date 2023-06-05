package discovery

import (
	"NMS-Plugins/utils"
	"errors"
	"fmt"
	g "github.com/gosnmp/gosnmp"
)

func Discovery(snmp g.GoSNMP) (systemName string, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			systemName = ""
		}

	}()

	oidSlice := []string{utils.CounterToOids["system.name"]}

	err = snmp.Connect()

	defer func() {
		err = snmp.Conn.Close()

		if err != nil {
			err = fmt.Errorf("error While closing the snmp connection")
		}

	}()

	if err != nil {
		systemName = ""
		return
	}

	result, err := snmp.Get(oidSlice)

	if err != nil {
		systemName = ""
		return
	}

	for _, variable := range result.Variables {
		return string(variable.Value.([]byte)), err

	}

	return
}
