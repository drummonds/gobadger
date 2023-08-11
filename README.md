# gobadger

Exploration of running go on the badger on a WSL Instance.  I am using Ubuntu 22.04

## Getting started.

- plug in badger to Windows USB
- install USB on WSL see below
- [install tinygo](https://tinygo.org/getting-started/install/linux/#ubuntu-debian)

This is a link to running tinygo on the badger https://tinygo.org/docs/reference/microcontrollers/badger2040/

The first step is getting example1 to run
### USB on WSL

To connect to the badger you need USB.  This is not mapped in wsl so you need to use usbip-win according to [Microsoft connect usb](https://learn.microsoft.com/en-us/windows/wsl/connect-usb).  This worked smoothly apart from having to reset the badger and use admin for windows command.

## Other Projects:

- https://github.com/strideynet/badger2040-go just getting started
- https://github.com/emmaly/tinygo-garden


