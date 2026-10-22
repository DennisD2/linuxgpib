# ICS4883 write code
This is test code to access a ICS 4883 instrument coupler.

Code uses linuxgpib go API. Device is controlled via GPIB.

## Build
```shell
go build .
```

## Run
Use ```address``` option for GPIB device address. 

#### Send a complete command string to device
Send a complete 10 digit command string like "0000010000" to device
Use "-blink" to do this in an endless loop. 

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
Use "-blink" to repeat in endless loop. The command string is sent, after 100 ms
a string with all zeros "0000010000" is sent to clear all pins, another 100 ms is waited.
The loop starts again. Interrupt with CTRL-C.