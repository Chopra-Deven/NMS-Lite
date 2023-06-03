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
	g.Default.Target = "172.16.8.2"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	//oids := []string{"1.1.3.6.1.2.1.1.5", "1.1.3.6.1.2.1.1.1"}

	instanceOids := make([]string, len(utils.INSTANCES))

	j := 0
	for _, value := range utils.SCALERS {
		instanceOids[j] = fmt.Sprintf("%v", value)
		j++
	}

	/*err2 := g.Default.BulkWalk(instanceOids, printValue) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}*/

	//for _, oid := range instanceOids {
	err2 := g.Default.BulkWalk(".1.3.6.1.2.1.2.2.1.6", printValue) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		panic(err2)
	}

}

func printValue(pdu g.SnmpPDU) error {
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
