# estim
A Go library for interfacing with e-stim units.

# Installation
```
go get -v github.com/biribiribiri/estim
```

# Connecting to the ET232 via Serial
It turns out that the audio/link port on the ET232 can be used to control the device. If the device completes a handshake when it powers on, then the port will be used as a serial input rather than as an audio input.

The device uses RS232, 19200/8/N/1. The pins are as follows: 
* Tip <-> RX
* Ring <-> TX
* Sleeve <-> Ground

Strangely enough, you can buy a suitable cable. [This is the one I bought.](https://www.amazon.com/gp/product/B004T9BBJC) You can connect it to a [USB-to-serial cable](https://www.amazon.com/gp/product/B0007OWNYA).