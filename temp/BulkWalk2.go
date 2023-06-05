package main

import (
	"NMS-Plugins/utils"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
)

func main() {
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	snmp.Target = "172.16.8.2"
	err := snmp.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer snmp.Conn.Close()

	instanceOids := make([]string, len(utils.INSTANCES))

	j := 0
	for _, value := range utils.SCALERS {
		instanceOids[j] = fmt.Sprintf("%v", value)
		j++
	}

	response := make(map[string]interface{})

	_ = response

	/*var printValue = func(pdu g.SnmpPDU) error {
		fmt.Printf("%s : ", pdu.Name)

		switch pdu.Type {
		case g.OctetString:
			b := pdu.Value.([]byte)
			//fmt.Printf("STRING: %s\n", string(b))
			fmt.Printf(" Type : %T - %s\n", pdu.Type, string(b))
			//fmt.Printf("STRING: %s\n", hex.EncodeToString(b))

		default:
			fmt.Printf(" %s\n", g.ToBigInt(pdu.Value))
		}
		return nil
	}
	*/
	var oidsSlice = []string{
		".1.3.6.1.2.1.2.2.1.1",
		".1.3.6.1.2.1.31.1.1.1.1",
	}

	for _, oid := range oidsSlice {
		err2 := snmp.BulkWalk(oid, func(dataUnit g.SnmpPDU) error {

			switch dataUnit.Type {
			case g.OctetString:
				b := dataUnit.Value.([]byte)
				//fmt.Printf("STRING: %s\n", string(b))
				fmt.Printf(" Type : %T - %s\n", dataUnit.Type, string(b))
				//fmt.Printf("STRING: %s\n", hex.EncodeToString(b))

			default:
				fmt.Printf(" %s\n", g.ToBigInt(dataUnit.Value))
			}

			return nil
		}) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			panic(err2)
		}
	}
}
