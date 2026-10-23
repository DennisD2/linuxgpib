/*
ICS4883 Instrument coupler test code
input tests - reading from input connector WITH SRQ
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
	count := flag.Int(
		"count", 1,
		"repeat count for query",
	)
	wait := flag.Int(
		"wait", 1,
		"wait time in seconds between queries",
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

	if *verbose {
		fmt.Printf("%v iterations, %vs wait time between iterations\n", *count, *wait)
	}

	// To get SRQ on input data available, we need to call 1x read before!
	var p []byte = make([]byte, 12)
	d.SetTimeout(100 * time.Millisecond)
	_, _ = d.Read(p)
	d.SetTimeout(1000 * time.Millisecond)
	//d.Clear()
	//d.Trigger()

	var num = 0

	wait_for_data_after_trigger := 0
	for true {
		num++
		if wait_for_data_after_trigger == 0 {
			// send device trigger
			//d.Trigger()
			wait_for_data_after_trigger = 1
		}
		time.Sleep(100 * time.Millisecond)

		/*state, glines := internal.Iblines(*board)
		if err := internal.Err(state); err != nil {
			error := err.Error()
			fmt.Printf("Error reading GPIB lines: %v\n", error)
			if strings.HasPrefix(error, "EDVR") {
				os.Exit(-1)
			}
		}
		lines := internal.FormatIblines(glines)

		srqStr := "LOW"
		if glines&0x2000 != 0 {
			srqStr = "HI"
		}
		fmt.Printf("%v: %x (%v) SRQ line is %v\n", num, glines, lines, srqStr)
		*/
		ret := internal.WaitSRQ(*board)
		if ret == 0 {
			fmt.Printf("%5v no srq...\n", num)
			continue
		} else {
			fmt.Printf("ret=%v\n", ret)
		}
		//sbyte := internal.ReadStatusByte(*board, *address)
		//fmt.Printf("Status byte: %x\n", sbyte)
		/*result := internal.WaitSRQ(*board)
		if result != 1 {
			println("no srq...\n")
			continue
		}*/

		/*haveSRQ := internal.TestSRQ(*board)
		if haveSRQ == 1 {
			println("HAVE SRQ")
		}*/
		/*var n int
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
		}*/
		//pins := prettyPrintPins(p)
		//println(pins)
		//wait_for_data_after_trigger = 0
		//if srqStr == "HI" {
		//sbyte := internal.ReadStatusByte(*board, *address)
		//ibsta, status := internal.Ibrsp(*address)
		println("Serial Poll...")
		var status byte
		status, err = d.Spoll()
		if err != nil {
			fmt.Printf("Error reading GPIB lines: %v\n", err)
		}
		fmt.Printf("Status byte: %x\n", status)
		if status == 0x41 {
			fmt.Printf("Data available!\n")
			_, _ = d.Read(p)
			fmt.Printf("Read: %x,%x,%x,%x,%x,%x,%x,%x,%x,%x,%x, %x\n", p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7], p[8], p[9], p[10], p[11])
			pins := prettyPrintPins(p)
			println(pins)
			internal.Ibonl(*address, 1)
		}
		//os.Exit(0)

		wait_for_data_after_trigger = 0
		//}

	}

}

func prettyPrintPins(p []byte) string {
	result := "Read: "
	i := 0
	for i < len(p)-2 {
		bitValue := bitvalue(p[i])
		result += fmt.Sprintf("%v,", bitValue)
		i++
	}
	//result = fmt.Sprintf("Read: %x,%x,%x,%x,%x,%x,%x,%x,%x,%x,%x, %x", p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7], p[8], p[9], p[10], p[11])
	return result
}

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
