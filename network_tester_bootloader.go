package main

import (
  "flag"
  "log"
  "io"
  "github.com/Omegaice/go-xmodem/xmodem"
)

func main() {
    firmware := flag.String("firmware", "", "Path to the firmware.")
    device := flag.String("device", "", "Path to the device on which we should commmunicate")
    verbose := flag.Bool("verbose", false, "Whether to show verbose/debug log or not.")

    flag.Parse()

    // TODO: Set up logger

    // TODO : Loop through all the found devices and upload firmware
}

func upload_firmware(dev_path, firmware_path string) {
    // TODO : Implement fucntion
}
