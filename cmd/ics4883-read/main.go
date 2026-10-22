/*
ICS4883 Instrument coupler test code
input tests - reading from input connector
*/

// Read data from ICS electronics 4880/4883 Instrument Coupler

// We just read in all 40 bit values in an endless loop.
// Pin 2 (Data available) needs a LO-HI edge to read in new values at pins.
// ~~Without LO-HI edge, last latched value will be returned by 4883.~~ <-- I do not think so anymore. Nothing is returned then!
// If there is no value, we will get timeout in read() and 0 bytes as return bytes.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/msiegen/linuxgpib"
)

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
		"The primary address of the GPIB device to query.",
	)

	flag.Parse()

	if *address < 0 {
		fmt.Fprintln(os.Stderr, "Please specify an -address!")
		os.Exit(1)
	}

	var opts []linuxgpib.Option
	if *verbose {
		opts = append(opts, linuxgpib.Log(log.Default()))
		opts = append(opts, linuxgpib.Timeout(1000*time.Millisecond))
	}

	d, err := linuxgpib.NewDevice(*board, *address, opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open device:", err)
		os.Exit(1)
	}
	defer d.Close()

	var p []byte = make([]byte, 12)
	loops := 1
	for true {
		var n int
		n, err = d.Read(p)
		if err != nil {
			e := fmt.Sprintf("%s", err)
			if e != "timed out" {
				// expected "error" is timeout. This is not really an error, but means
				// that just no new data arrived (i.e. no LO->HO transition at pin 2)
				// All other errors are real, unexpected errors
				fmt.Fprintln(os.Stderr, "Failed to read from device:", err)
				os.Exit(1)
			}
			time.Sleep(time.Second * 1)
		}
		if n == 0 {
			continue
		}
		pins := prettyPrintPins(p)
		fmt.Printf("%v\tRead %d bytes: %v\n", loops, n, pins)

		loops++
		time.Sleep(100 * time.Millisecond)
	}
}

// prettyPrintPins print bit values
func prettyPrintPins(p []byte) string {
	result := ""
	i := 0
	for i < len(p)-2 {
		bitValue := bitvalue(p[i])
		commaStr := ""
		if i < len(p)-3 {
			commaStr = ","
		}
		result += fmt.Sprintf("%v%v", bitValue, commaStr)
		i++
	}
	return result
}

// bitvalue format bit value according to 4883 coupler output data
func bitvalue(b byte) string {
	bvalue := fmt.Sprintf("%x", b)
	switch b {
	case 0x30:
		bvalue = "0000"
	case 0x31:
		bvalue = "0001"
	case 0x32:
		bvalue = "0010"
	case 0x33:
		bvalue = "0011"
	case 0x34:
		bvalue = "0100"
	case 0x35:
		bvalue = "0101"
	case 0x36:
		bvalue = "0110"
	case 0x37:
		bvalue = "0111"
	case 0x38:
		bvalue = "1000"
	case 0x39:
		bvalue = "1001"
	case 0x26:
		bvalue = "1010"
	case 0x2d:
		bvalue = "1011"
	case 0x2e:
		bvalue = "1100"
	case 0x2c:
		bvalue = "1101"
	case 0x45:
		bvalue = "1110"
	case 0x20:
		bvalue = "1111"
	}
	return bvalue
}
