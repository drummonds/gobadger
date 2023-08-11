# gobadger

Exploration of running go on the badger on a WSL Instance.  I am using Ubuntu 22.04

## Getting started.

- plug in badger to Windows USB
- install USB on WSL see below
- [install tinygo](https://tinygo.org/getting-started/install/linux/#ubuntu-debian)
- build local
- Manual deploy

This is a link to running tinygo on the badger https://tinygo.org/docs/reference/microcontrollers/badger2040/

The first step is getting example1 to run
### USB on WSL

To connect to the badger you need USB.  This is not mapped in wsl so you need to use usbip-win according to [Microsoft connect usb](https://learn.microsoft.com/en-us/windows/wsl/connect-usb).  This worked smoothly apart from having to reset the badger and use admin for windows command.

## Other Projects:

- https://github.com/strideynet/badger2040-go just getting started
- https://github.com/emmaly/tinygo-garden

## Buid local

Using `tinygo build -target=badger2040 -o example1 .`  Have to use . as target rather than ./main.go as shapes is part of main package.  End up with  a 734k package example1.

Then to flash it is:

`tinygo flash -target=badger2040  .`


### Getting the serial port to work
The first time I tired the flash command it didn't work.  I had attached the USB device as per  [Microsoft connect usb][]  and when I did (every time you reset the Badger board it reset):

`lsbusb`

I get a listing showing the board is connected:

``` bash
Bus 001 Device 004: ID 2e8a:0005 MicroPython Board in FS mode`
```

Trying to find the tty port:

`dmesg | grep tty`

Installed minicom and changed the setting port to dev/ttyACM0 and that did connect although I couldn't do anything.  So retry flashing

`tinygo flash -target=badger2040 -port=/dev/ttyACM0  .`

And now get:

`error: failed to flash /tmp/tinygo2787660771/main.uf2: unable to locate any volume: [RPI-RP2]`

 tinygo flash -target=badger2040 -o /run/media/$USER/PYBADGEBOOT/firmware.uf2 -port=/dev/ttyACM0  .

Maybe this is it -
As with any RP2040 based wotsit you just need to hold down the boot select (labelled BOOT/USR on Badger 2040) while you reset the board. That will bring it up in DFU mode which appears as a disk and you can drag and drop the .uf2 onto it.

## Manual deploy
A way of loading images is t use the UF2.  If you build with tinygo

`tinygo build -target=badger2040 -o badger2040_button.uf2 .`

- Then move to the windows machine
- then reboot the badger (hold down BOOT/USR then press reset)  the drive pops as a dmass storage device.
- Then copy the UF2 file to the drive
- the rpi reboots using the image
- Success!




