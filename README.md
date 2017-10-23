# estim
A Go library for interfacing with e-stim units. Only works with the Erostek ET232 for now.

## Disclaimers
This software is still in early development, and has not been well tested.
Please carefully test before using. 

Until I'm more satisfied with the code, the external API is subject to change at any time. 

This project is not in any way associated or affiliated with Erostek. You may void your warranty, etc etc.

## Installation
```
go get -v github.com/biribiribiri/estim
```

## Example
TODO

## Connecting to the ET232 via Serial
It turns out that the audio/link port on the ET232 can be used to control the device. If the device completes a handshake when it powers on, then the port will be used as a serial input rather than as an audio input.

The device uses RS232, 19200/8/N/1. The pins are as follows: 
* Tip <-> RX
* Ring <-> TX
* Sleeve <-> Ground

Strangely enough, you can buy a suitable cable. [This is the one I bought.](https://www.amazon.com/gp/product/B004T9BBJC) You can connect it to a [USB-to-serial cable](https://www.amazon.com/gp/product/B0007OWNYA).


## Thanks
Many thanks to the Buttshock project for providing [documentation on the ET232](https://github.com/metafetish/buttshock-protocol-docs/blob/master/doc/et232-protocol.org).