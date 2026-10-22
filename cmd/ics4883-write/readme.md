# ICS4883 write code
This is test code to access a ICS 4883 instrument coupler in write mode.

Code uses linuxgpib go API. Device is controlled via GPIB.
Information on that API: https://github.com/msiegen/linuxgpib

ICS4883 accepts strings in formast "dddddddddd\n\r". 'd' is some digit like
0-9 or a few alphanumeric chars. Only the lower 4 bits of a character byte
are used and forwarded to the output pins. So each character send corresponds 
to 4 bits at the output connector.

This small test app allows to send such command strings, which will switch on/off
the pins at output connector.

## Build
```shell
go build .
```

## Run
Use ```address``` option for GPIB device address. 

#### Send a complete command string to device
Send a complete 10 digit command string like "0000010000" to device
Use ```data``` to define the command string.
Use ```blink``` to do this in an endless loop. 

Example:
```shell
./ics4883-write -address=16 -data 0000010000 -blink -verbose=true
```

#### Switch on/off a pin
Switch on/off a pin by giving its number and a value (0 or 1).

Example:
```shell
./ics4883-write -address=16 -pin 13 -value=1 -verbose=true
./ics4883-write -address=16 -pin 13 -value=0 -verbose=true
```

#### Blinking
Use ```blink``` to repeat in endless loop. The command string is sent, after 100 ms
a string with all zeros "0000000000" is sent to clear all pins, another 100 ms is waited.
The loop starts again. Interrupt with CTRL-C.

## Notes, restrictions
Besides the 40 data pins, there are 10 other pins. Two (1,5) have GND and +5V.
4 others are output pins for handshake, and there are 3 input pins for handshake.
Finally, there is one "handshake enable" like input pin. It is not known how to use these
handshake lines.
