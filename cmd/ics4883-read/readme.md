# ICS4883 read code
This is test code to access a ICS 4883 instrument coupler in read mode.

Code uses linuxgpib go API. Device is controlled via GPIB.
Information on that API: https://github.com/msiegen/linuxgpib

ICS4883 needs a LO->HI edge at pin 2 of input connector to read in data
at input pins. So pin 2 is a "DATA READY" input pin. The external instrument
signals with this pin that it has provided stable data at data input pins.

If there is no data available, the read is timing out,

If there is data available (i.e. we had a LO->HI edge at pin 2), the data 
is returned in a 12 byte sequence in format "dddddddddd\n\r". 'd' is some digit like
0-9 or a few alphanumeric chars. Each digit corresponds 
to 4 bits at the input connector.

This small test app allows to read such strings.

There is also the possibility to use GPIB SRQ mechanism for reading (getting a service
request from device if data is available). This mechanism is tested in another test app.

## Build
```shell
go build .
```

## Run
Use ```address``` option for GPIB device address. 

Example:
```shell
./ics4883-read -address=16 -verbose=true
```
Example output:
```shell
./ics4883-read -address=16 -verbose=true
2026/10/22 15:43:55 linuxgpib.go:132 Opened board 0 with version 4.3.7
2026/10/22 15:43:55 linuxgpib.go:210 Opened address 16 (16/0) on board 0 as device 16
2026/10/22 15:43:56 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
2026/10/22 15:43:58 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
2026/10/22 15:43:59 linuxgpib.go:333 Read "  E       \r\n" in 8ms from address 16
Read 12 bytes
1:      Read: 1111,1111,1110,1111,1111,1111,1111,1111,1111,1111,
2026/10/22 15:44:00 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
2026/10/22 15:44:01 linuxgpib.go:333 Read "  E       \r\n" in 9ms from address 16
Read 12 bytes
2:      Read: 1111,1111,1110,1111,1111,1111,1111,1111,1111,1111,
2026/10/22 15:44:02 linuxgpib.go:333 Read "  E       \r\n" in 690ms from address 16
Read 12 bytes
3:      Read: 1111,1111,1110,1111,1111,1111,1111,1111,1111,1111,
2026/10/22 15:44:02 linuxgpib.go:333 Read "  E       \r\n" in 9ms from address 16

...
2026/10/22 15:53:39 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
2026/10/22 15:53:41 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
2026/10/22 15:53:56 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
2026/10/22 15:53:58 linuxgpib.go:331 Failed to read from address 16 device 16: timed out
...
```
## Notes, restrictions
Besides the 40 data pins, there are 10 other pins. Two (1,5) have GND and +5V.
4 others are output pins for handshake, and there are 4 input pins for handshake.
Pin 2 is "DATA READY", pin 6 is "/DATA READY". The functionality of remaining 6 pins is yet unknown.

## Links
* https://linux-gpib.sourceforge.io/doc_html/t1.html
* 