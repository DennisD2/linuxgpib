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

	d, err := linuxgpib.NewDevice(*board, *address, opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open device:", err)
		os.Exit(1)
	}
	defer d.Close()

	if *blink == false {
		fmt.Printf("Set pin %v to value %v \n", *pin, *value)
	} else {
		fmt.Printf("blinking pin %v\n", *pin)
	}

	if *blink == false {
		//queryStr := ">>>>>>>>>>"
		queryStr := "0000000000"
		if *value == 1 {
			queryStr = createQueryString(*pin)
		}

		fmt.Printf("Query: %v - %v\n", *pin, queryStr)
		if _, err := fmt.Fprintln(d, queryStr); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write to device:", err)
			os.Exit(1)
		}
	} else {
		offStr := "0000000000"
		onStr := createQueryString(*pin)
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

func createQueryString(num int) string {
	//queryStr := "0080000000"
	cpos := num / 4
	bpos := num % 4
	queryStr := ""
	var i int
	for i = 0; i < cpos; i++ {
		queryStr += "0"
	}
	queryStr += bitPosToNibble(bpos)
	for j := i + 1; j < 10; j++ {
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

func doContinue(num int, c int) bool {
	if c != 0 {
		return num < c
	} else {
		return true
	}
}
