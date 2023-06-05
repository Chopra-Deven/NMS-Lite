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

	//oids := []string{utils.[]}
	scalerOids := make([]string, len(utils.SCALERS))
	//instanceOids := make([]string, len(utils.INSTANCES))

	i := 0

	for oid := range utils.SCALERS {
		scalerOids[i] = fmt.Sprintf("%v", oid)
		i++
	}

	result, err2 := snmp.Get(scalerOids) // Get() accepts up to g.MAX_OIDS

	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	//fmt.Println("\n\n", result)

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case g.OctetString:
			fmt.Printf("%s : %s\n", variable.Name, string(variable.Value.([]byte)))
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
}
