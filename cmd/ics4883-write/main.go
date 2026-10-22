/*
ICS4883 Instrument coupler test code
output tests - writung to output connector
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/msiegen/linuxgpib"
	"github.com/msiegen/linuxgpib/internal"
)

var CONNECTOR_SIZE = 50

type ConnectorNumbering struct {
	byte int
	bit  int
	info string
}

var outputConn = [50 + 1]ConnectorNumbering{}

func fillOutputConnectorPinMap() {

	outputConn[0] = ConnectorNumbering{0, 0, "N.A."} // no index 0 at connector, starts with 1
	outputConn[1] = ConnectorNumbering{0, 0, "GND"}  // GND
	outputConn[2] = ConnectorNumbering{3, 1, ""}
	outputConn[3] = ConnectorNumbering{3, 2, ""}
	outputConn[4] = ConnectorNumbering{3, 3, ""}
	outputConn[5] = ConnectorNumbering{0, 0, "VCC"} // VCC
	outputConn[6] = ConnectorNumbering{3, 4, ""}

	outputConn[7] = ConnectorNumbering{4, 4, ""} // out-of-order

	outputConn[8] = ConnectorNumbering{2, 4, ""}
	outputConn[9] = ConnectorNumbering{2, 3, ""}
	outputConn[10] = ConnectorNumbering{2, 2, ""}
	outputConn[11] = ConnectorNumbering{2, 1, ""}

	outputConn[12] = ConnectorNumbering{8, 1, ""} // out-of-order

	outputConn[13] = ConnectorNumbering{1, 1, ""}
	outputConn[14] = ConnectorNumbering{1, 2, ""}
	outputConn[15] = ConnectorNumbering{1, 3, ""}
	outputConn[16] = ConnectorNumbering{1, 4, ""}

	outputConn[17] = ConnectorNumbering{7, 4, ""}
	outputConn[18] = ConnectorNumbering{7, 3, ""}
	outputConn[19] = ConnectorNumbering{7, 2, ""}
	outputConn[20] = ConnectorNumbering{7, 1, ""}

	outputConn[21] = ConnectorNumbering{0, 0, "HANDSHAKE_ENABLE"} // HANDSHAKE ENABLE
	outputConn[22] = ConnectorNumbering{0, 0, "CONTROL OUT 1"}    // CONTROL OUT 1
	outputConn[23] = ConnectorNumbering{0, 0, "CONTROL OUT 2"}    // CONTROL OUT 2
	outputConn[24] = ConnectorNumbering{0, 0, "CONTROL OUT 3"}    // CONTROL OUT 3
	outputConn[25] = ConnectorNumbering{0, 0, "CONTROL OUT 4"}    // CONTROL OUT 4

	outputConn[26] = ConnectorNumbering{5, 2, ""}
	outputConn[27] = ConnectorNumbering{5, 3, ""}
	outputConn[28] = ConnectorNumbering{5, 4, ""}
	outputConn[29] = ConnectorNumbering{5, 1, ""}

	outputConn[30] = ConnectorNumbering{10, 4, ""}
	outputConn[31] = ConnectorNumbering{10, 2, ""}
	outputConn[32] = ConnectorNumbering{10, 3, ""}
	outputConn[33] = ConnectorNumbering{10, 1, ""}

	outputConn[34] = ConnectorNumbering{9, 1, ""}
	outputConn[35] = ConnectorNumbering{9, 2, ""}
	outputConn[36] = ConnectorNumbering{9, 4, ""}
	outputConn[37] = ConnectorNumbering{9, 3, ""}

	outputConn[38] = ConnectorNumbering{4, 1, ""}
	outputConn[39] = ConnectorNumbering{4, 2, ""}
	outputConn[40] = ConnectorNumbering{4, 3, ""}

	outputConn[41] = ConnectorNumbering{8, 4, ""}
	outputConn[42] = ConnectorNumbering{8, 3, ""}
	outputConn[43] = ConnectorNumbering{8, 2, ""}

	outputConn[44] = ConnectorNumbering{6, 4, ""}
	outputConn[45] = ConnectorNumbering{6, 3, ""}
	outputConn[46] = ConnectorNumbering{6, 2, ""}
	outputConn[47] = ConnectorNumbering{6, 1, ""}

	outputConn[48] = ConnectorNumbering{0, 0, "CONTROL IN 1"} // CONTROL IN 1
	outputConn[49] = ConnectorNumbering{0, 0, "CONTROL IN 2"} // CONTROL IN 2
	outputConn[50] = ConnectorNumbering{0, 0, "CONTROL IN 3"} // CONTROL IN 3

	for i := 0; i < len(outputConn); i++ {
		if i == 0 {
			continue
		}
		if outputConn[i].byte != 0 && outputConn[i].bit != 0 {
			fmt.Printf("pin[%v] = { byte=%v, bit=%v }\n", i, outputConn[i].byte, outputConn[i].bit)
		} else {
			fmt.Printf("pin[%v] %v\n", i, outputConn[i].info)
		}
	}
}

func main() {
	verbose := flag.Bool(
		"verbose", false,
		"Turn on logging. Without this, only the result or first error is printed.",
	)
	board := flag.Int(
		"board", 0,
		"The board number. Defaults to zero, which corresponds to /dev/gpib0.",
	)
	address := flag.Int(
		"address", -1,
		"The primary address of the GPIB device (ICS4883) to query.",
	)
	pin := flag.Int(
		"pin", -1,
		"pin to use [1..40]",
	)
	value := flag.Int(
		"value", -1,
		"value for pin [0,1]",
	)
	dataString := flag.String(
		"data", "",
		"data string like 0000000000",
	)
	blink := flag.Bool(
		"blink", false,
		"switch pin on and off and on and off and...",
	)

	flag.Parse()

	if *address < 0 {
		fmt.Fprintln(os.Stderr, "Please specify an -address!")
		os.Exit(1)
	}

	var opts []linuxgpib.Option
	if *verbose {
		opts = append(opts, linuxgpib.Log(log.Default()))
	}

	fillOutputConnectorPinMap()

	d, err := linuxgpib.NewDevice(*board, *address, opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open device:", err)
		os.Exit(1)
	}
	defer d.Close()

	if *blink == false {
	} else {
		fmt.Printf("blinking pin %v\n", *pin)
	}
	if *dataString != "" {
		fmt.Printf("data %v\n", *dataString)
	}

	if *blink == false {
		// set or clear a pin
		// Write a single data string based on pin number
		dString := "0000000000"
		if *value == 1 {
			dString = createDataString(*pin)
		}

		fmt.Printf("Data string for pin %v:  %v\n", *pin, dString)
		if _, err := fmt.Fprintln(d, dString); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write to device:", err)
			os.Exit(1)
		}
	} else {
		// blink a pin, f ~ 5 Hz
		offStr := "0000000000"
		onStr := *dataString
		for true {
			if _, err := fmt.Fprintln(d, onStr); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write to device:", err)
				os.Exit(1)
			}
			time.Sleep(100 * time.Millisecond)
			if _, err := fmt.Fprintln(d, offStr); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write to device:", err)
				os.Exit(1)
			}
			time.Sleep(100 * time.Millisecond)
		}

	}

	lines := internal.IblinesString(*board)
	fmt.Printf("GPIB Lines: %v\n", lines)
}

func createDataString(num int) string {
	// String looks like "0080000000"
	// num is in [1..50]
	cpos := outputConn[num].byte
	bpos := outputConn[num].bit
	if cpos == 0 && bpos == 0 {
		fmt.Printf("invalid pin position, no data pin: %v\n", num)
		os.Exit(1)
	}
	queryStr := ""
	var i int
	for i = 1; i < cpos; i++ {
		queryStr += "0"
	}
	queryStr += bitPosToNibble(bpos - 1)
	for j := i + 1; j <= 10; j++ {
		queryStr += "0"
	}
	return queryStr
}

func bitPosToNibble(bpos int) string {
	ret := ""
	switch bpos {
	case 0:
		ret = "1"
		break
	case 1:
		ret = "2"
		break
	case 2:
		ret = "4"
		break
	case 3:
		ret = "8"
		break
	}
	return ret
}
