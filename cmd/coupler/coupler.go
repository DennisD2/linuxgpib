// Copyright 2023 Google LLC
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// version 2 as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

/*
Identify queries the identification string of a GPIB device.

It works for [SCPI] devices that support the *IDN? query.

Usage:

	identify [-verbose] [-board=BOARD] -address=ADDRESS

The flags are:

	-verbose
		Turn on logging. Without this, only the result or first error is printed.

	-board
		The board number. Defaults to zero, which corresponds to /dev/gpib0.

	-address
		The primary address of the GPIB device to query.

Examples:

	$ identify -address 22
	HEWLETT-PACKARD,34401A,0,10-5-2
	$

	$ identify -verbose -address 22
	2023/10/08 14:26:53 linuxgpib.go:132 Opened board 0 with version 4.3.4
	2023/10/08 14:26:53 linuxgpib.go:210 Opened address 22 (22/0) on board 0 as device 16
	2023/10/08 14:26:53 linuxgpib.go:434 Wrote "*IDN?\n" in 2ms to address 22
	2023/10/08 14:26:53 linuxgpib.go:333 Read "HEWLETT-PACKARD,34401A,0,10-5-2\n" in 35ms from address 22
	HEWLETT-PACKARD,34401A,0,10-5-2
	2023/10/08 14:26:53 linuxgpib.go:286 Closing address 22
	$

[SCPI]: https://en.wikipedia.org/wiki/Standard_Commands_for_Programmable_Instruments
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

	var num = 0
	vmax := 10 * 4
	var timesum time.Duration = 0.0
	for num < vmax {
		start := time.Now()

		// flip string
		//queryStr := ">>>>>>>>>>"
		//queryStr := "0080000000"
		queryStr := createQueryString(num)
		fmt.Printf("Query: %v - %v\n", num, queryStr)
		//if num%2 != 0 {
		//	queryStr = "0000000000"
		//}
		if _, err := fmt.Fprintln(d, queryStr); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write to device:", err)
			os.Exit(1)
		}

		elapsed := time.Since(start)
		timesum += elapsed
		num++

		if *count == 0 {
			time.Sleep(time.Duration(*wait) * time.Millisecond)
		} else {
			if num < *count {
				time.Sleep(time.Duration(*wait) * time.Millisecond)
			}
		}
	}
	var pp float32 = float32(timesum.Microseconds())
	pp = pp / float32(num)
	pp = pp / 1000.0
	if *verbose {
		fmt.Printf("Medium time in %v iterations: %v ms\n", num, pp)

		lines := internal.IblinesString(*board)
		fmt.Printf("GPIB Lines: %v\n", lines)
	}

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
